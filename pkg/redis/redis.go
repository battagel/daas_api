package redis

import (
	"context"
	"daas_api/pkg/logger"
	"encoding/json"
	"sync"

	"github.com/redis/go-redis/v9"
)

type CloseFunc func() error

type Redis struct {
	logger logger.Logger
	client *redis.Client
	ctx context.Context
}

func CreateRedis(logger logger.Logger, ctx context.Context, address string, password string) (*Redis, CloseFunc, error){
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	status, err := client.Ping(ctx).Result()
    if (err != nil || status != "PONG") {
        logger.Errorln("Redis connection was refused")
		return nil, nil, err
    }
	return &Redis{
		logger: logger,
		client: client,
		ctx: ctx,
	}, client.Close, nil
}

func (r *Redis) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-r.ctx.Done():
		err := r.client.Close();
		if err != nil {
			r.logger.Errorln("Failed to close Redis")
		}
		r.logger.Debugln("Redis gracefully shut down")
		return
	}
}

func (r *Redis) Insert(key string, jsonValue map[string]interface{}) error {
	value, err := json.Marshal(jsonValue)
	if err != nil {
		return err
	}
	err = r.client.Set(r.ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(key string) (map[string]interface{}, error) {
	value, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Debugw("Key does not exist",
				"key", key,
			)
			return nil, nil
		}
		// Other errors
		r.logger.Errorw("Error querying from Redis",
			"error", err,
		)
		return nil, err
	}


	var jsonValue map[string]interface{}
	err = json.Unmarshal([]byte(value), &jsonValue)
	if err != nil {
		r.logger.Errorw("Could not unmarshall",
			"error", err,
			"value", value,
		)
		return nil, err
	}

	return jsonValue, nil
}

func (r *Redis) GetAll() ([]map[string]interface{}, error) {
	keys, err := r.GetKeys()
	if err != nil {
		r.logger.Errorw("Error querying from Redis",
			"error", err,
		)
		return nil, err
	}
	if len(keys) < 1 {
		r.logger.Debugln("No keys found")
		return nil, nil
	}
	r.logger.Debugw("Got keys",
		"keys", keys,
	)

    values, err := r.client.MGet(r.ctx, keys...).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Debugw("Keys do not exist",
				"keys", keys,
			)
			return nil, nil
		}
		// Other errors
		r.logger.Errorw("Error querying from Redis",
			"error", err,
		)
		return nil, err
	}

    var result []map[string]interface{}

    for _, v := range values {
        if v == nil {
            continue
        }

        // Convert the JSON string to a map
        var jsonValue map[string]interface{}
        err := json.Unmarshal([]byte(v.(string)), &jsonValue)
        if err != nil {
            r.logger.Errorw("Could not unmarshall",
                "error", err,
                "value", v,
            )
            return nil, err
        }

        result = append(result, jsonValue)
    }

    return result, nil
}

func (r *Redis) GetKeys() ([]string, error) {
	// Maybe convery this to scan to reduce load
	keys, err := r.client.Keys(r.ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	return keys,nil
}

func (r *Redis) Delete(key string) error {
	_, err := r.client.Del(r.ctx, key).Result()
	if err != nil {
		return err
	}
	r.logger.Debugw("Deleted phrase",
		"phrase", key,
	)
	return nil
}

package dynamoDB

import (
	"daas_api/pkg/logger"
	"fmt"
)

type DynamoDB struct {
	logger logger.Logger
}


func CreateDynamoDB(logger logger.Logger) (*DynamoDB, error) {
	return &DynamoDB{logger: logger}, nil
}

func (db *DynamoDB) Insert() error {
	fmt.Println("Debug database insert")
	return nil

}

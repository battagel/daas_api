package phrase

import (
	"errors"
	"fmt"
)


type Explanation struct {
    Definition  string   `json:"definition"`
    Code        []string `json:"code"`
    References  []string `json:"references"`
}

type Phrase struct {
    Phrase      string         `json:"phrase"`
    Terms       []string       `json:"terms"`
    LastUpdate  string         `json:"last_update"`
    Relevance   float64        `json:"relevance"`
    // Accuracy    float64       `json:"accuracy"`
    Tags         []string      `json:"tags"`
    Explanations []Explanation `json:"explanations"`
}

func (p *Phrase) ToMap() map[string]interface{} {
    phraseMap := make(map[string]interface{})
    phraseMap["phrase"] = p.Phrase
    phraseMap["terms"] = p.Terms
    phraseMap["last_update"] = p.LastUpdate
    phraseMap["relevance"] = p.Relevance
    // phraseMap["accuracy"] = p.Accuracy
    phraseMap["tags"] = p.Tags

    // Convert the Explanation slice to a slice of maps
    explanations := make([]map[string]interface{}, len(p.Explanations))
    for i, exp := range p.Explanations {
        explanationMap := make(map[string]interface{})
        explanationMap["definition"] = exp.Definition
        explanationMap["code"] = exp.Code
        explanationMap["references"] = exp.References
        explanations[i] = explanationMap
    }
    phraseMap["explanations"] = explanations

    return phraseMap
}

// Turn raw database data into a phrase
func (p *Phrase) ToPhrase(rawData interface{}) error {
	switch data := rawData.(type) {
	case map[string]interface{}:
		// Use type assertions to convert the data from the map into the struct

		// Check and assign "phrase" field
		if val, ok := data["phrase"].(string); ok {
			p.Phrase = val
		} else {
			return errors.New("Invalid or missing 'phrase' field")
		}

		// Check and assign "terms" field
		p.Terms = toStringSlice(data["terms"].([]interface{}))

		// Check and assign "last_update" field
		if val, ok := data["last_update"].(string); ok {
			p.LastUpdate = val
		} else {
			return errors.New("Invalid or missing 'last_update' field")
		}

		// Check and assign "relevance" field
		if val, ok := data["relevance"].(float64); ok {
			p.Relevance = val
		} else {
			return errors.New("Invalid or missing 'relevance' field")
		}

		// // Check and assign "accuracy" field
		// if val, ok := data["accuracy"].(float64); ok {
		// 	phrase.Accuracy = int(val)
		// } else {
		// 	return errors.New("Invalid or missing 'accuracy' field")
		// }

		// Check and assign "tag" field
		p.Tags = toStringSlice(data["tags"].([]interface{}))

		// Convert the "explanation" field (slice of maps) back to the Explanation struct
		if expData, ok := data["explanations"].([]interface{}); ok {
			for _, exp := range expData {
				if expMap, ok := exp.(map[string]interface{}); ok {
					explanation := Explanation{
						Definition: expMap["definition"].(string),
						Code:       toStringSlice(expMap["code"].([]interface{})),
						References: toStringSlice(expMap["references"].([]interface{})),
					}
					p.Explanations = append(p.Explanations, explanation)
				} else {
					return errors.New("Invalid explanation data")
				}
			}
		} else {
			return errors.New("Invalid or missing 'explanation' field")
		}

	default:
		return errors.New("Unexpected type for rawData")
	}

	fmt.Println(p)
	return nil
}

// toStringSlice converts an interface{} slice to a []string slice
func toStringSlice(data []interface{}) []string {
    result := make([]string, len(data))
    for i, v := range data {
        result[i] = v.(string)
    }
    return result
}

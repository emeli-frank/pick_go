package validation

import (
	"encoding/json"
)

type Errors []struct{
	Key string
	Message string
}

func (e Errors) Error() string {
	return "validation error"
}

func (e *Errors) Add(key string, message string, op string, err error) {
	*e = append(*e, struct {
		Key string
		Message string
	}{Key:key, Message:message})
}

func (e Errors) MarshalJSON() ([]byte, error) {
	type errorlet struct {
		Key string `json:"key"`
		Message string `json:"message"`
	}

	var elets []errorlet
	for _, elet := range e {
		elets = append(elets, errorlet{Key: elet.Key, Message: elet.Message})
	}

	return json.Marshal(struct {
		Type string `json:"type"`
		Errors []errorlet `json:"errors"`
	}{Type:"validation", Errors:elets})
}

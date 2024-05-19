package models

import (
	"encoding/json"
)

type Message struct {
    Id      int    
    Name    string `json:"name"`   
    Message string `json:"message"`
}

func ParseJSONMessage(jsonData []byte) (Message, error) {
    var message Message
    err := json.Unmarshal(jsonData, &message)
    return message, err
}
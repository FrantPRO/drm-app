package data

import (
	"context"
	"encoding/json"
)

type Command struct {
	Action   string                 `json:"action"`
	Entity   string                 `json:"entity"`
	Data     map[string]interface{} `json:"data"`
	UserID   string                 `json:"user_id"`
	UserRole string                 `json:"user_role"`
}

type DataExecutor interface {
	ExecuteCommand(ctx context.Context, command *Command) (interface{}, error)
}

func formatData(data map[string]interface{}) string {
	if data == nil {
		return "{}"
	}
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

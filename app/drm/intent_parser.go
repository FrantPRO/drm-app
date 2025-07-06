package drm

import (
	"encoding/json"
	"fmt"
	"strings"
)

type IntentParser struct{}

func NewIntentParser() *IntentParser {
	return &IntentParser{}
}

func (p *IntentParser) Parse(query string) (*Command, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	var command Command

	if strings.Contains(query, "create") || strings.Contains(query, "add") {
		command.Action = "create"
	} else if strings.Contains(query, "read") || strings.Contains(query, "get") || strings.Contains(query, "list") || strings.Contains(query, "show") {
		command.Action = "read"
	} else if strings.Contains(query, "update") || strings.Contains(query, "modify") || strings.Contains(query, "change") {
		command.Action = "update"
	} else if strings.Contains(query, "delete") || strings.Contains(query, "remove") {
		command.Action = "delete"
	} else {
		command.Action = "read"
	}

	// Extract the command part (before json:)
	commandPart := query
	if jsonIndex := strings.Index(query, "json:"); jsonIndex != -1 {
		commandPart = query[:jsonIndex]
	}
	
	if strings.Contains(commandPart, "user") {
		command.Entity = "user"
	} else if strings.Contains(commandPart, "product") {
		command.Entity = "product"
	} else if strings.Contains(commandPart, "order") {
		command.Entity = "order"
	} else {
		return nil, fmt.Errorf("unknown entity in query: %s", query)
	}

	command.Data = make(map[string]interface{})
	
	if strings.Contains(query, "json:") {
		jsonStart := strings.Index(query, "json:")
		if jsonStart != -1 {
			jsonData := query[jsonStart+5:]
			if err := json.Unmarshal([]byte(jsonData), &command.Data); err != nil {
				return nil, fmt.Errorf("invalid JSON data: %w", err)
			}
		}
	}

	return &command, nil
}
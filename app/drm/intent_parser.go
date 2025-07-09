package drm

import (
	"encoding/json"
	"fmt"
	"strings"

	"drm-app/app/data"
)

type IntentParser struct{}

func NewIntentParser() *IntentParser {
	return &IntentParser{}
}

func (p *IntentParser) Parse(query string) (*data.Command, error) {
	query = strings.TrimSpace(strings.ToLower(query))

	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	var command data.Command

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
	commandPart := strings.TrimSpace(query)
	if jsonIndex := strings.Index(query, "json:"); jsonIndex != -1 {
		commandPart = strings.TrimSpace(query[:jsonIndex])
	}

	// Use word boundary detection for more precise matching
	commandWords := strings.Fields(commandPart)
	entityFound := false

	for _, word := range commandWords {
		switch strings.ToLower(word) {
		case "user", "users":
			command.Entity = "user"
			entityFound = true
		case "product", "products":
			command.Entity = "product"
			entityFound = true
		case "order", "orders":
			command.Entity = "order"
			entityFound = true
		}
	}

	if !entityFound {
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

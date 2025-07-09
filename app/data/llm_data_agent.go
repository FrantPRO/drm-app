package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

type LLMDataAgent struct {
	data   map[string]map[string]interface{}
	client *api.Client
}

func NewLLMDataAgent() *LLMDataAgent {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		client = nil
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		
		err = client.Heartbeat(ctx)
		if err != nil {
			client = nil
		}
	}

	return &LLMDataAgent{
		data: map[string]map[string]interface{}{
			"user": {
				"1": map[string]interface{}{
					"id":         "1",
					"name":       "John Doe",
					"email":      "john@example.com",
					"created_at": time.Now(),
				},
				"2": map[string]interface{}{
					"id":         "2",
					"name":       "Jane Smith",
					"email":      "jane@example.com",
					"created_at": time.Now(),
				},
			},
			"product": {
				"1": map[string]interface{}{
					"id":          "1",
					"name":        "Laptop",
					"price":       999.99,
					"description": "Gaming laptop",
					"created_at":  time.Now(),
				},
				"2": map[string]interface{}{
					"id":          "2",
					"name":        "Mouse",
					"price":       29.99,
					"description": "Wireless mouse",
					"created_at":  time.Now(),
				},
			},
			"order": {},
		},
		client: client,
	}
}

func (d *LLMDataAgent) ExecuteCommand(ctx context.Context, command *Command) (interface{}, error) {
	if d.client == nil {
		return d.fallbackExecution(command)
	}
	
	prompt := d.buildPrompt(command)
	
	response, err := d.queryLLM(ctx, prompt)
	if err != nil {
		return d.fallbackExecution(command)
	}
	
	return d.executeFromLLMResponse(command, response)
}

func (d *LLMDataAgent) buildPrompt(command *Command) string {
	dataJSON, _ := json.Marshal(d.data)
	
	prompt := fmt.Sprintf(`You are a data management assistant. Based on the command provided, execute the appropriate CRUD operation.

Current data state:
%s

Command details:
- Action: %s
- Entity: %s
- Data: %s
- UserID: %s
- UserRole: %s

Instructions:
1. Analyze the command and determine the appropriate action
2. For CREATE: Generate a new ID and add timestamps
3. For READ: Return the requested data or all items if no ID specified
4. For UPDATE: Modify existing data with new values and add updated_at timestamp
5. For DELETE: Remove the specified item
6. Return your response in JSON format with the following structure:
   {"action": "create|read|update|delete", "success": true|false, "data": {...}, "error": "error message if any"}

Respond only with valid JSON.`, 
		string(dataJSON), command.Action, command.Entity, 
		formatData(command.Data), command.UserID, command.UserRole)
	
	return prompt
}

func (d *LLMDataAgent) queryLLM(ctx context.Context, prompt string) (string, error) {
	timeout := time.Second * 5
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	
	req := &api.GenerateRequest{
		Model:  "llama3.2:1b",
		Prompt: prompt,
		Stream: &[]bool{false}[0],
	}
	
	var response strings.Builder
	err := d.client.Generate(timeoutCtx, req, func(resp api.GenerateResponse) error {
		response.WriteString(resp.Response)
		return nil
	})
	
	if err != nil {
		return "", err
	}
	
	return response.String(), nil
}

func (d *LLMDataAgent) executeFromLLMResponse(command *Command, llmResponse string) (interface{}, error) {
	var response struct {
		Action  string      `json:"action"`
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Error   string      `json:"error"`
	}
	
	if err := json.Unmarshal([]byte(llmResponse), &response); err != nil {
		return d.fallbackExecution(command)
	}
	
	if !response.Success {
		return nil, fmt.Errorf("LLM execution failed: %s", response.Error)
	}
	
	switch command.Action {
	case "create":
		return d.create(command.Entity, command.Data)
	case "read":
		return d.read(command.Entity, command.Data)
	case "update":
		return d.update(command.Entity, command.Data)
	case "delete":
		return d.delete(command.Entity, command.Data)
	default:
		return nil, fmt.Errorf("unsupported action: %s", command.Action)
	}
}

func (d *LLMDataAgent) fallbackExecution(command *Command) (interface{}, error) {
	switch command.Action {
	case "create":
		return d.create(command.Entity, command.Data)
	case "read":
		return d.read(command.Entity, command.Data)
	case "update":
		return d.update(command.Entity, command.Data)
	case "delete":
		return d.delete(command.Entity, command.Data)
	default:
		return nil, fmt.Errorf("unsupported action: %s", command.Action)
	}
}

func (d *LLMDataAgent) create(entity string, data map[string]interface{}) (interface{}, error) {
	if d.data[entity] == nil {
		d.data[entity] = make(map[string]interface{})
	}
	
	id := fmt.Sprintf("%d", len(d.data[entity])+1)
	data["id"] = id
	data["created_at"] = time.Now()
	
	d.data[entity][id] = data
	return data, nil
}

func (d *LLMDataAgent) read(entity string, data map[string]interface{}) (interface{}, error) {
	if id, ok := data["id"].(string); ok {
		if item, exists := d.data[entity][id]; exists {
			return item, nil
		}
		return nil, fmt.Errorf("item not found")
	}
	
	var results []interface{}
	for _, item := range d.data[entity] {
		results = append(results, item)
	}
	return results, nil
}

func (d *LLMDataAgent) update(entity string, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required for update")
	}
	
	if _, exists := d.data[entity][id]; !exists {
		return nil, fmt.Errorf("item not found")
	}
	
	for key, value := range data {
		if key != "id" {
			d.data[entity][id].(map[string]interface{})[key] = value
		}
	}
	d.data[entity][id].(map[string]interface{})["updated_at"] = time.Now()
	
	return d.data[entity][id], nil
}

func (d *LLMDataAgent) delete(entity string, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required for delete")
	}
	
	if _, exists := d.data[entity][id]; !exists {
		return nil, fmt.Errorf("item not found")
	}
	
	delete(d.data[entity], id)
	return map[string]string{"message": "deleted successfully"}, nil
}

func formatData(data map[string]interface{}) string {
	if data == nil {
		return "{}"
	}
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
package data

import (
	"context"
	"fmt"
	"time"
)

type Command struct {
	Action   string                 `json:"action"`
	Entity   string                 `json:"entity"`
	Data     map[string]interface{} `json:"data"`
	UserID   string                 `json:"user_id"`
	UserRole string                 `json:"user_role"`
}

type DataAgent struct {
	data map[string]map[string]interface{}
}

func NewDataAgent() *DataAgent {
	return &DataAgent{
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
	}
}

func (d *DataAgent) ExecuteCommand(ctx context.Context, command *Command) (interface{}, error) {
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

func (d *DataAgent) create(entity string, data map[string]interface{}) (interface{}, error) {
	if d.data[entity] == nil {
		d.data[entity] = make(map[string]interface{})
	}
	
	id := fmt.Sprintf("%d", len(d.data[entity])+1)
	data["id"] = id
	data["created_at"] = time.Now()
	
	d.data[entity][id] = data
	return data, nil
}

func (d *DataAgent) read(entity string, data map[string]interface{}) (interface{}, error) {
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

func (d *DataAgent) update(entity string, data map[string]interface{}) (interface{}, error) {
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

func (d *DataAgent) delete(entity string, data map[string]interface{}) (interface{}, error) {
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
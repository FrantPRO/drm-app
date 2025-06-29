package drm

import (
	"fmt"
)

type LogicAgent struct {
	rules map[string]map[string]func(map[string]interface{}) error
}

func NewLogicAgent() *LogicAgent {
	agent := &LogicAgent{
		rules: make(map[string]map[string]func(map[string]interface{}) error),
	}
	
	agent.loadRules()
	return agent
}

func (l *LogicAgent) loadRules() {
	l.rules["user"] = map[string]func(map[string]interface{}) error{
		"create": l.validateUserCreate,
		"update": l.validateUserUpdate,
	}
	
	l.rules["product"] = map[string]func(map[string]interface{}) error{
		"create": l.validateProductCreate,
		"update": l.validateProductUpdate,
	}
	
	l.rules["order"] = map[string]func(map[string]interface{}) error{
		"create": l.validateOrderCreate,
	}
}

func (l *LogicAgent) ValidateCommand(command *Command) error {
	entityRules, exists := l.rules[command.Entity]
	if !exists {
		return nil
	}
	
	validator, exists := entityRules[command.Action]
	if !exists {
		return nil
	}
	
	return validator(command.Data)
}

func (l *LogicAgent) validateUserCreate(data map[string]interface{}) error {
	if name, ok := data["name"].(string); !ok || name == "" {
		return fmt.Errorf("user name is required")
	}
	if email, ok := data["email"].(string); !ok || email == "" {
		return fmt.Errorf("user email is required")
	}
	return nil
}

func (l *LogicAgent) validateUserUpdate(data map[string]interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided for update")
	}
	return nil
}

func (l *LogicAgent) validateProductCreate(data map[string]interface{}) error {
	if name, ok := data["name"].(string); !ok || name == "" {
		return fmt.Errorf("product name is required")
	}
	if price, ok := data["price"].(float64); !ok || price <= 0 {
		return fmt.Errorf("product price must be positive")
	}
	return nil
}

func (l *LogicAgent) validateProductUpdate(data map[string]interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided for update")
	}
	return nil
}

func (l *LogicAgent) validateOrderCreate(data map[string]interface{}) error {
	if items, ok := data["items"].([]interface{}); !ok || len(items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}
	return nil
}
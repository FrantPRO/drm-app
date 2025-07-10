package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ollama/ollama/api"
	"strings"
	"time"

	"drm-app/app/db"
)

type PostgresLLMDataAgent struct {
	db     *db.Database
	client *api.Client
}

func NewPostgresLLMDataAgent(database *db.Database) *PostgresLLMDataAgent {
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

	return &PostgresLLMDataAgent{
		db:     database,
		client: client,
	}
}

func (p *PostgresLLMDataAgent) ExecuteCommand(ctx context.Context, command *Command) (interface{}, error) {
	if p.client == nil {
		return p.fallbackExecution(ctx, command)
	}

	prompt, err := p.buildPrompt(ctx, command)
	if err != nil {
		return p.fallbackExecution(ctx, command)
	}

	response, err := p.queryLLM(ctx, prompt)
	if err != nil {
		return p.fallbackExecution(ctx, command)
	}

	return p.executeFromLLMResponse(ctx, command, response)
}

func (p *PostgresLLMDataAgent) buildPrompt(ctx context.Context, command *Command) (string, error) {
	currentData, err := p.getCurrentDataState(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get current data state: %w", err)
	}

	dataJSON, _ := json.Marshal(currentData)

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

	return prompt, nil
}

func (p *PostgresLLMDataAgent) getCurrentDataState(ctx context.Context) (map[string]interface{}, error) {
	currentData := map[string]interface{}{}

	// Get sample of users
	var users []User
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id LIMIT 5`
	err := p.db.DB.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	currentData["users"] = users

	// Get sample of products
	var products []Product
	query = `SELECT id, name, price, description, created_at, updated_at FROM products ORDER BY id LIMIT 5`
	err = p.db.DB.SelectContext(ctx, &products, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	currentData["products"] = products

	// Get sample of orders
	var orders []Order
	query = `SELECT id, user_id, items, total_amount, status, created_at, updated_at FROM orders ORDER BY id LIMIT 5`
	err = p.db.DB.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	currentData["orders"] = orders

	return currentData, nil
}

func (p *PostgresLLMDataAgent) queryLLM(ctx context.Context, prompt string) (string, error) {
	timeout := time.Second * 5
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &api.GenerateRequest{
		Model:  "llama3.2:1b",
		Prompt: prompt,
		Stream: &[]bool{false}[0],
	}

	var response strings.Builder
	err := p.client.Generate(timeoutCtx, req, func(resp api.GenerateResponse) error {
		response.WriteString(resp.Response)
		return nil
	})

	if err != nil {
		return "", err
	}

	return response.String(), nil
}

func (p *PostgresLLMDataAgent) executeFromLLMResponse(ctx context.Context, command *Command, llmResponse string) (interface{}, error) {
	var response struct {
		Action  string      `json:"action"`
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Error   string      `json:"error"`
	}

	if err := json.Unmarshal([]byte(llmResponse), &response); err != nil {
		return p.fallbackExecution(ctx, command)
	}

	if !response.Success {
		return nil, fmt.Errorf("LLM execution failed: %s", response.Error)
	}

	// Use the PostgreSQL agent to execute the command
	pgAgent := NewPostgresDataAgent(p.db)
	return pgAgent.ExecuteCommand(ctx, command)
}

func (p *PostgresLLMDataAgent) fallbackExecution(ctx context.Context, command *Command) (interface{}, error) {
	// Use the PostgreSQL agent to execute the command
	pgAgent := NewPostgresDataAgent(p.db)
	return pgAgent.ExecuteCommand(ctx, command)
}

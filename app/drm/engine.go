package drm

import (
	"context"
	"fmt"

	"drm-app/app/data"
	"drm-app/app/db"
)

type Engine struct {
	AuthAgent         *AuthAgent
	AccessPolicyAgent *AccessPolicyAgent
	IntentParser      *IntentParser
	LogicAgent        *LogicAgent
	DataAgent         data.DataExecutor
	Database          *db.Database
}

func NewEngine() (*Engine, error) {
	database, err := db.NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Engine{
		AuthAgent:         NewAuthAgent(),
		AccessPolicyAgent: NewAccessPolicyAgent(),
		IntentParser:      NewIntentParser(),
		LogicAgent:        NewLogicAgent(),
		DataAgent:         data.NewPostgresLLMDataAgent(database),
		Database:          database,
	}, nil
}

func (e *Engine) Close() {
	if e.Database != nil {
		e.Database.Close()
	}
}

func NewTestEngine() *Engine {
	return &Engine{
		AuthAgent:         NewAuthAgent(),
		AccessPolicyAgent: NewAccessPolicyAgent(),
		IntentParser:      NewIntentParser(),
		LogicAgent:        NewLogicAgent(),
		DataAgent:         data.NewTestDataAgent(),
		Database:          nil,
	}
}

func (e *Engine) ProcessRequest(ctx context.Context, query string, token string) (interface{}, error) {
	user, err := e.AuthAgent.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	command, err := e.IntentParser.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	command.UserID = user.ID
	command.UserRole = user.Role

	if !e.AccessPolicyAgent.CheckAccess(command) {
		return nil, fmt.Errorf("access denied for action %s on entity %s", command.Action, command.Entity)
	}

	if err := e.LogicAgent.ValidateCommand(command); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	result, err := e.DataAgent.ExecuteCommand(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	return result, nil
}

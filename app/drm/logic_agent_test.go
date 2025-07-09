package drm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"drm-app/app/data"
)

type LogicAgentTestSuite struct {
	suite.Suite
	agent *LogicAgent
}

func (s *LogicAgentTestSuite) SetupSuite() {
	s.agent = NewLogicAgent()
}

func (s *LogicAgentTestSuite) TestValidateUserCreateValid() {
	command := &data.Command{
		Action: "create",
		Entity: "user",
		Data: map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.NoError(s.T(), err)
}

func (s *LogicAgentTestSuite) TestValidateUserCreateMissingName() {
	command := &data.Command{
		Action: "create",
		Entity: "user",
		Data: map[string]interface{}{
			"email": "john@example.com",
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "user name is required")
}

func (s *LogicAgentTestSuite) TestValidateUserCreateEmptyName() {
	command := &data.Command{
		Action: "create",
		Entity: "user",
		Data: map[string]interface{}{
			"name":  "",
			"email": "john@example.com",
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "user name is required")
}

func (s *LogicAgentTestSuite) TestValidateUserCreateMissingEmail() {
	command := &data.Command{
		Action: "create",
		Entity: "user",
		Data: map[string]interface{}{
			"name": "John Doe",
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "user email is required")
}

func (s *LogicAgentTestSuite) TestValidateUserUpdateValid() {
	command := &data.Command{
		Action: "update",
		Entity: "user",
		Data: map[string]interface{}{
			"id":   "1",
			"name": "Updated Name",
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.NoError(s.T(), err)
}

func (s *LogicAgentTestSuite) TestValidateUserUpdateEmpty() {
	command := &data.Command{
		Action: "update",
		Entity: "user",
		Data:   map[string]interface{}{},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "no data provided for update")
}

func (s *LogicAgentTestSuite) TestValidateProductCreateValid() {
	command := &data.Command{
		Action: "create",
		Entity: "product",
		Data: map[string]interface{}{
			"name":  "Test Product",
			"price": 99.99,
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.NoError(s.T(), err)
}

func (s *LogicAgentTestSuite) TestValidateProductCreateMissingName() {
	command := &data.Command{
		Action: "create",
		Entity: "product",
		Data: map[string]interface{}{
			"price": 99.99,
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "product name is required")
}

func (s *LogicAgentTestSuite) TestValidateProductCreateInvalidPrice() {
	command := &data.Command{
		Action: "create",
		Entity: "product",
		Data: map[string]interface{}{
			"name":  "Test Product",
			"price": -10.0,
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "product price must be positive")
}

func (s *LogicAgentTestSuite) TestValidateProductCreateZeroPrice() {
	command := &data.Command{
		Action: "create",
		Entity: "product",
		Data: map[string]interface{}{
			"name":  "Test Product",
			"price": 0.0,
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "product price must be positive")
}

func (s *LogicAgentTestSuite) TestValidateOrderCreateValid() {
	command := &data.Command{
		Action: "create",
		Entity: "order",
		Data: map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{
					"product_id": "1",
					"quantity":   2,
				},
			},
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.NoError(s.T(), err)
}

func (s *LogicAgentTestSuite) TestValidateOrderCreateNoItems() {
	command := &data.Command{
		Action: "create",
		Entity: "order",
		Data: map[string]interface{}{
			"items": []interface{}{},
		},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "order must have at least one item")
}

func (s *LogicAgentTestSuite) TestValidateOrderCreateMissingItems() {
	command := &data.Command{
		Action: "create",
		Entity: "order",
		Data:   map[string]interface{}{},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "order must have at least one item")
}

func (s *LogicAgentTestSuite) TestValidateUnknownEntity() {
	command := &data.Command{
		Action: "create",
		Entity: "unknown",
		Data:   map[string]interface{}{},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.NoError(s.T(), err)
}

func (s *LogicAgentTestSuite) TestValidateUnknownAction() {
	command := &data.Command{
		Action: "unknown",
		Entity: "user",
		Data:   map[string]interface{}{},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.NoError(s.T(), err)
}

func (s *LogicAgentTestSuite) TestValidateReadAction() {
	command := &data.Command{
		Action: "read",
		Entity: "user",
		Data:   map[string]interface{}{},
	}
	
	err := s.agent.ValidateCommand(command)
	assert.NoError(s.T(), err)
}

func TestLogicAgentTestSuite(t *testing.T) {
	suite.Run(t, new(LogicAgentTestSuite))
}
package drm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EngineTestSuite struct {
	suite.Suite
	engine *Engine
	ctx    context.Context
}

func (s *EngineTestSuite) SetupSuite() {
	s.engine = NewTestEngine()
	s.ctx = context.Background()
}

func (s *EngineTestSuite) TearDownSuite() {
	if s.engine != nil {
		s.engine.Close()
	}
}

func (s *EngineTestSuite) TestProcessRequestValidAuth() {
	result, err := s.engine.ProcessRequest(s.ctx, "list users", "admin-token")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
}

func (s *EngineTestSuite) TestProcessRequestInvalidAuth() {
	result, err := s.engine.ProcessRequest(s.ctx, "list users", "invalid-token")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "authentication failed")
}

func (s *EngineTestSuite) TestProcessRequestParsingError() {
	result, err := s.engine.ProcessRequest(s.ctx, "invalid query", "admin-token")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "parsing failed")
}

func (s *EngineTestSuite) TestProcessRequestAccessDenied() {
	result, err := s.engine.ProcessRequest(s.ctx, "delete user json:{\"id\":\"1\"}", "guest-token")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "access denied")
}

func (s *EngineTestSuite) TestProcessRequestValidationError() {
	result, err := s.engine.ProcessRequest(s.ctx, "create user json:{\"name\":\"\"}", "admin-token")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "validation failed")
}

func (s *EngineTestSuite) TestProcessRequestSuccess() {
	result, err := s.engine.ProcessRequest(s.ctx, "create user json:{\"name\":\"Test User\",\"email\":\"test@example.com\"}", "admin-token")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)

	resultMap, ok := result.(map[string]interface{})
	assert.True(s.T(), ok)
	assert.Equal(s.T(), "test user", resultMap["name"])
	assert.Equal(s.T(), "test@example.com", resultMap["email"])
	assert.Contains(s.T(), resultMap, "id")
	assert.Contains(s.T(), resultMap, "created_at")
}

func TestEngineTestSuite(t *testing.T) {
	suite.Run(t, new(EngineTestSuite))
}

package drm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IntentParserTestSuite struct {
	suite.Suite
	parser *IntentParser
}

func (s *IntentParserTestSuite) SetupSuite() {
	s.parser = NewIntentParser()
}

func (s *IntentParserTestSuite) TestParseCreateUser() {
	command, err := s.parser.Parse("create user json:{\"name\":\"John\",\"email\":\"john@example.com\"}")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "create", command.Action)
	assert.Equal(s.T(), "user", command.Entity)
	assert.Equal(s.T(), "john", command.Data["name"])
	assert.Equal(s.T(), "john@example.com", command.Data["email"])
}

func (s *IntentParserTestSuite) TestParseReadUser() {
	command, err := s.parser.Parse("read user json:{\"id\":\"1\"}")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "read", command.Action)
	assert.Equal(s.T(), "user", command.Entity)
	assert.Equal(s.T(), "1", command.Data["id"])
}

func (s *IntentParserTestSuite) TestParseUpdateProduct() {
	command, err := s.parser.Parse("update product json:{\"id\":\"1\",\"price\":99.99}")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "update", command.Action)
	assert.Equal(s.T(), "product", command.Entity)
	assert.Equal(s.T(), "1", command.Data["id"])
	assert.Equal(s.T(), 99.99, command.Data["price"])
}

func (s *IntentParserTestSuite) TestParseDeleteOrder() {
	command, err := s.parser.Parse("delete order json:{\"id\":\"1\"}")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "delete", command.Action)
	assert.Equal(s.T(), "order", command.Entity)
	assert.Equal(s.T(), "1", command.Data["id"])
}

func (s *IntentParserTestSuite) TestParseListUsers() {
	command, err := s.parser.Parse("list users")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "read", command.Action)
	assert.Equal(s.T(), "user", command.Entity)
	assert.Empty(s.T(), command.Data)
}

func (s *IntentParserTestSuite) TestParseShowProducts() {
	command, err := s.parser.Parse("show products")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "read", command.Action)
	assert.Equal(s.T(), "product", command.Entity)
}

func (s *IntentParserTestSuite) TestParseGetOrders() {
	command, err := s.parser.Parse("get orders")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "read", command.Action)
	assert.Equal(s.T(), "order", command.Entity)
}

func (s *IntentParserTestSuite) TestParseAddUser() {
	command, err := s.parser.Parse("add user json:{\"name\":\"Jane\"}")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "create", command.Action)
	assert.Equal(s.T(), "user", command.Entity)
	assert.Equal(s.T(), "jane", command.Data["name"])
}

func (s *IntentParserTestSuite) TestParseModifyProduct() {
	command, err := s.parser.Parse("modify product json:{\"id\":\"1\"}")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "update", command.Action)
	assert.Equal(s.T(), "product", command.Entity)
}

func (s *IntentParserTestSuite) TestParseRemoveOrder() {
	command, err := s.parser.Parse("remove order json:{\"id\":\"1\"}")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "delete", command.Action)
	assert.Equal(s.T(), "order", command.Entity)
}

func (s *IntentParserTestSuite) TestParseEmptyQuery() {
	command, err := s.parser.Parse("")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), command)
	assert.Contains(s.T(), err.Error(), "empty query")
}

func (s *IntentParserTestSuite) TestParseUnknownEntity() {
	command, err := s.parser.Parse("create unknown")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), command)
	assert.Contains(s.T(), err.Error(), "unknown entity")
}

func (s *IntentParserTestSuite) TestParseDefaultAction() {
	command, err := s.parser.Parse("user")
	
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "read", command.Action)
	assert.Equal(s.T(), "user", command.Entity)
}

func (s *IntentParserTestSuite) TestParseInvalidJSON() {
	command, err := s.parser.Parse("create user json:{invalid json}")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), command)
	assert.Contains(s.T(), err.Error(), "invalid JSON data")
}

func (s *IntentParserTestSuite) TestParseWhitespaceQuery() {
	command, err := s.parser.Parse("   ")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), command)
	assert.Contains(s.T(), err.Error(), "empty query")
}

func TestIntentParserTestSuite(t *testing.T) {
	suite.Run(t, new(IntentParserTestSuite))
}
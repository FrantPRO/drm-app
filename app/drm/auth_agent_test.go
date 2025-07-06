package drm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthAgentTestSuite struct {
	suite.Suite
	agent *AuthAgent
}

func (s *AuthAgentTestSuite) SetupSuite() {
	s.agent = NewAuthAgent()
}

func (s *AuthAgentTestSuite) TestValidateValidAdminToken() {
	user, err := s.agent.ValidateToken("admin-token")
	
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), "1", user.ID)
	assert.Equal(s.T(), "Admin", user.Name)
	assert.Equal(s.T(), "admin", user.Role)
}

func (s *AuthAgentTestSuite) TestValidateValidUserToken() {
	user, err := s.agent.ValidateToken("user-token")
	
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), "2", user.ID)
	assert.Equal(s.T(), "User", user.Name)
	assert.Equal(s.T(), "user", user.Role)
}

func (s *AuthAgentTestSuite) TestValidateValidGuestToken() {
	user, err := s.agent.ValidateToken("guest-token")
	
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), "3", user.ID)
	assert.Equal(s.T(), "Guest", user.Name)
	assert.Equal(s.T(), "guest", user.Role)
}

func (s *AuthAgentTestSuite) TestValidateInvalidToken() {
	user, err := s.agent.ValidateToken("invalid-token")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
	assert.Contains(s.T(), err.Error(), "invalid token")
}

func (s *AuthAgentTestSuite) TestValidateEmptyToken() {
	user, err := s.agent.ValidateToken("")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
	assert.Contains(s.T(), err.Error(), "token is required")
}

func (s *AuthAgentTestSuite) TestValidateWhitespaceToken() {
	user, err := s.agent.ValidateToken("   ")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
	assert.Contains(s.T(), err.Error(), "token is required")
}

func (s *AuthAgentTestSuite) TestValidateTokenWithWhitespace() {
	user, err := s.agent.ValidateToken("  admin-token  ")
	
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), "admin", user.Role)
}

func (s *AuthAgentTestSuite) TestValidateNonExistentToken() {
	user, err := s.agent.ValidateToken("nonexistent-token")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
	assert.Contains(s.T(), err.Error(), "invalid token")
}

func (s *AuthAgentTestSuite) TestValidatePartialToken() {
	user, err := s.agent.ValidateToken("admin")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
	assert.Contains(s.T(), err.Error(), "invalid token")
}

func (s *AuthAgentTestSuite) TestValidateCaseSensitiveToken() {
	user, err := s.agent.ValidateToken("ADMIN-TOKEN")
	
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
	assert.Contains(s.T(), err.Error(), "invalid token")
}

func TestAuthAgentTestSuite(t *testing.T) {
	suite.Run(t, new(AuthAgentTestSuite))
}
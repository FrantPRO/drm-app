package drm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AccessPolicyAgentTestSuite struct {
	suite.Suite
	agent *AccessPolicyAgent
}

func (s *AccessPolicyAgentTestSuite) SetupSuite() {
	s.agent = NewAccessPolicyAgent()
}

func (s *AccessPolicyAgentTestSuite) TestAdminCanCreateUser() {
	command := &Command{
		Action:   "create",
		Entity:   "user",
		UserRole: "admin",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestAdminCanReadUser() {
	command := &Command{
		Action:   "read",
		Entity:   "user",
		UserRole: "admin",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestAdminCanUpdateUser() {
	command := &Command{
		Action:   "update",
		Entity:   "user",
		UserRole: "admin",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestAdminCanDeleteUser() {
	command := &Command{
		Action:   "delete",
		Entity:   "user",
		UserRole: "admin",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCanReadUser() {
	command := &Command{
		Action:   "read",
		Entity:   "user",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCanUpdateUser() {
	command := &Command{
		Action:   "update",
		Entity:   "user",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCannotCreateUser() {
	command := &Command{
		Action:   "create",
		Entity:   "user",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCannotDeleteUser() {
	command := &Command{
		Action:   "delete",
		Entity:   "user",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCanReadProduct() {
	command := &Command{
		Action:   "read",
		Entity:   "product",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCannotCreateProduct() {
	command := &Command{
		Action:   "create",
		Entity:   "product",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCanCreateOrder() {
	command := &Command{
		Action:   "create",
		Entity:   "order",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCanReadOrder() {
	command := &Command{
		Action:   "read",
		Entity:   "order",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUserCannotUpdateOrder() {
	command := &Command{
		Action:   "update",
		Entity:   "order",
		UserRole: "user",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestGuestCanReadProduct() {
	command := &Command{
		Action:   "read",
		Entity:   "product",
		UserRole: "guest",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.True(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestGuestCannotReadUser() {
	command := &Command{
		Action:   "read",
		Entity:   "user",
		UserRole: "guest",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestGuestCannotCreateProduct() {
	command := &Command{
		Action:   "create",
		Entity:   "product",
		UserRole: "guest",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestGuestCannotAccessOrder() {
	command := &Command{
		Action:   "read",
		Entity:   "order",
		UserRole: "guest",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUnknownRole() {
	command := &Command{
		Action:   "read",
		Entity:   "user",
		UserRole: "unknown",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUnknownEntity() {
	command := &Command{
		Action:   "read",
		Entity:   "unknown",
		UserRole: "admin",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func (s *AccessPolicyAgentTestSuite) TestUnknownAction() {
	command := &Command{
		Action:   "unknown",
		Entity:   "user",
		UserRole: "admin",
	}
	
	hasAccess := s.agent.CheckAccess(command)
	assert.False(s.T(), hasAccess)
}

func TestAccessPolicyAgentTestSuite(t *testing.T) {
	suite.Run(t, new(AccessPolicyAgentTestSuite))
}
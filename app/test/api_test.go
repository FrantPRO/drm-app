package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	testApp *TestApp
}

func (s *APITestSuite) SetupSuite() {
	s.testApp = NewTestApp(s.T())
}

func (s *APITestSuite) TestHealthCheck() {
	resp := s.testApp.PostRequest(TestQueries.ListUsers, AdminToken)
	obj := AssertSuccessResponse(s.T(), resp)
	obj.Value("result").Array().Length().Gt(0)
}

func (s *APITestSuite) TestValidUserList() {
	resp := s.testApp.PostRequest(TestQueries.ListUsers, AdminToken)
	obj := AssertSuccessResponse(s.T(), resp)
	obj.Value("result").Array().Length().Gt(0)
}

func (s *APITestSuite) TestCreateUser() {
	resp := s.testApp.PostRequest(TestQueries.CreateUser, AdminToken)
	obj := AssertSuccessResponse(s.T(), resp)
	
	result := obj.Value("result").Object()
	result.ContainsKey("id")
	result.ContainsKey("name")
	result.ContainsKey("email")
	result.Value("name").String().IsEqual("test user")
	result.Value("email").String().IsEqual("test@example.com")
}

func (s *APITestSuite) TestReadUser() {
	resp := s.testApp.PostRequest(TestQueries.ReadUser, AdminToken)
	obj := AssertSuccessResponse(s.T(), resp)
	
	result := obj.Value("result").Object()
	result.ContainsKey("id")
	result.ContainsKey("name")
	result.ContainsKey("email")
}

func (s *APITestSuite) TestUpdateUser() {
	resp := s.testApp.PostRequest(TestQueries.UpdateUser, AdminToken)
	obj := AssertSuccessResponse(s.T(), resp)
	
	result := obj.Value("result").Object()
	result.Value("name").String().IsEqual("updated name")
}

func (s *APITestSuite) TestDeleteUser() {
	deleteQuery := "delete user json:{\"id\":\"2\"}"
	resp := s.testApp.PostRequest(deleteQuery, AdminToken)
	obj := AssertSuccessResponse(s.T(), resp)
	
	result := obj.Value("result").Object()
	result.ContainsKey("message")
}

func (s *APITestSuite) TestCreateProduct() {
	resp := s.testApp.PostRequest(TestQueries.CreateProduct, AdminToken)
	obj := AssertSuccessResponse(s.T(), resp)
	
	result := obj.Value("result").Object()
	result.ContainsKey("id")
	result.ContainsKey("name")
	result.ContainsKey("price")
	result.Value("name").String().IsEqual("test product")
	result.Value("price").Number().IsEqual(99.99)
}

func (s *APITestSuite) TestListProducts() {
	resp := s.testApp.PostRequest(TestQueries.ListProducts, UserToken)
	obj := AssertSuccessResponse(s.T(), resp)
	obj.Value("result").Array().Length().Gt(0)
}

func (s *APITestSuite) TestCreateOrder() {
	resp := s.testApp.PostRequest(TestQueries.CreateOrder, UserToken)
	obj := AssertSuccessResponse(s.T(), resp)
	
	result := obj.Value("result").Object()
	result.ContainsKey("id")
	result.ContainsKey("items")
}

func (s *APITestSuite) TestInvalidToken() {
	resp := s.testApp.PostRequest(TestQueries.ListUsers, "invalid-token")
	AssertAuthError(s.T(), resp)
}

func (s *APITestSuite) TestMissingToken() {
	resp := s.testApp.PostRequestWithoutToken(TestQueries.ListUsers)
	AssertBadRequestError(s.T(), resp, "Token is required")
}

func (s *APITestSuite) TestMissingQuery() {
	resp := s.testApp.PostRequestWithoutQuery(AdminToken)
	AssertBadRequestError(s.T(), resp, "Query is required")
}

func (s *APITestSuite) TestInvalidQuery() {
	resp := s.testApp.PostRequest(TestQueries.InvalidQuery, AdminToken)
	AssertParsingError(s.T(), resp)
}

func (s *APITestSuite) TestAccessDenied() {
	resp := s.testApp.PostRequest(TestQueries.DeleteUser, GuestToken)
	AssertAccessDeniedError(s.T(), resp)
}

func (s *APITestSuite) TestValidationError() {
	resp := s.testApp.PostRequest(TestQueries.CreateUserNoName, AdminToken)
	AssertValidationError(s.T(), resp)
}

func (s *APITestSuite) TestInvalidJSON() {
	resp := s.testApp.PostInvalidJSON()
	AssertBadRequestError(s.T(), resp, "Invalid request body")
}

func (s *APITestSuite) TestUserRolePermissions() {
	resp := s.testApp.PostRequest(TestQueries.CreateUser, UserToken)
	AssertAccessDeniedError(s.T(), resp)
}

func (s *APITestSuite) TestGuestRolePermissions() {
	resp := s.testApp.PostRequest(TestQueries.ListProducts, GuestToken)
	AssertSuccessResponse(s.T(), resp)
}

func (s *APITestSuite) TestUserCanCreateOrder() {
	resp := s.testApp.PostRequest(TestQueries.CreateOrder, UserToken)
	AssertSuccessResponse(s.T(), resp)
}

func (s *APITestSuite) TestUserCannotDeleteUser() {
	resp := s.testApp.PostRequest(TestQueries.DeleteUser, UserToken)
	AssertAccessDeniedError(s.T(), resp)
}

func (s *APITestSuite) TestGuestCannotCreateProduct() {
	resp := s.testApp.PostRequest(TestQueries.CreateProduct, GuestToken)
	AssertAccessDeniedError(s.T(), resp)
}

func (s *APITestSuite) TestGuestCannotAccessUsers() {
	resp := s.testApp.PostRequest(TestQueries.ListUsers, GuestToken)
	AssertAccessDeniedError(s.T(), resp)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
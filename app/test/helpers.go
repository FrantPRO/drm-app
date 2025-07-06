package test

import (
	"net/http"
	"testing"

	"drm-app/app/drm"
	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
)

type TestApp struct {
	App    *fiber.App
	Engine *drm.Engine
	Client *httpexpect.Expect
}

func NewTestApp(t *testing.T) *TestApp {
	engine := drm.NewEngine()
	
	app := fiber.New(fiber.Config{
		AppName: "DRM Core Test v1.0.0",
	})
	
	app.Post("/request", createHandleRequest(engine))
	
	client := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(app.Handler()),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})
	
	return &TestApp{
		App:    app,
		Engine: engine,
		Client: client,
	}
}

func createHandleRequest(engine *drm.Engine) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handleRequest(c, engine)
	}
}

type RequestBody struct {
	Query string `json:"query"`
	Token string `json:"token"`
}

func handleRequest(c *fiber.Ctx, engine *drm.Engine) error {
	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query is required",
		})
	}

	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is required",
		})
	}

	result, err := engine.ProcessRequest(c.Context(), req.Query, req.Token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
		"status": "success",
	})
}

func (ta *TestApp) PostRequest(query, token string) *httpexpect.Response {
	return ta.Client.POST("/request").
		WithJSON(map[string]string{
			"query": query,
			"token": token,
		}).
		Expect()
}

func (ta *TestApp) PostRequestWithoutToken(query string) *httpexpect.Response {
	return ta.Client.POST("/request").
		WithJSON(map[string]string{
			"query": query,
		}).
		Expect()
}

func (ta *TestApp) PostRequestWithoutQuery(token string) *httpexpect.Response {
	return ta.Client.POST("/request").
		WithJSON(map[string]string{
			"token": token,
		}).
		Expect()
}

func (ta *TestApp) PostInvalidJSON() *httpexpect.Response {
	return ta.Client.POST("/request").
		WithText("invalid json").
		Expect()
}

func AssertSuccessResponse(t *testing.T, resp *httpexpect.Response) *httpexpect.Object {
	if resp.Raw().StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response body: %s", resp.Raw().StatusCode, resp.Body().Raw())
		t.FailNow()
	}
	
	obj := resp.Status(http.StatusOK).
		JSON().
		Object()
	
	obj.Value("status").String().IsEqual("success")
	obj.ContainsKey("result")
	
	return obj
}

func AssertErrorResponse(t *testing.T, resp *httpexpect.Response, statusCode int, errorMsg string) {
	obj := resp.Status(statusCode).
		JSON().
		Object()
	
	obj.ContainsKey("error")
	if errorMsg != "" {
		obj.Value("error").String().Contains(errorMsg)
	}
}

func AssertValidationError(t *testing.T, resp *httpexpect.Response) {
	AssertErrorResponse(t, resp, http.StatusInternalServerError, "validation failed")
}

func AssertAuthError(t *testing.T, resp *httpexpect.Response) {
	AssertErrorResponse(t, resp, http.StatusInternalServerError, "authentication failed")
}

func AssertAccessDeniedError(t *testing.T, resp *httpexpect.Response) {
	AssertErrorResponse(t, resp, http.StatusInternalServerError, "access denied")
}

func AssertParsingError(t *testing.T, resp *httpexpect.Response) {
	AssertErrorResponse(t, resp, http.StatusInternalServerError, "parsing failed")
}

func AssertBadRequestError(t *testing.T, resp *httpexpect.Response, errorMsg string) {
	AssertErrorResponse(t, resp, http.StatusBadRequest, errorMsg)
}

const (
	AdminToken = "admin-token"
	UserToken  = "user-token"
	GuestToken = "guest-token"
)

var TestQueries = struct {
	ListUsers        string
	ListProducts     string
	ListOrders       string
	CreateUser       string
	CreateProduct    string
	CreateOrder      string
	ReadUser         string
	UpdateUser       string
	DeleteUser       string
	InvalidQuery     string
	CreateUserNoName string
}{
	ListUsers:        "list users",
	ListProducts:     "list products",
	ListOrders:       "list orders",
	CreateUser:       "create user json:{\"name\":\"Test User\",\"email\":\"test@example.com\"}",
	CreateProduct:    "create product json:{\"name\":\"Test Product\",\"price\":99.99}",
	CreateOrder:      "create order json:{\"items\":[{\"product_id\":\"1\",\"quantity\":2}]}",
	ReadUser:         "read user json:{\"id\":\"1\"}",
	UpdateUser:       "update user json:{\"id\":\"1\",\"name\":\"Updated Name\"}",
	DeleteUser:       "delete user json:{\"id\":\"1\"}",
	InvalidQuery:     "invalid query without entity",
	CreateUserNoName: "create user json:{\"name\":\"\"}",
}
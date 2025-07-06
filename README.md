# DRM (Declarative-Relation Mapping Engine) Application

## Overview
###### DRM (Declarative-Relation Mapping) is an experimental backend architecture built on the principle of declarative logic execution, where intents expressed in natural language are interpreted and executed by a central reasoning engine. Instead of exposing dozens of static API endpoints, DRM provides a single /request interface that dynamically processes user instructions using structured entity declarations and agent-based logic.

#### This MVP demonstrates:
* Agent-based architecture with declarative schemas
* A unified natural-language interface (via /request)
* Secure role-based access control (RBAC)
* Executable business rules described as YAML policies
* Integration-ready for LLM-based parsing and reasoning

## Install dependencies
```console
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1
```

## Architecture
#### System Components
 
| Layer          | Component           | Responsibility                           |
|:---------------|:--------------------|:-----------------------------------------|
| API            | `/request` endpoint | Accepts user queries in natural language |
| Auth           | `AuthAgent`         | Validates tokens and extracts user roles |
| Access Control | `AccessPolicyAgent` | Enforces entity-level access rules       |
| Parsing        | `IntentParser`      | Converts user query → structured Command |
| Logic          | `LogicAgent`        | Validates data against YAML rules        |
| Data           | `DataAgent`         | Executes Create/Read/Update/Delete       |
| Storage        | `PostgreSQL`        | Persistent data store                    |

#### Directory Structure
```
drm-app/
├── app/
│   ├── main.go                  # entry point
│   ├── drm/                     # core DRM engine
│   ├── handlers/                # HTTP request routing
│   ├── parser/                  # intent parsing logic
│   ├── logic/                   # rule execution and validation
│   ├── access/                  # access policies
│   ├── data/                    # data storage abstraction
│   ├── schemas/                 # YAML-based entity declarations
│   └── utils/                   # helper functions
├── db-data/                     # volume-mounted folder for Postgres
├── docker/
│   ├── Dockerfile               # Go backend container build
│   └── init.sql (optional)      # schema bootstrap for Postgres
├── docker-compose.yml           # full service orchestration
├── go.mod / go.sum              # Go module definition
└── README.md
```

#### Technologies & Libraries
| Purpose               | Tool / Library         | Version | Notes                         |
|:----------------------|:-----------------------|:-------:|:------------------------------|
| Language              | Go                     | 1.24.4  | Latest stable version         |
| Web framework         | Fiber                  | v2.52.8 | Fast and minimal Express-like |
| YAML parser           | yaml.v3                | v3.0.1  | For entity/rule declarations  |
| PostgreSQL driver     | pgx                    | v5.7.5  | Native driver for PostgreSQL  |
| SQL helper            | sqlx                   | v1.4.0  | Lightweight ORM-less access   |
| HTTP client (LLM)     | resty                  | v2.16.5 | Optional: LLM API integration |
| Testing framework     | testify                | v1.10.0 | Unit/integration testing      |
| OpenAI API (optional) | go-openai              | v1.40.3 | GPT-4 support                 |
| Database              | PostgreSQL             |   17+   | Local or Docker volume        |
| Containerization      | Docker, Docker Compose | latest  | Two containers: app + db      |


### MVP Features
* POST /request — accepts natural-language instructions
* Built-in auth (token → user + role)
* Schema-defined entity structure and rules
* Flexible YAML-based entity declarations (schemas/*.yaml)
* Role-based access policies (access/*.yaml)
* Modular agent architecture (Parser, Access, Logic, Data)
* Easy future integration with LLMs (GPT/OpenRouter/Ollama)

## API Usage

### Authentication
The application uses token-based authentication with three predefined roles:

| Token | Role | Permissions |
|-------|------|-------------|
| `admin-token` | Admin | Full access: create, read, update, delete all entities |
| `user-token` | User | Limited: read/update users, read products, create/read orders |
| `guest-token` | Guest | Read-only: products only |

### Endpoint
**POST** `/request`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "query": "natural language query with optional json:{...}",
  "token": "admin-token|user-token|guest-token"
}
```

### Request Examples

#### Admin Role Examples (admin-token)
Full access to all operations:

```bash
# Create a new user
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "create user json:{\"name\":\"John Doe\",\"email\":\"john@example.com\"}", "token": "admin-token"}'

# List all users
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "list users", "token": "admin-token"}'

# Update a user
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "update user json:{\"id\":\"1\",\"name\":\"Jane Smith\"}", "token": "admin-token"}'

# Delete a user
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "delete user json:{\"id\":\"2\"}", "token": "admin-token"}'

# Create a product
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "create product json:{\"name\":\"Gaming Laptop\",\"price\":1299.99}", "token": "admin-token"}'
```

#### User Role Examples (user-token)
Limited access - can manage own profile and orders:

```bash
# Read user information
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "read user json:{\"id\":\"2\"}", "token": "user-token"}'

# Update own profile
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "update user json:{\"id\":\"2\",\"name\":\"Updated Name\"}", "token": "user-token"}'

# View products
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "list products", "token": "user-token"}'

# Create an order
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "create order json:{\"items\":[{\"product_id\":\"1\",\"quantity\":2}]}", "token": "user-token"}'

# View orders
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "list orders", "token": "user-token"}'
```

#### Guest Role Examples (guest-token)
Read-only access to products:

```bash
# View products (allowed)
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "list products", "token": "guest-token"}'

# Read specific product
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "read product json:{\"id\":\"1\"}", "token": "guest-token"}'

# Access denied examples:
# Trying to access users (will fail)
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"query": "list users", "token": "guest-token"}'
# Response: {"error":"access denied for action read on entity user"}
```

### Response Format

**Success Response:**
```json
{
  "result": {
    "id": "1",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-01-06T12:00:00Z"
  },
  "status": "success"
}
```

**Error Response:**
```json
{
  "error": "authentication failed: invalid token"
}
```

### Natural Language Query Format

The query format supports:
- **Actions**: create, add, read, get, list, show, update, modify, change, delete, remove
- **Entities**: user, product, order
- **Data**: `json:{...}` for structured data

**Examples:**
- `"list users"` - Read all users
- `"create user json:{\"name\":\"...\",\"email\":\"...\"}"`
- `"update product json:{\"id\":\"1\",\"price\":99.99}"`
- `"delete order json:{\"id\":\"1\"}"`

### Planned Extensions
* Add LLM-based IntentParser using OpenAI or local models
* Event-based agent (EventAgent)
* Admin Web UI (React)
* Exportable history log + audit
* Plugin-style agent registration
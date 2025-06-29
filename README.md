# DRM (Declarative-Relation Mapping Engine) Application

## Overview
###### DRM (Declarative-Relation Mapping) is an experimental backend architecture built on the principle of declarative logic execution, where intents expressed in natural language are interpreted and executed by a central reasoning engine. Instead of exposing dozens of static API endpoints, DRM provides a single /request interface that dynamically processes user instructions using structured entity declarations and agent-based logic.

#### This MVP demonstrates:
* Agent-based architecture with declarative schemas
* A unified natural-language interface (via /request)
* Secure role-based access control (RBAC)
* Executable business rules described as YAML policies
* Integration-ready for LLM-based parsing and reasoning

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

### Planned Extensions
* Add LLM-based IntentParser using OpenAI or local models
* Event-based agent (EventAgent)
* Admin Web UI (React)
* Exportable history log + audit
* Plugin-style agent registration
# Project Structure Guide

This document explains how the UniEntrega backend is organized.

## 📁 Directory Structure

```
unientrega/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handlers/
│   │   ├── health.go            # Health & root route handlers
│   │   └── api.go               # API route handlers
│   ├── routes/
│   │   └── routes.go            # Route definitions
│   ├── middleware/              # (future) Custom middleware
│   ├── models/                  # (future) Data models
│   ├── services/                # (future) Business logic
│   └── repository/              # (future) Database access
├── docker/
│   ├── api.Dockerfile           # Production Dockerfile
│   ├── api.dev.Dockerfile       # Development Dockerfile
│   ├── docker-compose.yml       # Production compose
│   └── docker-compose.dev.yml   # Development compose
├── .env                         # Environment variables (local)
├── .env.example                 # Environment template
├── .air.toml                    # Hot reload configuration
└── go.mod                       # Go dependencies
```

## 🏗️ Architecture Layers

### 1. **cmd/server/main.go** - Entry Point
- Loads configuration
- Initializes Gin router
- Sets up routes
- Starts HTTP server

**Responsibilities:**
- Application bootstrap
- Server configuration
- Graceful startup/shutdown

### 2. **internal/routes/** - Route Definitions
- Defines all HTTP routes
- Groups related routes
- Maps routes to handlers

**Example:**
```go
r.GET("/", healthHandler.GetRoot)
r.GET("/health", healthHandler.GetHealth)
```

### 3. **internal/handlers/** - Request Handlers
- Handles HTTP requests
- Validates input
- Calls services
- Returns responses

**Structure:**
```go
type HealthHandler struct {
    config *config.Config
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
}
```

### 4. **internal/config/** - Configuration
- Loads environment variables
- Provides type-safe config
- Manages .env files

### 5. **internal/middleware/** (future)
- Authentication
- Logging
- CORS
- Rate limiting

### 6. **internal/models/** (future)
- Data structures
- Database models
- DTOs (Data Transfer Objects)

### 7. **internal/services/** (future)
- Business logic
- Data processing
- External API calls

### 8. **internal/repository/** (future)
- Database queries
- Data access layer
- CRUD operations

## 🔄 Request Flow

```
HTTP Request
    ↓
main.go (Gin Router)
    ↓
routes/routes.go (Route Matching)
    ↓
handlers/*.go (Request Handler)
    ↓
services/*.go (Business Logic) [future]
    ↓
repository/*.go (Database) [future]
    ↓
Response
```

## ✅ Benefits of This Structure

### **Separation of Concerns**
- Each layer has a single responsibility
- Easy to understand and maintain
- Changes in one layer don't affect others

### **Testability**
- Handlers can be tested independently
- Services can be mocked
- Clear dependencies

### **Scalability**
- Easy to add new routes
- Simple to add new handlers
- Clear where new code belongs

### **Team Collaboration**
- Multiple developers can work on different layers
- Reduces merge conflicts
- Standard structure everyone understands

## 📝 Adding New Features

### Adding a New Route

1. **Create handler** in `internal/handlers/`:
```go
// internal/handlers/user.go
type UserHandler struct {
    config *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
    return &UserHandler{config: cfg}
}

func (h *UserHandler) GetUser(c *gin.Context) {
    c.JSON(200, gin.H{"user": "example"})
}
```

2. **Register route** in `internal/routes/routes.go`:
```go
userHandler := handlers.NewUserHandler(cfg)
v1.GET("/users/:id", userHandler.GetUser)
```

3. **Done!** The route is now available at `/api/v1/users/:id`

## 🎯 Best Practices

### **DO:**
- ✅ Keep handlers thin (just HTTP logic)
- ✅ Put business logic in services
- ✅ Use dependency injection
- ✅ Return proper HTTP status codes
- ✅ Validate input in handlers

### **DON'T:**
- ❌ Put business logic in handlers
- ❌ Access database directly from handlers
- ❌ Mix concerns across layers
- ❌ Use global variables
- ❌ Hardcode configuration

## 📚 Further Reading

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [Clean Architecture in Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

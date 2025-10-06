# GateHide API

A RESTful API built with Go and Gin framework following SOLID principles and clean architecture.

## 🚀 Features

- ✅ Clean architecture with separation of concerns
- ✅ SOLID principles implementation
- ✅ Health check endpoint
- ✅ Custom logging middleware
- ✅ CORS support
- ✅ Security headers
- ✅ Environment-based configuration
- ✅ Graceful error handling
- ✅ Hot reloading for development

## 📁 Project Structure

```
gatehide-api/
├── cmd/
│   └── app/
│       └── main.go           # Application entry point
├── config/
│   └── config.go             # Configuration management
├── internal/
│   ├── handlers/             # HTTP request handlers
│   │   └── health_handler.go
│   ├── middlewares/          # Custom middlewares
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── security.go
│   ├── models/               # Data models
│   │   └── health.go
│   └── routes/               # Route definitions
│       └── routes.go
├── .env                      # Environment variables
├── .env.example              # Example environment variables
├── .air.toml                 # Air hot reload configuration
├── .gitignore
├── go.mod                    # Go module definition
├── README.md                 # Full documentation
├── QUICKSTART.md             # Quick start guide
└── HOTRELOAD.md              # Hot reloading guide
```

## 🛠️ Prerequisites

- Go 1.22 or higher
- Git

## 📦 Installation

1. Clone the repository (if not already done):
```bash
cd /Applications/XAMPP/xamppfiles/htdocs/gatehide/gatehide-api
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the example environment file:
```bash
cp .env.example .env
```

4. Update the `.env` file with your configuration.

## 🏃 Running the Application

### Development Mode with Hot Reload

```bash
# Using Air for hot reloading (recommended for development)
make dev
# or
make hot
# or
air
```

### Standard Development Mode

```bash
go run cmd/app/main.go
```

### Build and Run

```bash
# Build the binary
go build -o bin/gatehide-api cmd/app/main.go

# Run the binary
./bin/gatehide-api
```

## 🧪 Testing the Health Endpoint

Once the server is running, you can test the health endpoint:

### Using curl
```bash
curl http://localhost:8080/health
```

### Using browser
Open your browser and navigate to:
```
http://localhost:8080/health
```

### Expected Response
```json
{
  "status": "healthy",
  "timestamp": "2025-10-06T10:30:00Z",
  "service": "GateHide API",
  "version": "1.0.0"
}
```

## 📚 API Endpoints

### Health Check

**Endpoint:** `GET /health` or `GET /api/v1/health`

**Description:** Check if the API is running and healthy

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-06T10:30:00Z",
  "service": "GateHide API",
  "version": "1.0.0"
}
```

## 🔧 Configuration

The application can be configured using environment variables in the `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| GIN_MODE | Gin mode (debug/release) | debug |
| APP_NAME | Application name | GateHide API |
| APP_VERSION | Application version | 1.0.0 |
| API_SECRET | API secret key | - |

## 🏗️ Architecture Principles

This project follows:

- **SOLID Principles**: Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, Dependency Inversion
- **Clean Code**: Meaningful names, single responsibility functions, explicit error handling
- **Security First**: Input validation, secure headers, no sensitive data logging
- **DRY**: Reusable code through proper abstraction
- **KISS**: Simple, maintainable solutions
- **Modular Architecture**: Clear separation of layers (handlers, models, routes, middlewares)

## 📝 Code Quality

### Formatting
```bash
go fmt ./...
```

### Linting
```bash
golangci-lint run
```

### Testing
```bash
go test ./...
```

## 🔒 Security Features

- Security HTTP headers (X-Content-Type-Options, X-Frame-Options, etc.)
- CORS configuration
- Request logging without sensitive data
- Environment-based secrets management

## 🤝 Contributing

1. Follow the SOLID principles and clean code practices
2. Ensure all tests pass
3. Update documentation as needed
4. Use meaningful commit messages

## 📄 License

[Add your license here]

## 👥 Authors

GateHide Team


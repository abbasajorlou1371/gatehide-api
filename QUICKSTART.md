# 🚀 Quick Start Guide

## Step 1: Install Dependencies

```bash
cd /Applications/XAMPP/xamppfiles/htdocs/gatehide/gatehide-api
go mod download
```

## Step 2: Run the Server

Choose one of the following methods:

### Method 1: Hot Reload Development (Recommended)
```bash
make dev
# or
make hot
# or
air
```

### Method 2: Standard Development
```bash
go run cmd/app/main.go
```

### Method 3: Using Make
```bash
make run
```

### Method 4: Build and Run Binary
```bash
make build
./bin/gatehide-api
```

## Step 3: Test the Health Endpoint

Once the server is running, open a new terminal and test:

### Using the Test Script
```bash
./test_health.sh
```

### Using curl
```bash
curl http://localhost:8080/health
```

### Using Browser
Open: http://localhost:8080/health

## Expected Output

```json
{
  "status": "healthy",
  "timestamp": "2025-10-06T10:30:00Z",
  "service": "GateHide API",
  "version": "1.0.0"
}
```

## Server Logs

When the server starts successfully, you should see:

```
🚀 Starting GateHide API v1.0.0
📡 Server running on port 8080
🔧 Environment: debug
🏥 Health check available at: http://localhost:8080/health
```

## Available Make Commands

```bash
make help       # Show all available commands
make install     # Install dependencies
make run         # Run the application
make dev         # Run with hot reload (recommended for development)
make hot         # Alias for dev command
make build       # Build the application
make test        # Run tests
make clean       # Clean build artifacts
make fmt         # Format code
make lint        # Run linter (requires golangci-lint)
```

## Troubleshooting

### Port Already in Use
If port 8080 is already in use, change the port in `.env`:
```
PORT=3000
```

### Missing Dependencies
Run:
```bash
go mod tidy
```

## Hot Reloading Features

🔥 **Hot Reloading** is now enabled! When you make changes to your Go files:

- ✅ **Automatic Restart**: Server restarts automatically on file changes
- ✅ **Fast Build**: Only rebuilds when necessary
- ✅ **File Watching**: Monitors `.go`, `.html`, `.env` files
- ✅ **Clean Output**: Clear build logs and error messages
- ✅ **Exclude Patterns**: Ignores test files and vendor directories

### Hot Reload Commands

```bash
# Start with hot reload (recommended for development)
make dev
make hot
air

# Stop hot reload: Ctrl+C
```

## Next Steps

1. ✅ Health endpoint is working
2. ✅ Hot reloading is configured
3. Add database connection
4. Add authentication middleware
5. Add business logic endpoints
6. Add tests

## Project Structure Overview

```
gatehide-api/
├── cmd/app/main.go              # Entry point
├── config/config.go             # Configuration
├── internal/
│   ├── handlers/                # Request handlers
│   ├── middlewares/             # HTTP middlewares
│   ├── models/                  # Data models
│   └── routes/                  # Route definitions
├── .env                         # Environment variables
├── Makefile                     # Build commands
└── README.md                    # Full documentation
```

## Support

For full documentation, see [README.md](README.md)


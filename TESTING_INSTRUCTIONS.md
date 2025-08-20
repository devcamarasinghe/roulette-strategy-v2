# Roulette Strategy Platform - Testing Instructions

## 🚀 Quick Start

### Prerequisites (REQUIRED!)
1. **Docker Desktop ONLY** - Download and install:
   - Windows: https://docs.docker.com/desktop/install/windows-install/  
   - Mac: https://docs.docker.com/desktop/install/mac-install/
   - Linux: https://docs.docker.com/desktop/install/linux-install/
2. **Start Docker Desktop** and wait for green/blue whale icon
3. **No need to install**: Python, Go, Node.js, or any programming languages! Everything runs in containers 🐳
4. **Test Docker works:**
   ```bash
   docker --version
   docker-compose --version
   ```

### Steps to Test

1. **Get the project files**
   ```bash
   # If from GitHub
   git clone <your-repo-url>
   cd roulette-strategy-v2/backend
   
   # If from files
   # Extract the project folder and navigate to backend/
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Test the application**
   - Main Gateway: http://localhost:8080/
   - API Documentation: http://localhost:8080/users/docs
   - Register User: POST http://localhost:8080/users/register
   - Login User: POST http://localhost:8080/users/login

4. **Stop all services**
   ```bash
   docker-compose down
   ```

## 🧪 Testing Endpoints

### Register New User
```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"testpass123"}'
```

### Login User  
```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}'
```

### List Users
```bash
curl http://localhost:8080/users/users
```

## 🐛 Troubleshooting

- **Port conflicts**: Make sure ports 8080 and 8001 are free
- **Docker issues**: Restart Docker Desktop
- **Network errors**: Run `docker-compose down` then `docker-compose up -d`

## 📋 Architecture

```
🌐 Your Browser
    ↓
🚂 API Gateway (Port 8080) - Routes requests
    ↓  
🏢 User Service (Port 8001) - Handles users
```

Both services run in Docker containers with automatic networking!

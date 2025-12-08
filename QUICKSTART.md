# 🚀 Quick Start Guide - bdSeeker Backend

Get the bdSeeker API running in **5 minutes**!

## Prerequisites Checklist

- [ ] Go 1.22+ installed (`go version`)
- [ ] PostgreSQL installed (`psql --version`)
- [ ] PostgreSQL service running

## 5-Minute Setup

### 1️⃣ Get the Code (30 seconds)

```bash
git clone https://github.com/bishworup11/bdSeeker-backend.git
cd bdSeeker-backend
```

### 2️⃣ Install Dependencies (1 minute)

```bash
go mod download
```

### 3️⃣ Setup Database (1 minute)

```bash
# Create database
psql -U postgres -c "CREATE DATABASE bdseeker;"
# Enter your PostgreSQL password when prompted
```

### 4️⃣ Configure Environment (30 seconds)

The application uses **Viper** for configuration. You have three options:

**Option A: Use existing .env file (easiest)**
```bash
# The .env file is already configured with defaults
# If your PostgreSQL password is different from `postgres`, edit it:
nano .env
```

**Option B: Set environment variables**
```bash
export DB_PASSWORD=your_password
export JWT_SECRET=your-secret-key
```

**Option C: Create config.yaml (optional)**
```yaml
dbhost: localhost
dbport: "5432"
dbuser: postgres
dbpassword: your_password
dbname: bdseeker
```

### 5️⃣ Run the Server (30 seconds)

```bash
go run main.go
```

You should see:
```
2025/XX/XX XX:XX:XX No config file found, using environment variables and defaults
✓ Database connection established successfully
✓ Database migrations completed successfully
🚀 Server starting on 0.0.0.0:9000
📚 API Documentation: http://0.0.0.0:9000/api/v1/health
🔧 Environment: development
```

### 6️⃣ Test It! (1 minute)

```bash
# Test health check
curl http://localhost:9000/api/v1/health

# Should return: {"status":"ok"}

# Run full test suite
chmod +x test_api.sh
./test_api.sh
```

## ✅ Success!

If all 26 tests pass, you're ready to go! 🎉

## 🎯 Next Steps

1. **Read the API docs**: Check out [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
2. **Try the endpoints**: Use Postman or curl to interact with the API
3. **Build your frontend**: Connect your React/Vue/Angular app to the API

## 📝 Quick API Examples

### Register a User
```bash
curl -X POST http://localhost:9000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User",
    "role": "developer"
  }'
```

### Login
```bash
curl -X POST http://localhost:9000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### List Jobs
```bash
curl http://localhost:9000/api/v1/jobs
```

## 🆘 Having Issues?

### Database Connection Error?
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Start it if needed
sudo systemctl start postgresql
```

### Port Already in Use?
```bash
# Change port in .env or via environment variable
export SERVER_PORT=8081
```

### Configuration Issues?
Viper reads configuration in this priority order:
1. Environment variables (highest priority)
2. config.yaml/config.json file
3. .env file variables
4. Default values (lowest priority)

### Need More Help?
Check the full [README.md](README.md) for detailed troubleshooting.

---

**That's it! You're all set! 🚀**

## ⚡ Performance Upgrade

This backend now uses **Gin framework** (40x faster than Gorilla Mux) and **Viper** for robust configuration management!

# 🚀 bdSeeker Backend - Job Searching Platform API

A comprehensive REST API backend for a job searching platform built with **Go, GORM, and PostgreSQL**. Similar to HackerRank Jobs, this platform connects developers with companies and provides features for job posting, company reviews, ratings, and more.

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Running the Application](#-running-the-application)
- [Testing the API](#-testing-the-api)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Troubleshooting](#-troubleshooting)

## ✨ Features

### Core Features
- 🔐 **JWT Authentication** - Secure user registration and login
- 👥 **Role-based Access Control** - Developer, Company, and Admin roles
- 🏢 **Company Profiles** - Company information, ratings, and reviews
- 💼 **Developer Profiles** - Bio, experience, education, and certificates
- 📝 **Job Posting** - Create, search, and filter job opportunities
- ⭐ **Ratings & Reviews** - Rate and review companies
- 💬 **Comments & Reactions** - Engage with job posts
- 🔍 **Advanced Search** - Filter by location, work mode, salary, experience
- 📄 **Pagination** - Efficient data handling for large datasets

### Technical Features
- Clean architecture with separation of concerns
- GORM for database operations with auto-migration
- Input validation using go-playground/validator
- Standardized API responses
- CORS support for frontend integration
- Comprehensive error handling

## 🛠 Tech Stack

- **Language**: Go 1.22+
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Router**: Gin Web Framework
- **Configuration**: Viper
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator

## 📦 Prerequisites

Before you begin, ensure you have the following installed on your system:

### 1. Go Programming Language
- **Version**: 1.22 or higher
- **Download**: https://golang.org/dl/
- **Verify installation**:
  ```bash
  go version
  ```
  You should see something like: `go version go1.22.4 linux/amd64`

### 2. PostgreSQL Database
- **Version**: 12 or higher
- **Download**: https://www.postgresql.org/download/
- **Verify installation**:
  ```bash
  psql --version
  ```
  You should see something like: `psql (PostgreSQL) 14.x`

### 3. Git (Optional but recommended)
- **Download**: https://git-scm.com/downloads
- **Verify installation**:
  ```bash
  git --version
  ```

## 🚀 Installation

### Step 1: Clone the Repository

```bash
# Using Git
git clone https://github.com/bishworup11/bdSeeker-backend.git
cd bdSeeker-backend

# OR download and extract the ZIP file, then navigate to the folder
cd bdSeeker-backend
```

### Step 2: Install Go Dependencies

```bash
# Download all required packages
go mod download

# Verify dependencies are installed
go mod tidy
```

This will install:
- GORM and PostgreSQL driver
- JWT library
- Gin web framework
- Viper configuration management
- Validator
- Bcrypt for password hashing
- And other dependencies

### Step 3: Set Up PostgreSQL Database

#### Option A: Using psql Command Line

1. **Start PostgreSQL service** (if not already running):
   ```bash
   # On Linux
   sudo systemctl start postgresql
   
   # On macOS (with Homebrew)
   brew services start postgresql
   
   # On Windows
   # PostgreSQL should start automatically, or use Services app
   ```

2. **Create the database**:
   ```bash
   # Connect to PostgreSQL
   psql -U postgres
   
   # You'll be prompted for password (default is often 'postgres')
   # Then run this SQL command:
   CREATE DATABASE bdseeker;
   
   # Exit psql
   \q
   ```

#### Option B: Using pgAdmin (GUI)

1. Open pgAdmin
2. Connect to your PostgreSQL server
3. Right-click on "Databases" → "Create" → "Database"
4. Enter database name: `bdseeker`
5. Click "Save"

### Step 4: Configure Environment Variables

1. **Copy the example environment file**:
   ```bash
   # The .env file should already exist, but if not:
   cp .env.example .env
   ```

2. **Edit the `.env` file** with your settings:
   ```bash
   # Open with your favorite text editor
   nano .env
   # or
   vim .env
   # or use any text editor
   ```

3. **Update the following values**:
   
   **Option A: Using Environment Variables (Recommended with Viper)**
   ```bash
   # Set environment variables directly
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=your_password
   export DB_NAME=bdseeker
   export DB_SSLMODE=disable
   export JWT_SECRET=your-super-secret-jwt-key
   export JWT_EXPIRY=24h
   export JWT_REFRESH_EXPIRY=168h
   export SERVER_PORT=9000
   export SERVER_HOST=0.0.0.0
   export ENV=development
   ```
   
   **Option B: Using .env file (Backward Compatible)**
   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres          # Your PostgreSQL username
   DB_PASSWORD=postgres      # Your PostgreSQL password
   DB_NAME=bdseeker
   DB_SSLMODE=disable

   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   JWT_EXPIRY=24h
   JWT_REFRESH_EXPIRY=168h

   # Server Configuration
   SERVER_PORT=9000
   SERVER_HOST=0.0.0.0

   # Environment
   ENV=development
   ```
   
   **Option C: Using config.yaml (New with Viper)**
   Create a `config.yaml` file:
   ```yaml
   dbhost: localhost
   dbport: "5432"
   dbuser: postgres
   dbpassword: postgres
   dbname: bdseeker
   dbsslmode: disable
   jwtsecret: your-super-secret-jwt-key
   jwtexpiry: 24h
   jwtrefreshexpiry: 168h
   serverport: "9000"
   serverhost: 0.0.0.0
   environment: development
   ```

   **Important**: 
   - Viper supports multiple configuration sources with priority: ENV vars > config file > defaults
   - Environment variables take highest priority
   - Replace `DB_USER` and `DB_PASSWORD` with your PostgreSQL credentials
   - Change `JWT_SECRET` to a strong, random string in production

## 🏃 Running the Application

### Method 1: Using `go run` (Development)

```bash
# Run directly without building
go run main.go
```

You should see:
```
✓ Database connection established successfully
✓ Database migrations completed successfully
🚀 Server starting on 0.0.0.0:9000
📚 API Documentation: http://0.0.0.0:9000/api/v1/health
```

### Method 2: Build and Run (Production-like)

```bash
# Build the executable
go build -o bdseeker-api main.go

# Run the executable
./bdseeker-api
```

### Verify the Server is Running

Open your browser or use curl:
```bash
curl http://localhost:9000/api/v1/health
```

You should see:
```json
{"status":"ok"}
```

## 🧪 Testing the API

### Automated Testing

We provide a comprehensive test script that tests all endpoints:

```bash
# Make the script executable (first time only)
chmod +x test_api.sh

# Run all tests
./test_api.sh
```

You should see output like:
```
==========================================
  bdSeeker REST API - Complete Test Suite
==========================================

Test 1: Health Check
✓ PASS

Test 2: Register Developer (new user)
✓ PASS

...

==========================================
           TEST SUMMARY
==========================================
Total Tests:  26
Passed:       26
Failed:       0

✓ All tests passed!
```

### Manual Testing with curl

#### 1. Register a New User
```bash
curl -X POST http://localhost:9000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe",
    "role": "developer"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:9000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Save the token from the response for authenticated requests.

#### 3. Get Current User (Protected Route)
```bash
curl -X GET http://localhost:9000/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Testing with Postman

**We provide a complete Postman collection!**

1. **Import the collection**:
   - Open Postman
   - Click **Import**
   - Select `bdSeeker-API.postman_collection.json`
   - Select `bdSeeker-Local.postman_environment.json`

2. **Select environment**: Choose "bdSeeker Local Environment" from the dropdown

3. **Start testing**: The collection includes 24 pre-configured requests with automatic token management!

For detailed Postman usage, see [POSTMAN_GUIDE.md](POSTMAN_GUIDE.md)

## 📚 API Documentation

### Quick Reference

- **Base URL**: `http://localhost:9000/api/v1`
- **Authentication**: JWT Bearer token in `Authorization` header
- **Content-Type**: `application/json`

### Main Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check | No |
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | Login user | No |
| GET | `/auth/me` | Get current user | Yes |
| GET | `/companies` | List companies | No |
| POST | `/companies` | Create company profile | Yes |
| GET | `/developers` | List developers | No |
| POST | `/developers` | Create developer profile | Yes |
| GET | `/jobs` | List jobs | No |
| POST | `/jobs` | Create job post | Yes (Company) |
| POST | `/jobs/:id/reactions` | React to job | Yes |
| POST | `/companies/:id/ratings` | Rate company | Yes |

For complete API documentation, see [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

## 📁 Project Structure

```
bdSeeker-backend/
├── cmd/api/              # Application entry points (optional)
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection & migrations
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # Authentication, CORS, error handling
│   ├── models/           # Database models (GORM)
│   ├── repositories/     # Data access layer
│   └── services/         # Business logic
├── pkg/utils/            # Utility functions (JWT, hashing, validation)
├── .env                  # Environment variables (create from .env.example)
├── main.go               # Application entry point
├── go.mod                # Go module definition
├── go.sum                # Go dependencies checksums
├── test_api.sh           # Automated test script
├── API_DOCUMENTATION.md  # Complete API documentation
└── README.md             # This file
```

## 🔧 Troubleshooting

### Common Issues and Solutions

#### 1. Database Connection Failed

**Error**: `failed to connect to database`

**Solutions**:
- Verify PostgreSQL is running: `sudo systemctl status postgresql`
- Check database exists: `psql -U postgres -l | grep bdseeker`
- Verify credentials in `.env` file match your PostgreSQL setup
- Try connecting manually: `psql -U postgres -d bdseeker`

#### 2. Port Already in Use

**Error**: `bind: address already in use`

**Solutions**:
```bash
# Find process using port 9000
lsof -i :9000

# Kill the process
kill -9 <PID>

# Or change the port in .env file
SERVER_PORT=8081
```

#### 3. Module Not Found

**Error**: `cannot find module`

**Solution**:
```bash
# Clean and reinstall dependencies
go clean -modcache
go mod download
go mod tidy
```

#### 4. Database Migration Errors

**Error**: `failed to run migrations`

**Solutions**:
```bash
# Drop and recreate database
psql -U postgres -c "DROP DATABASE bdseeker;"
psql -U postgres -c "CREATE DATABASE bdseeker;"

# Restart the application (migrations run automatically)
./bdseeker-api
```

#### 5. JWT Token Expired

**Error**: `Invalid or expired token`

**Solution**:
- Login again to get a new token
- Tokens expire after 24 hours by default (configurable in `.env`)

### Getting Help

If you encounter issues not listed here:

1. Check the server logs for detailed error messages
2. Verify all prerequisites are correctly installed
3. Ensure `.env` file is properly configured
4. Review the [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
5. Open an issue on GitHub with:
   - Error message
   - Steps to reproduce
   - Your environment (OS, Go version, PostgreSQL version)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

**Zubair Ahmed Rafi**
- GitHub: [@walleeva2018](https://github.com/walleeva2018)

**Bishworup**
- GitHub: [@bishworup11](https://github.com/bishworup11)

## 🙏 Acknowledgments

- Built following clean architecture principles
- Inspired by HackerRank Jobs platform
- Uses industry-standard Go libraries and best practices

---

**Happy Coding! 🎉**

For questions or support, please open an issue on GitHub.

# API Endpoints Summary

## Base URL
`http://localhost:8080/api/v1`

## Authentication Endpoints

### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe",
  "role": "developer" // or "company" or "admin"
}
```

### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Get Current User (Protected)
```http
GET /auth/me
Authorization: Bearer <token>
```

## Company Endpoints

### List Companies
```http
GET /companies?page=1&limit=10&location=New%20York
```

### Get Company Details
```http
GET /companies/:id
```

### Create Company Profile (Protected)
```http
POST /companies
Authorization: Bearer <token>
Content-Type: application/json

{
  "company_name": "TechCorp Inc",
  "description": "Leading tech company",
  "website": "https://techcorp.com",
  "location": "New York, NY"
}
```

### Rate Company (Protected)
```http
POST /companies/:id/ratings
Authorization: Bearer <token>
Content-Type: application/json

{
  "rating": 5
}
```

### Review Company (Protected)
```http
POST /companies/:id/reviews
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Great company to work with!"
}
```

### List Company Reviews
```http
GET /companies/:id/reviews?page=1&limit=10
```

## Developer Endpoints

### List Developers
```http
GET /developers?page=1&limit=10
```

### Get Developer Details
```http
GET /developers/:id
```

### Create Developer Profile (Protected)
```http
POST /developers
Authorization: Bearer <token>
Content-Type: application/json

{
  "bio": "Experienced full-stack developer"
}
```

## Job Endpoints

### List Jobs
```http
GET /jobs?page=1&limit=10&work_mode=remote&location=Remote&search=backend
```

### Get Job Details
```http
GET /jobs/:id
```

### Create Job Post (Protected - Company only)
```http
POST /jobs
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Senior Backend Developer",
  "description": "Looking for experienced Go developer",
  "salary_min": 100000,
  "salary_max": 150000,
  "experience_min_years": 5,
  "experience_max_years": 10,
  "work_mode": "remote",
  "location": "Remote"
}
```

### React to Job (Protected)
```http
POST /jobs/:id/reactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "type": "apply" // or "like" or "bookmark"
}
```

### Comment on Job (Protected)
```http
POST /jobs/:id/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "This is a great opportunity!"
}
```

## Technology & Language Endpoints

### List Technologies
```http
GET /technologies?search=react
```

### List Programming Languages
```http
GET /languages?search=go
```

## Health Check

### Server Health
```http
GET /health
```

## Query Parameters

### Pagination
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)

### Job Filters
- `work_mode` - Filter by work mode (office/hybrid/remote/onsite)
- `location` - Filter by location (ILIKE search)
- `search` - Search in title and description
- `sort_by` - Sort results (created_desc, created_asc, salary_desc, salary_asc)

### Company Filters
- `location` - Filter by location
- `tech_ids` - Filter by technology IDs (comma-separated)

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": {
    "data": [ ... ],
    "total_count": 100,
    "page": 1,
    "limit": 10,
    "total_pages": 10
  }
}
```

## HTTP Status Codes

- `200 OK` - Successful GET request
- `201 Created` - Successful POST request
- `400 Bad Request` - Validation error or malformed request
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Testing

Run the comprehensive test suite:
```bash
./test_api.sh
```

All 26 tests should pass, covering:
- Authentication flow
- Profile creation
- CRUD operations
- Pagination
- Filtering and search
- Protected routes
- Error handling

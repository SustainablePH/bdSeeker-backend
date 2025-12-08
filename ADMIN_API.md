# 🔐 Admin API Documentation

## Default Admin Credentials

**Email**: `admin@bdseeker.com`  
**Password**: `admin123`

> ⚠️ **IMPORTANT**: Change the admin password in production!

The default admin user is automatically created when the application starts for the first time.

---

## Admin Endpoints

All admin endpoints require:
- **Authentication**: JWT Bearer token
- **Authorization**: Admin role
- **Base URL**: `http://localhost:9000/api/v1/admin`

### Authentication

First, login as admin to get the token:

```bash
curl -X POST http://localhost:9000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@bdseeker.com",
    "password": "admin123"
  }'
```

Use the returned token in the `Authorization` header for all admin requests.

---

## 📊 Statistics

### Get Platform Statistics

```http
GET /api/v1/admin/stats
Authorization: Bearer <admin_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "total_users": 10,
    "total_developers": 5,
    "total_companies": 3,
    "total_jobs": 15,
    "total_reports": 2,
    "pending_reviews": 8
  }
}
```

**Example:**
```bash
curl -X GET http://localhost:9000/api/v1/admin/stats \
  -H "Authorization: Bearer <admin_token>"
```

---

## 👥 User Management

### List All Users

```http
GET /api/v1/admin/users?page=1&limit=10&role=developer
Authorization: Bearer <admin_token>
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `limit` (optional) - Items per page (default: 10)
- `role` (optional) - Filter by role (developer/company/admin)

**Response:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "email": "user@example.com",
        "full_name": "John Doe",
        "role": "developer",
        "created_at": "2025-12-07T00:00:00Z",
        "updated_at": "2025-12-07T00:00:00Z"
      }
    ],
    "total_count": 10,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

**Example:**
```bash
curl -X GET "http://localhost:9000/api/v1/admin/users?page=1&limit=10" \
  -H "Authorization: Bearer <admin_token>"
```

### Delete User

```http
DELETE /api/v1/admin/users/:id
Authorization: Bearer <admin_token>
```

**Example:**
```bash
curl -X DELETE http://localhost:9000/api/v1/admin/users/5 \
  -H "Authorization: Bearer <admin_token>"
```

---

## ⭐ Review Management

### List Pending Reviews

```http
GET /api/v1/admin/reviews/pending?page=1&limit=10
Authorization: Bearer <admin_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Pending reviews retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "company_id": 2,
        "user_id": 3,
        "content": "Great company to work with!",
        "is_approved": false,
        "created_at": "2025-12-07T00:00:00Z"
      }
    ],
    "total_count": 5,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

**Example:**
```bash
curl -X GET "http://localhost:9000/api/v1/admin/reviews/pending?page=1&limit=10" \
  -H "Authorization: Bearer <admin_token>"
```

### Approve Review

```http
PUT /api/v1/admin/reviews/:id/approve
Authorization: Bearer <admin_token>
```

**Example:**
```bash
curl -X PUT http://localhost:9000/api/v1/admin/reviews/1/approve \
  -H "Authorization: Bearer <admin_token>"
```

**Response:**
```json
{
  "success": true,
  "message": "Review approved successfully",
  "data": {
    "id": 1,
    "company_id": 2,
    "user_id": 3,
    "content": "Great company to work with!",
    "is_approved": true,
    "created_at": "2025-12-07T00:00:00Z",
    "updated_at": "2025-12-07T01:00:00Z"
  }
}
```

### Reject Review

```http
DELETE /api/v1/admin/reviews/:id/reject
Authorization: Bearer <admin_token>
```

**Example:**
```bash
curl -X DELETE http://localhost:9000/api/v1/admin/reviews/1/reject \
  -H "Authorization: Bearer <admin_token>"
```

---

## 💬 Comment Management

### Approve Comment

```http
PUT /api/v1/admin/comments/:id/approve
Authorization: Bearer <admin_token>
```

**Example:**
```bash
curl -X PUT http://localhost:9000/api/v1/admin/comments/5/approve \
  -H "Authorization: Bearer <admin_token>"
```

---

## 🚨 Report Management

### List Reports

```http
GET /api/v1/admin/reports?page=1&limit=10&status=pending&type=user
Authorization: Bearer <admin_token>
```

**Query Parameters:**
- `page` (optional) - Page number
- `limit` (optional) - Items per page
- `status` (optional) - Filter by status (pending/reviewed/resolved/dismissed)
- `type` (optional) - Filter by report type

**Response:**
```json
{
  "success": true,
  "message": "Reports retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "reporter_id": 3,
        "reported_id": 5,
        "report_type": "user",
        "reason": "Inappropriate behavior",
        "status": "pending",
        "created_at": "2025-12-07T00:00:00Z"
      }
    ],
    "total_count": 2,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

**Example:**
```bash
curl -X GET "http://localhost:9000/api/v1/admin/reports?status=pending" \
  -H "Authorization: Bearer <admin_token>"
```

### Update Report Status

```http
PUT /api/v1/admin/reports/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "status": "resolved"
}
```

**Status Values:**
- `reviewed` - Report has been reviewed
- `resolved` - Issue has been resolved
- `dismissed` - Report was dismissed

**Example:**
```bash
curl -X PUT http://localhost:9000/api/v1/admin/reports/1 \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"status": "resolved"}'
```

**Response:**
```json
{
  "success": true,
  "message": "Report status updated successfully",
  "data": {
    "id": 1,
    "reporter_id": 3,
    "reported_id": 5,
    "report_type": "user",
    "reason": "Inappropriate behavior",
    "status": "resolved",
    "reviewed_by": 10,
    "reviewed_at": "2025-12-07T01:00:00Z",
    "created_at": "2025-12-07T00:00:00Z",
    "updated_at": "2025-12-07T01:00:00Z"
  }
}
```

---

## 📋 Complete Admin API Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/stats` | Get platform statistics |
| GET | `/admin/users` | List all users (with filters) |
| DELETE | `/admin/users/:id` | Delete a user |
| GET | `/admin/reviews/pending` | List pending reviews |
| PUT | `/admin/reviews/:id/approve` | Approve a review |
| DELETE | `/admin/reviews/:id/reject` | Reject/delete a review |
| PUT | `/admin/comments/:id/approve` | Approve a comment |
| GET | `/admin/reports` | List reports (with filters) |
| PUT | `/admin/reports/:id` | Update report status |

---

## 🔒 Security Notes

1. **Admin Only**: All endpoints require admin role
2. **JWT Required**: Valid JWT token must be provided
3. **Change Default Password**: Always change the default admin password in production
4. **Audit Logs**: Consider implementing audit logging for admin actions
5. **Rate Limiting**: Implement rate limiting for admin endpoints in production

---

## 🧪 Testing Admin APIs

### Quick Test Script

```bash
#!/bin/bash

# Login as admin
ADMIN_TOKEN=$(curl -s -X POST http://localhost:9000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@bdseeker.com","password":"admin123"}' \
  | jq -r '.data.token')

echo "Admin Token: $ADMIN_TOKEN"

# Get stats
echo "\n=== Platform Statistics ==="
curl -s -X GET http://localhost:9000/api/v1/admin/stats \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# List users
echo "\n=== All Users ==="
curl -s -X GET "http://localhost:9000/api/v1/admin/users?page=1&limit=5" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# List pending reviews
echo "\n=== Pending Reviews ==="
curl -s -X GET "http://localhost:9000/api/v1/admin/reviews/pending?page=1&limit=5" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# List reports
echo "\n=== Reports ==="
curl -s -X GET "http://localhost:9000/api/v1/admin/reports?page=1&limit=5" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
```

Save as `test_admin.sh`, make executable with `chmod +x test_admin.sh`, and run!

---

## 💡 Common Use Cases

### 1. Moderate Company Reviews
```bash
# List pending reviews
curl -X GET "http://localhost:9000/api/v1/admin/reviews/pending" \
  -H "Authorization: Bearer <admin_token>"

# Approve a review
curl -X PUT http://localhost:9000/api/v1/admin/reviews/1/approve \
  -H "Authorization: Bearer <admin_token>"
```

### 2. Handle User Reports
```bash
# List pending reports
curl -X GET "http://localhost:9000/api/v1/admin/reports?status=pending" \
  -H "Authorization: Bearer <admin_token>"

# Resolve a report
curl -X PUT http://localhost:9000/api/v1/admin/reports/1 \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"status": "resolved"}'
```

### 3. Monitor Platform Growth
```bash
# Get current statistics
curl -X GET http://localhost:9000/api/v1/admin/stats \
  -H "Authorization: Bearer <admin_token>"
```

### 4. Manage Users
```bash
# List all developers
curl -X GET "http://localhost:9000/api/v1/admin/users?role=developer" \
  -H "Authorization: Bearer <admin_token>"

# Delete a problematic user
curl -X DELETE http://localhost:9000/api/v1/admin/users/5 \
  -H "Authorization: Bearer <admin_token>"
```

---

## 🚀 Next Steps

1. **Change Admin Password**: Use the `/auth/me` endpoint to update password
2. **Create Additional Admins**: Register new users with admin role
3. **Implement Audit Logging**: Track all admin actions
4. **Add Email Notifications**: Notify users when reviews are approved/rejected
5. **Dashboard UI**: Build an admin dashboard for easier management

---

For more information, see the main [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

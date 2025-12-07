

# 🧠 Backend Architecture Instructions (Go + GORM + PostgreSQL)

This document defines the backend architecture for a **job searching web platform** similar to HackerRank jobs.
The system must allow:

* Guest access (read-only)
* Authenticated users (Developer or Company)
* Admin role (full management & moderation)

---

## 📦 Tech Stack

| Layer      | Technology                    |
| ---------- | ----------------------------- |
| Language   | Go (>= 1.21 recommended)      |
| ORM        | GORM                          |
| DB         | PostgreSQL                    |
| Auth       | JWT (Access + Refresh)        |
| API        | REST  |
| Env config | Dotenv / Viper                |
| Migration  | GORM AutoMigrate or Goose     |

---

### 1.3 Create Project Structure
```
project-root/
├── cmd/api/main.go
├── internal/
│   ├── config/         # Environment config
│   ├── models/         # Database models
│   ├── handlers/       # HTTP handlers
│   ├── services/       # Business logic
│   ├── repositories/   # Data access layer
│   ├── middleware/     # Auth, CORS, error handling
│   └── database/       # DB connection & migrations
├── pkg/utils/          # JWT, hash, validators
└── .env
```

### 1.4 Environment Variables (.env)
```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
JWT_SECRET, JWT_EXPIRY
SERVER_PORT
```

---

## Phase 2: Database Models (GORM)

### 2.1 Core Models to Create

**User Model** (`internal/models/user.go`)
- Fields: ID, Email (unique), PasswordHash, FullName, Role (developer/company/admin), CreatedAt, UpdatedAt
- Relations: HasOne CompanyProfile, HasOne DeveloperProfile

**Technology & ProgrammingLanguage** (`internal/models/user.go`)
- Simple lookup tables with ID and Name

**CompanyProfile** (`internal/models/company.go`)
- Fields: ID, UserID (FK), CompanyName, Description, Website, Location, timestamps
- Relations: BelongsTo User, ManyToMany Technologies, HasMany JobPosts, Ratings, Reviews

**CompanyRating, CompanyReview** (`internal/models/company.go`)
- Rating: CompanyID, UserID, Rating (1-5)
- Review: CompanyID, UserID, Content, IsApproved (admin approval), timestamps
- Review Relations: HasMany Reactions, Comments

**CompanyReviewReaction, CompanyReviewComment, CompanyReviewReply**
- Reaction: ReviewID, UserID, Type (like/useful)
- Comment: ReviewID, UserID, Content, IsApproved
- Reply: ReviewCommentID, UserID, Content

**DeveloperProfile** (`internal/models/developer.go`)
- Fields: ID, UserID (FK), Bio, timestamps
- Relations: BelongsTo User, HasMany Experiences/Educations/Certificates, ManyToMany Technologies/Languages

**DeveloperExperience, DeveloperEducation, DeveloperCertificate**
- Experience: DeveloperID, Title, CompanyName, StartDate, EndDate, Description
- Education: DeveloperID, Institution, Degree, FieldOfStudy, dates, Grade, Description
- Certificate: DeveloperID, CertificateName, IssuingOrganization, dates, CredentialID, CertificateLink

**JobPost** (`internal/models/job.go`)
- Fields: ID, CompanyID (FK), Title, Description, SalaryMin/Max, ExperienceMin/MaxYears, WorkMode (office/hybrid/remote/onsite), Location, timestamps
- Relations: BelongsTo Company, HasMany Reactions, Comments

**PostReaction, PostComment, CommentReply** (`internal/models/job.go`)
- Reaction: UserID, JobPostID, Type (like/bookmark/apply)
- Comment: UserID, JobPostID, Content
- Reply: CommentID, UserID, Content (no reply-to-reply)

**UserReport** (`internal/models/report.go`)
- Fields: ReporterID, ReportedID, ReportType, Description, Status (pending/reviewed/resolved/dismissed), ReviewedBy (admin), ReviewedAt, timestamps
- Relations: Reporter (User), Reported (User), Reviewer (User)

### 2.2 GORM Tags to Use
- Primary keys: `gorm:"primaryKey"`
- Foreign keys: `gorm:"not null"` with appropriate field names
- Unique constraints: `gorm:"uniqueIndex"`
- Size limits: `gorm:"size:255"`
- Default values: `gorm:"default:false"`
- JSON serialization: `json:"field_name"` and `json:"-"` for sensitive fields

### 2.3 Many-to-Many Relationships
Implement using GORM's `many2many` tag:
- Company <-> Technologies
- Developer <-> Technologies
- Developer <-> ProgrammingLanguages

---

## Phase 3: Database Layer

### 3.1 Database Connection (`internal/database/database.go`)
- Create Connect() function with PostgreSQL DSN
- Export global DB variable
- Implement proper error handling
- Add logging for connection status

### 3.2 Auto Migration (`internal/database/database.go`)
- Create Migrate() function
- Run DB.AutoMigrate() for all models in correct order
- Handle foreign key dependencies

---

## Phase 4: Authentication & Utilities

### 4.1 JWT Utilities (`pkg/utils/jwt.go`)
**Functions to implement:**
- `GenerateToken(userID, email, role, secret, expiry)` - Returns JWT string
- `ValidateToken(tokenString, secret)` - Returns claims or error
- Define JWTClaims struct with UserID, Email, Role, and jwt.RegisteredClaims

### 4.2 Password Hashing (`pkg/utils/hash.go`)
- `HashPassword(password)` - Returns bcrypt hash
- `CheckPassword(password, hash)` - Returns bool

### 4.3 Validator (`pkg/utils/validator.go`)
- Initialize go-playground/validator instance
- Create helper functions for validation errors

---

## Phase 5: Middleware

### 5.1 Authentication Middleware (`internal/middleware/auth.go`)
**AuthMiddleware:**
- Extract Bearer token from Authorization header
- Validate JWT token
- Set user_id, user_email, user_role in context
- Return 401 if invalid

**RoleMiddleware:**
- Accept allowed roles as parameters
- Check user_role from context
- Return 403 if role not allowed

### 5.2 CORS Middleware
- Allow origins, methods, headers
- Handle OPTIONS preflight

### 5.3 Error Handler Middleware
- Catch panics
- Format error responses consistently

---

## Phase 6: API Endpoints Implementation

### 6.1 Authentication Endpoints (`internal/handlers/auth_handler.go`)

**POST /api/v1/auth/register**
- Accept: email, password, full_name, role
- Validate input
- Hash password
- Create user
- Return JWT token

**POST /api/v1/auth/login**
- Accept: email, password
- Find user by email
- Verify password
- Generate JWT token
- Return token and user info

**GET /api/v1/auth/me** (Protected)
- Get user_id from context
- Fetch user with profile (company or developer)
- Return user data

---

### 6.2 Company Endpoints (`internal/handlers/company_handler.go`)

**GET /api/v1/companies** (Public)
- Support pagination (page, limit)
- Filter by location, technologies
- Preload technologies
- Return list with metadata

**GET /api/v1/companies/:id** (Public)
- Get company by ID
- Preload: User, Technologies, Ratings, Reviews (approved only)
- Calculate average rating
- Return full company profile

**POST /api/v1/companies** (Auth: developer/company)
- Accept: company_name, description, website, location, technology_ids
- Validate user doesn't have existing company profile
- Create company profile
- Associate technologies
- Return created profile

**PUT /api/v1/companies/:id** (Auth: owner/admin)
- Check ownership (user_id matches or admin role)
- Update allowed fields
- Update technologies association
- Return updated profile

**DELETE /api/v1/companies/:id** (Auth: owner/admin)
- Check ownership
- Soft delete or hard delete
- Return success message

---

### 6.3 Company Rating & Review Endpoints

**POST /api/v1/companies/:id/ratings** (Auth)
- Accept: rating (1-5)
- Check if user already rated (upsert logic)
- Create or update rating
- Return rating

**POST /api/v1/companies/:id/reviews** (Auth)
- Accept: content
- Create review with is_approved=false
- Return review (pending approval)

**GET /api/v1/companies/:id/reviews** (Public)
- Get approved reviews only
- Support pagination
- Preload: User, Reactions, Comments
- Return reviews list

**PUT /api/v1/reviews/:id** (Auth: owner/admin)
- Check ownership
- Update content
- Reset is_approved to false if content changed
- Return updated review

**DELETE /api/v1/reviews/:id** (Auth: owner/admin)
- Check ownership or admin
- Delete review
- Return success

**POST /api/v1/reviews/:id/reactions** (Auth)
- Accept: type (like/useful)
- Create or update reaction (upsert)
- Return reaction

**POST /api/v1/reviews/:id/comments** (Auth)
- Accept: content
- Create comment with is_approved=false
- Return comment

**POST /api/v1/review-comments/:id/replies** (Auth)
- Accept: content
- Create reply (no nested replies)
- Return reply

---

### 6.4 Developer Endpoints (`internal/handlers/developer_handler.go`)

**GET /api/v1/developers** (Public)
- Support pagination
- Filter by technologies, languages, experience
- Preload: Technologies, Languages
- Return list

**GET /api/v1/developers/:id** (Public)
- Get developer by ID
- Preload: User, Experiences, Educations, Certificates, Technologies, Languages
- Return full profile

**POST /api/v1/developers** (Auth: developer)
- Accept: bio
- Validate user doesn't have existing developer profile
- Create developer profile
- Return created profile

**PUT /api/v1/developers/:id** (Auth: owner)
- Check ownership
- Update bio, technologies, languages
- Return updated profile

**POST /api/v1/developers/:id/experiences** (Auth: owner)
- Accept: title, company_name, start_date, end_date, description
- Create experience entry
- Return experience

**PUT /api/v1/experiences/:id** (Auth: owner)
- Check ownership via developer_id
- Update experience
- Return updated experience

**DELETE /api/v1/experiences/:id** (Auth: owner)
- Check ownership
- Delete experience
- Return success

**POST /api/v1/developers/:id/educations** (Auth: owner)
- Accept: institution, degree, field_of_study, dates, grade, description
- Create education entry
- Return education

**PUT /api/v1/educations/:id** (Auth: owner)
- Check ownership
- Update education
- Return updated education

**DELETE /api/v1/educations/:id** (Auth: owner)
- Check ownership
- Delete education
- Return success

**POST /api/v1/developers/:id/certificates** (Auth: owner)
- Accept: certificate_name, issuing_organization, dates, credential_id, certificate_link, description
- Validate certificate_link URL
- Create certificate entry
- Return certificate

**PUT /api/v1/certificates/:id** (Auth: owner)
- Check ownership
- Update certificate
- Return updated certificate

**DELETE /api/v1/certificates/:id** (Auth: owner)
- Check ownership
- Delete certificate
- Return success

---

### 6.5 Job Post Endpoints (`internal/handlers/job_handler.go`)

**GET /api/v1/jobs** (Public)
- Support pagination
- Filter by: company_id, work_mode, location, salary_range, experience_range
- Search by title/description (ILIKE)
- Sort by: created_at, salary
- Preload Company
- Return list with metadata

**GET /api/v1/jobs/:id** (Public)
- Get job by ID
- Preload: Company (with User), Reactions, Comments (with Replies)
- Return full job details

**POST /api/v1/jobs** (Auth: company/admin)
- Accept: title, description, salary_min/max, experience_min/max_years, work_mode, location
- Validate user has company profile (if not admin)
- Get company_id from user's company profile
- Create job post
- Return created job

**PUT /api/v1/jobs/:id** (Auth: owner/admin)
- Check ownership (company_id matches user's company)
- Update job fields
- Return updated job

**DELETE /api/v1/jobs/:id** (Auth: owner/admin)
- Check ownership
- Delete job post
- Return success

**POST /api/v1/jobs/:id/reactions** (Auth)
- Accept: type (like/bookmark/apply)
- Create or update reaction (upsert by user_id + job_post_id)
- Return reaction

**DELETE /api/v1/jobs/:id/reactions** (Auth)
- Delete user's reaction to job
- Return success

**POST /api/v1/jobs/:id/comments** (Auth)
- Accept: content
- Create comment
- Return comment with user info

**PUT /api/v1/comments/:id** (Auth: owner)
- Check ownership
- Update comment content
- Return updated comment

**DELETE /api/v1/comments/:id** (Auth: owner/admin)
- Check ownership or admin
- Delete comment
- Return success

**POST /api/v1/comments/:id/replies** (Auth)
- Accept: content
- Create reply (no nested replies)
- Return reply with user info

---

### 6.6 Report Endpoints (`internal/handlers/report_handler.go`)

**POST /api/v1/reports** (Auth)
- Accept: reported_id, report_type, description
- Validate reporter_id != reported_id
- Create report with status="pending"
- Return created report

**GET /api/v1/reports** (Auth: admin)
- Support pagination
- Filter by status, report_type
- Preload: Reporter, Reported, Reviewer
- Return reports list

**GET /api/v1/reports/:id** (Auth: admin)
- Get report by ID
- Preload all relations
- Return full report details

**PUT /api/v1/reports/:id** (Auth: admin)
- Accept: status (reviewed/resolved/dismissed)
- Update status, reviewed_by, reviewed_at
- Return updated report

---

### 6.7 Technology & Language Endpoints (`internal/handlers/tech_handler.go`)

**GET /api/v1/technologies** (Public)
- List all technologies
- Optional search by name
- Return list

**POST /api/v1/technologies** (Auth: admin)
- Accept: name
- Check uniqueness
- Create technology
- Return created technology

**GET /api/v1/languages** (Public)
- List all programming languages
- Optional search by name
- Return list

**POST /api/v1/languages** (Auth: admin)
- Accept: name
- Check uniqueness
- Create language
- Return created language

---

## Phase 7: Service Layer (Business Logic)

### 7.1 Create Services for Each Domain
**Services to implement:**
- `auth_service.go` - Register, Login, Token refresh
- `user_service.go` - User CRUD, profile management
- `company_service.go` - Company CRUD, rating/review logic
- `developer_service.go` - Developer CRUD, experience/education/certificate logic
- `job_service.go` - Job CRUD, reaction/comment logic
- `report_service.go` - Report management, admin actions

### 7.2 Service Layer Responsibilities
- Input validation (using validator)
- Business rule enforcement
- Transaction management for complex operations
- Call repository layer for data access
- Return DTOs (Data Transfer Objects) not raw models

---

## Phase 8: Repository Layer (Data Access)

### 8.1 Create Repositories
**Repositories to implement:**
- `user_repository.go` - User database operations
- `company_repository.go` - Company database operations
- `developer_repository.go` - Developer database operations
- `job_repository.go` - Job database operations
- `report_repository.go` - Report database operations

### 8.2 Repository Responsibilities
- Raw GORM queries
- Complex joins and preloads
- Pagination helpers
- Filter building
- Transaction support

---

## Phase 9: Advanced Features

### 9.1 Search & Filtering
**Implement for:**
- Job posts: by title, description, location, salary, experience, work_mode
- Companies: by name, location, technologies
- Developers: by technologies, languages, experience level

**Use GORM features:**
- `Where()` for conditions
- `Preload()` for relations
- `Joins()` for complex queries
- `Order()` for sorting
- `Limit()` and `Offset()` for pagination

### 9.2 Pagination Helper
Create utility function:
```go
type PaginationParams struct {
    Page  int
    Limit int
}

type PaginationResult struct {
    Data       interface{}
    TotalCount int64
    Page       int
    Limit      int
    TotalPages int
}
```

### 9.3 Response Formatting
Standardize API responses:
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

---

## Phase 10: Admin Features

### 10.1 Admin Dashboard Endpoints
**Additional endpoints for admin:**
- `GET /api/v1/admin/stats` - Platform statistics
- `GET /api/v1/admin/users` - List all users with filters
- `PUT /api/v1/admin/users/:id/status` - Ban/unban users
- `PUT /api/v1/admin/reviews/:id/approve` - Approve/reject reviews
- `PUT /api/v1/admin/review-comments/:id/approve` - Approve/reject comments

### 10.2 Moderation System
Implement approval workflows for:
- Company reviews (is_approved flag)
- Review comments (is_approved flag)
- User reports (status field)

---

## Phase 11: Testing

### 11.1 Unit Tests
Write tests for:
- Utilities (JWT, hash, validators)
- Services (business logic)
- Middleware

### 11.2 Integration Tests
Test:
- Authentication flow
- CRUD operations for all entities
- Complex queries and filters
- Relationship management

### 11.3 API Testing
Use tools:
- Postman/Insomnia collections
- Automated API tests with Go testing package

---

## Phase 12: Documentation & Deployment

### 12.1 API Documentation
**Generate using:**
- Swagger/OpenAPI specification
- Document all endpoints, request/response schemas
- Include authentication requirements
- Provide example requests

### 12.2 Deployment Preparation
- Create Dockerfile
- Docker Compose for local development
- Environment-specific configs
- Database migration scripts
- CI/CD pipeline setup

### 12.3 Security Checklist
- [ ] Rate limiting implemented
- [ ] Input validation on all endpoints
- [ ] SQL injection protection (GORM handles this)
- [ ] XSS protection
- [ ] CSRF protection
- [ ] Secure password requirements
- [ ] JWT token expiration
- [ ] HTTPS in production
- [ ] Environment variables for secrets

---

## Key Implementation Notes

### GORM Best Practices
1. Use `Preload()` strategically to avoid N+1 queries
2. Use `Select()` to limit fields returned
3. Use transactions for multi-step operations
4. Implement soft deletes where appropriate
5. Add indexes for frequently queried fields

### Error Handling
1. Create custom error types
2. Use proper HTTP status codes
3. Don't expose internal errors to clients
4. Log errors server-side

### Security Considerations
1. Never return password hashes in API responses
2. Validate all user input
3. Check ownership before update/delete operations
4. Implement rate limiting for sensitive endpoints
5. Use parameterized queries (GORM handles this)

### Performance Optimization
1. Add database indexes on foreign keys and frequently searched fields
2. Implement caching for lookup tables (technologies, languages)
3. Use database connection pooling
4. Optimize complex queries with proper joins
5. Implement pagination on all list endpoints

---

## Quick Start Commands

```bash
# Run migrations
go run cmd/api/main.go migrate

# Start server
go run cmd/api/main.go serve

# Run tests
go test ./...

# Build for production
go build -o bin/api cmd/api/main.go
```


# 📮 Postman Collection Guide

## Files Included

1. **bdSeeker-API.postman_collection.json** - Complete API collection with all endpoints
2. **bdSeeker-Local.postman_environment.json** - Environment variables for local development

## 🚀 Quick Setup

### Step 1: Import Collection

1. Open Postman
2. Click **Import** button (top left)
3. Drag and drop `bdSeeker-API.postman_collection.json`
4. Click **Import**

### Step 2: Import Environment

1. Click **Import** button again
2. Drag and drop `bdSeeker-Local.postman_environment.json`
3. Click **Import**
4. Select **bdSeeker Local Environment** from the environment dropdown (top right)

### Step 3: Start Testing!

The collection is organized into folders:
- **Authentication** - Register, login, get current user
- **Companies** - Company CRUD, ratings, reviews
- **Developers** - Developer profiles
- **Jobs** - Job postings, reactions, comments
- **Technologies & Languages** - Lookup data
- **System** - Health check

## 🔄 Automated Token Management

The collection includes automatic token management:

1. **Register Developer** - Automatically saves `dev_token` and `dev_user_id`
2. **Register Company** - Automatically saves `company_token` and `company_user_id`
3. **Login** - Automatically saves `auth_token`
4. **Create Developer Profile** - Automatically saves `developer_id`
5. **Create Company Profile** - Automatically saves `company_id`
6. **Create Job Post** - Automatically saves `job_id`

These tokens are automatically used in subsequent requests!

## 📝 Usage Workflow

### For Testing Complete Flow:

1. **Register Users**
   - Run "Register Developer" (saves dev_token)
   - Run "Register Company" (saves company_token)

2. **Create Profiles**
   - Run "Create Developer Profile" (uses dev_token)
   - Run "Create Company Profile" (uses company_token)

3. **Create Job**
   - Run "Create Job Post" (uses company_token, saves job_id)

4. **Interact with Job**
   - Run "Apply to Job" (uses dev_token and job_id)
   - Run "Bookmark Job"
   - Run "Comment on Job"

5. **Rate & Review**
   - Run "Rate Company" (uses dev_token and company_id)
   - Run "Review Company"

## 🔧 Environment Variables

The environment includes these variables:

| Variable | Description | Auto-set |
|----------|-------------|----------|
| `base_url` | API base URL | No |
| `auth_token` | General auth token | Yes (on login) |
| `dev_token` | Developer user token | Yes (on register) |
| `company_token` | Company user token | Yes (on register) |
| `dev_user_id` | Developer user ID | Yes |
| `company_user_id` | Company user ID | Yes |
| `developer_id` | Developer profile ID | Yes |
| `company_id` | Company profile ID | Yes |
| `job_id` | Job post ID | Yes |

## 📋 Available Endpoints

### Authentication (4 requests)
- ✅ Register Developer
- ✅ Register Company
- ✅ Login
- ✅ Get Current User

### Companies (6 requests)
- ✅ List Companies
- ✅ Get Company Details
- ✅ Create Company Profile
- ✅ Rate Company
- ✅ Review Company
- ✅ List Company Reviews

### Developers (3 requests)
- ✅ List Developers
- ✅ Get Developer Details
- ✅ Create Developer Profile

### Jobs (8 requests)
- ✅ List Jobs
- ✅ Search Jobs (Remote)
- ✅ Get Job Details
- ✅ Create Job Post
- ✅ Apply to Job
- ✅ Bookmark Job
- ✅ Like Job
- ✅ Comment on Job

### Technologies & Languages (2 requests)
- ✅ List Technologies
- ✅ List Programming Languages

### System (1 request)
- ✅ Health Check

**Total: 24 pre-configured requests**

## 🎯 Query Parameters

Many endpoints support query parameters (already configured but disabled by default):

### List Endpoints
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)

### Jobs
- `work_mode` - Filter by remote/hybrid/office/onsite
- `location` - Filter by location
- `search` - Search in title/description
- `sort_by` - Sort results (created_desc, salary_desc, etc.)

### Technologies/Languages
- `search` - Search by name

To use them, enable the query parameter in Postman!

## 🔐 Authentication

Protected endpoints automatically use the appropriate token:
- Developer endpoints use `{{dev_token}}`
- Company endpoints use `{{company_token}}`
- General protected endpoints use `{{auth_token}}`

## 💡 Tips

1. **Run in Order**: For first-time testing, run requests in order within each folder
2. **Check Environment**: Ensure "bdSeeker Local Environment" is selected
3. **View Tokens**: Click the eye icon next to environment to see saved tokens
4. **Reset Tokens**: Clear environment variables to start fresh
5. **Change Base URL**: Update `base_url` if server runs on different port

## 🐛 Troubleshooting

### "Invalid or expired token"
- Re-run the login or register request to get a new token

### "User already has a profile"
- Use different email or clear database

### "Connection refused"
- Make sure the server is running: `./bdseeker-api`

### Variables not saving
- Check the "Tests" tab in each request for the auto-save scripts

## 📚 Additional Resources

- **API Documentation**: See `API_DOCUMENTATION.md`
- **Setup Guide**: See `README.md`
- **Quick Start**: See `QUICKSTART.md`

---

**Happy Testing! 🚀**

For issues or questions, check the main README.md file.

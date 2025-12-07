#!/bin/bash

# Comprehensive API Testing Script
BASE_URL="http://localhost:8080/api/v1"

echo "=========================================="
echo "  bdSeeker REST API - Complete Test Suite"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

test_count=0
pass_count=0
fail_count=0

run_test() {
    test_count=$((test_count + 1))
    echo -e "${BLUE}Test $test_count: $1${NC}"
}

pass() {
    pass_count=$((pass_count + 1))
    echo -e "${GREEN}✓ PASS${NC}"
    echo ""
}

fail() {
    fail_count=$((fail_count + 1))
    echo -e "${RED}✗ FAIL: $1${NC}"
    echo ""
}

# Test 1: Health Check
run_test "Health Check"
HEALTH=$(curl -s $BASE_URL/health)
if echo $HEALTH | grep -q "ok"; then
    pass
else
    fail "Health check failed"
fi

# Test 2: Register New Developer
run_test "Register Developer (new user)"
TIMESTAMP=$(date +%s)
DEV_EMAIL="dev${TIMESTAMP}@test.com"
DEV_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$DEV_EMAIL\",\"password\":\"password123\",\"full_name\":\"Test Developer\",\"role\":\"developer\"}")

if echo $DEV_RESPONSE | grep -q "success.*true"; then
    DEV_TOKEN=$(echo $DEV_RESPONSE | jq -r '.data.token')
    DEV_ID=$(echo $DEV_RESPONSE | jq -r '.data.user.id')
    echo "Developer ID: $DEV_ID"
    echo "Token: ${DEV_TOKEN:0:50}..."
    pass
else
    fail "Developer registration failed"
    echo $DEV_RESPONSE | jq .
fi

# Test 3: Register New Company
run_test "Register Company User (new user)"
COMPANY_EMAIL="company${TIMESTAMP}@test.com"
COMPANY_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$COMPANY_EMAIL\",\"password\":\"password123\",\"full_name\":\"Test Company\",\"role\":\"company\"}")

if echo $COMPANY_RESPONSE | grep -q "success.*true"; then
    COMPANY_TOKEN=$(echo $COMPANY_RESPONSE | jq -r '.data.token')
    COMPANY_ID=$(echo $COMPANY_RESPONSE | jq -r '.data.user.id')
    echo "Company ID: $COMPANY_ID"
    echo "Token: ${COMPANY_TOKEN:0:50}..."
    pass
else
    fail "Company registration failed"
    echo $COMPANY_RESPONSE | jq .
fi

# Test 4: Login
run_test "Login with Developer Credentials"
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$DEV_EMAIL\",\"password\":\"password123\"}")

if echo $LOGIN_RESPONSE | grep -q "success.*true"; then
    pass
else
    fail "Login failed"
fi

# Test 5: Get Current User (Protected)
run_test "Get Current User Info (Protected Route)"
ME_RESPONSE=$(curl -s -X GET $BASE_URL/auth/me \
  -H "Authorization: Bearer $DEV_TOKEN")

if echo $ME_RESPONSE | grep -q "success.*true"; then
    echo $ME_RESPONSE | jq '.data | {id, email, role}'
    pass
else
    fail "Get current user failed"
fi

# Test 6: Create Developer Profile
run_test "Create Developer Profile"
DEV_PROFILE=$(curl -s -X POST $BASE_URL/developers \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"bio":"Experienced full-stack developer with expertise in Go and React"}')

if echo $DEV_PROFILE | grep -q "success.*true"; then
    DEV_PROFILE_ID=$(echo $DEV_PROFILE | jq -r '.data.id')
    echo "Developer Profile ID: $DEV_PROFILE_ID"
    pass
else
    fail "Developer profile creation failed"
    echo $DEV_PROFILE | jq .
fi

# Test 7: Create Company Profile
run_test "Create Company Profile"
COMPANY_PROFILE=$(curl -s -X POST $BASE_URL/companies \
  -H "Authorization: Bearer $COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"company_name":"TechCorp Inc","description":"Leading technology company","website":"https://techcorp.com","location":"New York, NY"}')

if echo $COMPANY_PROFILE | grep -q "success.*true"; then
    COMPANY_PROFILE_ID=$(echo $COMPANY_PROFILE | jq -r '.data.id')
    echo "Company Profile ID: $COMPANY_PROFILE_ID"
    pass
else
    fail "Company profile creation failed"
    echo $COMPANY_PROFILE | jq .
fi

# Test 8: List Companies
run_test "List All Companies (Pagination)"
COMPANIES=$(curl -s "$BASE_URL/companies?page=1&limit=5")
if echo $COMPANIES | grep -q "success.*true"; then
    COMPANY_COUNT=$(echo $COMPANIES | jq '.data.total_count')
    echo "Total companies: $COMPANY_COUNT"
    pass
else
    fail "List companies failed"
fi

# Test 9: Get Company Details
run_test "Get Specific Company Details"
COMPANY_DETAIL=$(curl -s "$BASE_URL/companies/$COMPANY_PROFILE_ID")
if echo $COMPANY_DETAIL | grep -q "success.*true"; then
    echo $COMPANY_DETAIL | jq '.data | {id, company_name, location}'
    pass
else
    fail "Get company details failed"
fi

# Test 10: List Developers
run_test "List All Developers (Pagination)"
DEVELOPERS=$(curl -s "$BASE_URL/developers?page=1&limit=5")
if echo $DEVELOPERS | grep -q "success.*true"; then
    DEV_COUNT=$(echo $DEVELOPERS | jq '.data.total_count')
    echo "Total developers: $DEV_COUNT"
    pass
else
    fail "List developers failed"
fi

# Test 11: Get Developer Details
run_test "Get Specific Developer Details"
DEV_DETAIL=$(curl -s "$BASE_URL/developers/$DEV_PROFILE_ID")
if echo $DEV_DETAIL | grep -q "success.*true"; then
    echo $DEV_DETAIL | jq '.data | {id, bio}'
    pass
else
    fail "Get developer details failed"
fi

# Test 12: Create Job Post
run_test "Create Job Post (Company)"
JOB_POST=$(curl -s -X POST $BASE_URL/jobs \
  -H "Authorization: Bearer $COMPANY_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Senior Backend Developer",
    "description":"Looking for experienced Go developer",
    "salary_min":100000,
    "salary_max":150000,
    "experience_min_years":5,
    "experience_max_years":10,
    "work_mode":"remote",
    "location":"Remote"
  }')

if echo $JOB_POST | grep -q "success.*true"; then
    JOB_ID=$(echo $JOB_POST | jq -r '.data.id')
    echo "Job ID: $JOB_ID"
    pass
else
    fail "Job post creation failed"
    echo $JOB_POST | jq .
fi

# Test 13: List Jobs
run_test "List All Jobs"
JOBS=$(curl -s "$BASE_URL/jobs?page=1&limit=10")
if echo $JOBS | grep -q "success.*true"; then
    JOB_COUNT=$(echo $JOBS | jq '.data.total_count')
    echo "Total jobs: $JOB_COUNT"
    pass
else
    fail "List jobs failed"
fi

# Test 14: Get Job Details
run_test "Get Specific Job Details"
JOB_DETAIL=$(curl -s "$BASE_URL/jobs/$JOB_ID")
if echo $JOB_DETAIL | grep -q "success.*true"; then
    echo $JOB_DETAIL | jq '.data | {id, title, work_mode, location}'
    pass
else
    fail "Get job details failed"
fi

# Test 15: React to Job (Apply)
run_test "React to Job - Apply"
REACTION=$(curl -s -X POST "$BASE_URL/jobs/$JOB_ID/reactions" \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type":"apply"}')

if echo $REACTION | grep -q "success.*true"; then
    pass
else
    fail "Job reaction failed"
fi

# Test 16: React to Job (Bookmark)
run_test "React to Job - Bookmark"
BOOKMARK=$(curl -s -X POST "$BASE_URL/jobs/$JOB_ID/reactions" \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type":"bookmark"}')

if echo $BOOKMARK | grep -q "success.*true"; then
    pass
else
    fail "Job bookmark failed"
fi

# Test 17: Comment on Job
run_test "Comment on Job Post"
COMMENT=$(curl -s -X POST "$BASE_URL/jobs/$JOB_ID/comments" \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"This is a great opportunity! When can I start?"}')

if echo $COMMENT | grep -q "success.*true"; then
    pass
else
    fail "Job comment failed"
fi

# Test 18: Rate Company
run_test "Rate Company (5 stars)"
RATING=$(curl -s -X POST "$BASE_URL/companies/$COMPANY_PROFILE_ID/ratings" \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"rating":5}')

if echo $RATING | grep -q "success.*true"; then
    pass
else
    fail "Company rating failed"
fi

# Test 19: Review Company
run_test "Review Company"
REVIEW=$(curl -s -X POST "$BASE_URL/companies/$COMPANY_PROFILE_ID/reviews" \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Excellent company culture and great benefits. Highly recommend!"}')

if echo $REVIEW | grep -q "success.*true"; then
    echo "Review created (pending approval)"
    pass
else
    fail "Company review failed"
fi

# Test 20: List Company Reviews
run_test "List Company Reviews"
REVIEWS=$(curl -s "$BASE_URL/companies/$COMPANY_PROFILE_ID/reviews?page=1&limit=10")
if echo $REVIEWS | grep -q "success.*true"; then
    REVIEW_COUNT=$(echo $REVIEWS | jq '.data.total_count')
    echo "Total reviews: $REVIEW_COUNT"
    pass
else
    fail "List reviews failed"
fi

# Test 21: Search Jobs by Work Mode
run_test "Search Jobs (work_mode=remote)"
REMOTE_JOBS=$(curl -s "$BASE_URL/jobs?work_mode=remote&page=1&limit=10")
if echo $REMOTE_JOBS | grep -q "success.*true"; then
    REMOTE_COUNT=$(echo $REMOTE_JOBS | jq '.data.total_count')
    echo "Remote jobs found: $REMOTE_COUNT"
    pass
else
    fail "Job search by work mode failed"
fi

# Test 22: Search Jobs by Location
run_test "Search Jobs by Location"
LOCATION_JOBS=$(curl -s "$BASE_URL/jobs?location=Remote&page=1&limit=10")
if echo $LOCATION_JOBS | grep -q "success.*true"; then
    pass
else
    fail "Job search by location failed"
fi

# Test 23: List Technologies
run_test "List Technologies"
TECHS=$(curl -s "$BASE_URL/technologies")
if echo $TECHS | grep -q "success.*true"; then
    pass
else
    fail "List technologies failed"
fi

# Test 24: List Programming Languages
run_test "List Programming Languages"
LANGS=$(curl -s "$BASE_URL/languages")
if echo $LANGS | grep -q "success.*true"; then
    pass
else
    fail "List languages failed"
fi

# Test 25: Invalid Token Test
run_test "Invalid Token Handling"
INVALID=$(curl -s -X GET $BASE_URL/auth/me \
  -H "Authorization: Bearer invalid_token_here")

if echo $INVALID | grep -q "Invalid or expired token"; then
    pass
else
    fail "Invalid token not handled properly"
fi

# Test 26: Missing Authorization Header
run_test "Missing Authorization Header"
NO_AUTH=$(curl -s -X POST $BASE_URL/developers \
  -H "Content-Type: application/json" \
  -d '{"bio":"test"}')

if echo $NO_AUTH | grep -q "Authorization header required"; then
    pass
else
    fail "Missing auth header not handled properly"
fi

# Summary
echo "=========================================="
echo "           TEST SUMMARY"
echo "=========================================="
echo -e "Total Tests:  $test_count"
echo -e "${GREEN}Passed:       $pass_count${NC}"
echo -e "${RED}Failed:       $fail_count${NC}"
echo ""

if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi

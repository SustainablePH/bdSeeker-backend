#!/bin/bash
set +e

BASE_URL="http://localhost:9000/api/v1"
PASSED=0
FAILED=0

echo "============================================"
echo "  COMPREHENSIVE API TEST SUITE - GIN FRAMEWORK"
echo "============================================"
echo ""

# Test 1: Health Check
echo "Test 1: Health Check"
RESPONSE=$(curl -s "$BASE_URL/health")
if echo "$RESPONSE" | grep -q '"status":"ok"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 2: Register Developer
echo "Test 2: Register Developer User"
DEV_REGISTER=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@test.com","password":"pass123","full_name":"Dev User","role":"developer"}')
if echo "$DEV_REGISTER" | grep -q '"token"'; then
    echo "✓ PASS"
    ((PASSED++))
    DEV_TOKEN=$(echo "$DEV_REGISTER" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
else
    echo "✗ FAIL"
    ((FAILED++))
    DEV_TOKEN=""
fi
echo ""

# Test 3: Register Company
echo "Test 3: Register Company User"
COMPANY_REGISTER=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"company@test.com","password":"pass123","full_name":"Company User","role":"company"}')
if echo "$COMPANY_REGISTER" | grep -q '"token"'; then
    echo "✓ PASS"
    ((PASSED++))
    COMPANY_TOKEN=$(echo "$COMPANY_REGISTER" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
else
    echo "✗ FAIL"
    ((FAILED++))
    COMPANY_TOKEN=""
fi
echo ""

# Test 4: Login Developer
echo "Test 4: Login Developer"
LOGIN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@test.com","password":"pass123"}')
if echo "$LOGIN" | grep -q '"token"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 5: Get Current User (Protected)
echo "Test 5: Get Current User (/auth/me)"
if [ -n "$DEV_TOKEN" ]; then
    ME=$(curl -s -X GET "$BASE_URL/auth/me" -H "Authorization: Bearer $DEV_TOKEN")
    if echo "$ME" | grep -q 'dev@test.com'; then
        echo "✓ PASS"
        ((PASSED++))
    else
        echo "✗ FAIL"
        ((FAILED++))
    fi
else
    echo "✗ FAIL: No token"
    ((FAILED++))
fi
echo ""

# Test 6: Create Developer Profile
echo "Test 6: Create Developer Profile (Protected)"
if [ -n "$DEV_TOKEN" ]; then
    DEV_PROFILE=$(curl -s -X POST "$BASE_URL/developers" \
      -H "Authorization: Bearer $DEV_TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"bio":"Experienced developer"}')
    if echo "$DEV_PROFILE" | grep -q '"message"'; then
        echo "✓ PASS"
        ((PASSED++))
    else
        echo "✗ FAIL"
        ((FAILED++))
    fi
else
    echo "✗ FAIL: No token"
    ((FAILED++))
fi
echo ""

# Test 7: List Developers
echo "Test 7: List Developers (Public)"
DEVS=$(curl -s "$BASE_URL/developers")
if echo "$DEVS" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 8: Create Company Profile
echo "Test 8: Create Company Profile (Protected)"
if [ -n "$COMPANY_TOKEN" ]; then
    COMPANY_PROFILE=$(curl -s -X POST "$BASE_URL/companies" \
      -H "Authorization: Bearer $COMPANY_TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"company_name":"Tech Corp","description":"Leading tech company","location":"Dhaka"}')
    if echo "$COMPANY_PROFILE" | grep -q '"message"'; then
        echo "✓ PASS"
        ((PASSED++))
    else
        echo "✗ FAIL"
        ((FAILED++))
    fi
else
    echo "✗ FAIL: No token"
    ((FAILED++))
fi
echo ""

# Test 9: List Companies
echo "Test 9: List Companies (Public)"
COMPANIES=$(curl -s "$BASE_URL/companies")
if echo "$COMPANIES" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 10: List Companies with Pagination
echo "Test 10: List Companies with Pagination"
COMPANIES_PAGE=$(curl -s "$BASE_URL/companies?page=1&limit=5")
if echo "$COMPANIES_PAGE" | grep -q '"total_pages"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 11: List Jobs
echo "Test 11: List Jobs (Public)"
JOBS=$(curl -s "$BASE_URL/jobs")
if echo "$JOBS" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 12: List Technologies
echo "Test 12: List Technologies"
TECH=$(curl -s "$BASE_URL/technologies")
if echo "$TECH" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 13: List Languages
echo "Test 13: List Languages"
LANGS=$(curl -s "$BASE_URL/languages")
if echo "$LANGS" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 14: Admin Login
echo "Test 14: Admin Login"
ADMIN_LOGIN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@bdseeker.com","password":"admin123"}')
if echo "$ADMIN_LOGIN" | grep -q '"token"'; then
    echo "✓ PASS"
    ((PASSED++))
    ADMIN_TOKEN=$(echo "$ADMIN_LOGIN" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
else
    echo "✗ FAIL"
    ((FAILED++))
    ADMIN_TOKEN=""
fi
echo ""

# Test 15: Get Admin Stats
echo "Test 15: Get Admin Stats (Admin Only)"
if [ -n "$ADMIN_TOKEN" ]; then
    STATS=$(curl -s -X GET "$BASE_URL/admin/stats" \
      -H "Authorization: Bearer $ADMIN_TOKEN")
    if echo "$STATS" | grep -q '"total_users"'; then
        echo "✓ PASS"
        ((PASSED++))
    else
        echo "✗ FAIL"
        ((FAILED++))
    fi
else
    echo "✗ FAIL: No admin token"
    ((FAILED++))
fi
echo ""

# Test 16: List Users (Admin)
echo "Test 16: List Users (Admin Only)"
if [ -n "$ADMIN_TOKEN" ]; then
    USERS=$(curl -s -X GET "$BASE_URL/admin/users" \
      -H "Authorization: Bearer $ADMIN_TOKEN")
    if echo "$USERS" | grep -q '"message"'; then
        echo "✓ PASS"
        ((PASSED++))
    else
        echo "✗ FAIL"
        ((FAILED++))
    fi
else
    echo "✗ FAIL: No admin token"
    ((FAILED++))
fi
echo ""

# Test 17: Unauthorized Access (should fail)
echo "Test 17: Unauthorized Access to Protected Route"
UNAUTH=$(curl -s -X GET "$BASE_URL/auth/me")
if echo "$UNAUTH" | grep -q '"error"'; then
    echo "✓ PASS (Correctly rejected)"
    ((PASSED++))
else
    echo "✗ FAIL (Should have been rejected)"
    ((FAILED++))
fi
echo ""

# Test 18: Invalid Login
echo "Test 18: Invalid Login Credentials"
INVALID=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"wrong@test.com","password":"wrongpass"}')
if echo "$INVALID" | grep -q '"error"'; then
    echo "✓ PASS (Correctly rejected)"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 19: CORS Headers
echo "Test 19: CORS Headers Present"
CORS=$(curl -s -I "$BASE_URL/health" | grep -i "access-control")
if [ -n "$CORS" ]; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Test 20: Logout
echo "Test 20: Logout"
LOGOUT=$(curl -s -X POST "$BASE_URL/auth/logout")
if echo "$LOGOUT" | grep -q 'Logged out'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL"
    ((FAILED++))
fi
echo ""

# Summary
echo "============================================"
echo "              TEST SUMMARY"
echo "============================================"
echo "Total Tests:  $((PASSED + FAILED))"
echo "Passed:       $PASSED"
echo "Failed:       $FAILED"
echo ""
if [ $FAILED -eq 0 ]; then
    echo "✅ ALL TESTS PASSED!"
    echo ""
    echo "Migration to Gin Framework: SUCCESS"
    echo "All APIs working correctly!"
else
    echo "⚠️  Some tests failed"
fi
echo "============================================"

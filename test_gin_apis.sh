#!/bin/bash
set +e  # Continue on errors

BASE_URL="http://localhost:9000/api/v1"
PASSED=0
FAILED=0

echo "=========================================="
echo "  Gin Framework API - Complete Test Suite"
echo "=========================================="
echo ""

# Test 1: Health Check
echo "Test 1: Health Check"
RESPONSE=$(curl -s "$BASE_URL/health")
if echo "$RESPONSE" | grep -q '"status":"ok"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL: $RESPONSE"
    ((FAILED++))
fi
echo ""

# Test 2: Register New User
echo "Test 2: Register Developer User"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_dev@example.com",
    "password": "password123",
    "full_name": "Test Developer",
    "role": "developer"
  }')
if echo "$REGISTER_RESPONSE" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
    TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
else
    echo "✗ FAIL: $REGISTER_RESPONSE"
    ((FAILED++))
    TOKEN=""
fi
echo ""

# Test 3: Login
echo "Test 3: Login with Registered User"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_dev@example.com",
    "password": "password123"
  }')
if echo "$LOGIN_RESPONSE" | grep -q '"token"'; then
    echo "✓ PASS"
    ((PASSED++))
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
else
    echo "✗ FAIL: $LOGIN_RESPONSE"
    ((FAILED++))
fi
echo ""

# Test 4: Get Current User (Protected Route)
echo "Test 4: Get Current User (/auth/me - Protected)"
if [ -n "$TOKEN" ]; then
    ME_RESPONSE=$(curl -s -X GET "$BASE_URL/auth/me" \
      -H "Authorization: Bearer $TOKEN")
    if echo "$ME_RESPONSE" | grep -q 'test_dev@example.com'; then
        echo "✓ PASS"
        ((PASSED++))
    else
        echo "✗ FAIL: $ME_RESPONSE"
        ((FAILED++))
    fi
else
    echo "✗ FAIL: No token available"
    ((FAILED++))
fi
echo ""

# Test 5: List Technologies (Public)
echo "Test 5: List Technologies (Public)"
TECH_RESPONSE=$(curl -s "$BASE_URL/technologies")
if echo "$TECH_RESPONSE" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL: $TECH_RESPONSE"
    ((FAILED++))
fi
echo ""

# Test 6: List Languages (Public)
echo "Test 6: List Languages (Public)"
LANG_RESPONSE=$(curl -s "$BASE_URL/languages")
if echo "$LANG_RESPONSE" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL: $LANG_RESPONSE"
    ((FAILED++))
fi
echo ""

# Test 7: List Companies (Public)
echo "Test 7: List Companies (Public)"
COMPANIES_RESPONSE=$(curl -s "$BASE_URL/companies")
if echo "$COMPANIES_RESPONSE" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL: $COMPANIES_RESPONSE"
    ((FAILED++))
fi
echo ""

# Test 8: List Jobs (Public)
echo "Test 8: List Jobs (Public)"
JOBS_RESPONSE=$(curl -s "$BASE_URL/jobs")
if echo "$JOBS_RESPONSE" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL: $JOBS_RESPONSE"
    ((FAILED++))
fi
echo ""

# Test 9: List Developers (Public)
echo "Test 9: List Developers (Public)"
DEVS_RESPONSE=$(curl -s "$BASE_URL/developers")
if echo "$DEVS_RESPONSE" | grep -q '"message"'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL: $DEVS_RESPONSE"
    ((FAILED++))
fi
echo ""

# Test 10: Logout
echo "Test 10: Logout"
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/logout")
if echo "$LOGOUT_RESPONSE" | grep -q 'Logged out successfully'; then
    echo "✓ PASS"
    ((PASSED++))
else
    echo "✗ FAIL: $LOGOUT_RESPONSE"
    ((FAILED++))
fi
echo ""

# Summary
echo "=========================================="
echo "           TEST SUMMARY"
echo "=========================================="
echo "Total Tests:  $((PASSED + FAILED))"
echo "Passed:       $PASSED"
echo "Failed:       $FAILED"
echo ""
if [ $FAILED -eq 0 ]; then
    echo "✓ All tests passed!"
else
    echo "✗ Some tests failed"
fi
echo "=========================================="

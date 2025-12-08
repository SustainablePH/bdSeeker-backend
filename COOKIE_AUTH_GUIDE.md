# 🍪 Cookie-Based Authentication Guide

## Overview

The bdSeeker API now supports **automatic cookie-based authentication** for seamless frontend integration. No need to manually manage tokens in your frontend code!

## ✨ Features

✅ **HTTP-Only Cookies** - Secure, cannot be accessed by JavaScript  
✅ **Automatic Token Management** - Backend sets/clears cookies automatically  
✅ **Dual Support** - Works with both cookies (browsers) and Authorization header (API clients)  
✅ **Session Management** - Automatic logout on token expiration  
✅ **CSRF Protection** - SameSite cookie attribute  

---

## 🔐 How It Works

### 1. Login/Register
When a user logs in or registers, the backend automatically:
- Generates a JWT token
- Sets an HTTP-only cookie (`auth_token`)
- Returns the token in the response (for compatibility)

### 2. Authenticated Requests
The middleware checks for authentication in this order:
1. **Cookie** (`auth_token`) - Preferred for browser clients
2. **Authorization Header** - Fallback for API clients (Postman, mobile apps)

### 3. Logout
When a user logs out:
- All authentication cookies are cleared
- User is automatically logged out

### 4. Token Expiration
If the token expires:
- Cookie is automatically cleared
- Frontend receives 401 Unauthorized
- Frontend should redirect to login page

---

## 🚀 Frontend Integration

### React/Vue/Angular Example

```javascript
// Login - No need to store token!
async function login(email, password) {
  const response = await fetch('http://localhost:9000/api/v1/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // IMPORTANT: Include cookies
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json();
  
  if (data.success) {
    // Cookie is automatically set!
    // Redirect to dashboard
    window.location.href = '/dashboard';
  }
}

// Make authenticated requests
async function getProfile() {
  const response = await fetch('http://localhost:9000/api/v1/auth/me', {
    method: 'GET',
    credentials: 'include', // IMPORTANT: Include cookies
  });

  if (response.status === 401) {
    // Token expired or invalid
    // Redirect to login
    window.location.href = '/login';
    return;
  }

  const data = await response.json();
  return data.data;
}

// Logout
async function logout() {
  await fetch('http://localhost:9000/api/v1/auth/logout', {
    method: 'POST',
    credentials: 'include', // IMPORTANT: Include cookies
  });

  // Cookie is automatically cleared!
  // Redirect to login
  window.location.href = '/login';
}
```

### Axios Example

```javascript
import axios from 'axios';

// Configure axios to include cookies
axios.defaults.withCredentials = true;
axios.defaults.baseURL = 'http://localhost:9000/api/v1';

// Login
async function login(email, password) {
  try {
    const response = await axios.post('/auth/login', { email, password });
    // Cookie is automatically set!
    return response.data;
  } catch (error) {
    console.error('Login failed:', error);
  }
}

// Get profile
async function getProfile() {
  try {
    const response = await axios.get('/auth/me');
    return response.data;
  } catch (error) {
    if (error.response?.status === 401) {
      // Redirect to login
      window.location.href = '/login';
    }
  }
}

// Logout
async function logout() {
  await axios.post('/auth/logout');
  window.location.href = '/login';
}

// Add response interceptor for automatic redirect
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Automatically redirect to login on 401
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### Fetch API with Error Handling

```javascript
// Create a wrapper function for authenticated requests
async function authenticatedFetch(url, options = {}) {
  const response = await fetch(url, {
    ...options,
    credentials: 'include', // Always include cookies
  });

  // Handle 401 Unauthorized
  if (response.status === 401) {
    // Clear any local state
    localStorage.clear();
    sessionStorage.clear();
    
    // Redirect to login
    window.location.href = '/login';
    throw new Error('Session expired');
  }

  return response;
}

// Usage
async function fetchJobs() {
  const response = await authenticatedFetch('http://localhost:9000/api/v1/jobs');
  const data = await response.json();
  return data;
}
```

---

## 🔧 API Endpoints

### Authentication Endpoints

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
- Sets `auth_token` cookie (HTTP-only, 24 hours)
- Returns user data and token

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe",
  "role": "developer"
}
```

**Response:**
- Sets `auth_token` cookie (HTTP-only, 24 hours)
- Returns user data and token

#### Logout
```http
POST /api/v1/auth/logout
```

**Response:**
- Clears all authentication cookies
- Returns success message

#### Get Current User
```http
GET /api/v1/auth/me
```

**Response:**
- Returns current user data
- 401 if not authenticated

---

## 🍪 Cookie Details

### auth_token Cookie

| Property | Value |
|----------|-------|
| Name | `auth_token` |
| Path | `/` |
| MaxAge | 86400 seconds (24 hours) |
| HttpOnly | `true` (cannot be accessed by JavaScript) |
| Secure | `false` (set to `true` in production with HTTPS) |
| SameSite | `Lax` (CSRF protection) |

---

## 🔒 Security Features

### 1. HTTP-Only Cookies
- Cookies cannot be accessed by JavaScript
- Protects against XSS attacks
- Token is never exposed to client-side code

### 2. SameSite Protection
- Prevents CSRF attacks
- Cookies only sent with same-site requests

### 3. Automatic Expiration
- Tokens expire after 24 hours
- Expired tokens are automatically rejected
- Cookies are cleared on expiration

### 4. Secure Flag (Production)
- Set `Secure: true` in production
- Cookies only sent over HTTPS
- Prevents man-in-the-middle attacks

---

## 🌐 CORS Configuration

For frontend apps running on different domains/ports, ensure CORS is configured:

```go
// Backend CORS middleware already configured
// Allows credentials (cookies) from frontend
```

Frontend must set:
```javascript
credentials: 'include' // or withCredentials: true for axios
```

---

## 📱 Mobile App Integration

For mobile apps (React Native, Flutter), you can still use the Authorization header:

```javascript
// Store token in secure storage
import AsyncStorage from '@react-native-async-storage/async-storage';

async function login(email, password) {
  const response = await fetch('http://localhost:9000/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json();
  
  if (data.success) {
    // Store token for mobile apps
    await AsyncStorage.setItem('auth_token', data.data.token);
  }
}

async function makeAuthenticatedRequest(url) {
  const token = await AsyncStorage.getItem('auth_token');
  
  const response = await fetch(url, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (response.status === 401) {
    // Token expired, redirect to login
    await AsyncStorage.removeItem('auth_token');
    // Navigate to login screen
  }

  return response;
}
```

---

## 🧪 Testing with curl

### Login and save cookie
```bash
curl -c cookies.txt -X POST http://localhost:9000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@bdseeker.com","password":"admin123"}'
```

### Use cookie for authenticated request
```bash
curl -b cookies.txt -X GET http://localhost:9000/api/v1/auth/me
```

### Logout
```bash
curl -b cookies.txt -c cookies.txt -X POST http://localhost:9000/api/v1/auth/logout
```

---

## 🎯 Best Practices

### Frontend

1. **Always use `credentials: 'include'`** for authenticated requests
2. **Handle 401 responses** by redirecting to login
3. **Don't store tokens** in localStorage or sessionStorage (cookies are automatic!)
4. **Use HTTPS** in production
5. **Implement global error handling** for 401 responses

### Backend (Already Implemented)

✅ HTTP-only cookies  
✅ SameSite protection  
✅ Automatic cookie clearing on logout  
✅ Token expiration handling  
✅ Dual authentication support (cookies + header)  

---

## 🚨 Common Issues

### Issue: Cookies not being set
**Solution:** Ensure `credentials: 'include'` is set in fetch/axios

### Issue: 401 on every request
**Solution:** Check CORS configuration and credentials setting

### Issue: Cookies not sent with requests
**Solution:** Verify `withCredentials: true` (axios) or `credentials: 'include'` (fetch)

### Issue: CORS errors
**Solution:** Backend CORS middleware must allow credentials from your frontend origin

---

## 📚 Additional Resources

- [MDN: HTTP Cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies)
- [OWASP: Session Management](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [Fetch API Credentials](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch#sending_a_request_with_credentials_included)

---

## ✅ Summary

With cookie-based authentication:
- ✅ No manual token management in frontend
- ✅ Automatic session handling
- ✅ Better security (HTTP-only cookies)
- ✅ Seamless user experience
- ✅ Easy logout implementation
- ✅ Automatic redirect on session expiration

**Your frontend code is now simpler and more secure!** 🎉

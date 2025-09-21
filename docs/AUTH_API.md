# Authentication API Documentation

## Overview

HEPIC App Server v2 provides a comprehensive authentication system with JWT tokens, role-based access control, and user management capabilities.

## Features

- ✅ **JWT Authentication** - Secure token-based authentication
- ✅ **User Registration & Login** - Complete user lifecycle management
- ✅ **Password Security** - bcrypt hashing with salt
- ✅ **Role-Based Access** - Admin and user roles
- ✅ **Profile Management** - Update user information
- ✅ **Password Management** - Change passwords securely
- ✅ **User Administration** - Admin-only user management
- ✅ **Input Validation** - Comprehensive request validation
- ✅ **ClickHouse Integration** - Persistent user storage

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "role": "user"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1234567890,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "User registered successfully"
}
```

#### Login User
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-16T10:30:00Z",
    "user": {
      "id": 1234567890,
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "last_login": "2024-01-15T10:30:00Z"
    }
  },
  "message": "Login successful"
}
```

### Protected Endpoints (JWT Token Required)

#### Get Current User Info
```http
GET /api/v1/auth/me
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1234567890,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "last_login": "2024-01-15T10:30:00Z"
  }
}
```

#### Update User Profile
```http
PUT /api/v1/auth/profile
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "username": "john_doe_updated",
  "email": "john.updated@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1234567890,
    "username": "john_doe_updated",
    "email": "john.updated@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "Profile updated successfully"
}
```

#### Change Password
```http
POST /api/v1/auth/change-password
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

### Admin Endpoints (Admin Role Required)

#### Get Users List
```http
GET /api/v1/auth/users?page=1&per_page=10&role=user
Authorization: Bearer <admin_jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1234567890,
        "username": "john_doe",
        "email": "john@example.com",
        "role": "user",
        "is_active": true,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z",
        "last_login": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "per_page": 10,
    "total_pages": 1
  }
}
```

#### Get User Statistics
```http
GET /api/v1/auth/stats
Authorization: Bearer <admin_jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_users": 150,
    "active_users": 120,
    "admin_users": 5,
    "regular_users": 145,
    "new_users_today": 3
  }
}
```

## Authentication Flow

### 1. User Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "password123"
  }'
```

### 2. User Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "password": "password123"
  }'
```

### 3. Using JWT Token
```bash
# Save token from login response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Use token for protected endpoints
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/auth/me
```

## JWT Token Details

### Token Structure
```json
{
  "user_id": 1234567890,
  "username": "john_doe",
  "role": "user",
  "exp": 1705312200,
  "iat": 1705225800,
  "jti": "unique-token-id"
}
```

### Token Claims
- `user_id` - User's unique identifier
- `username` - User's username
- `role` - User's role (admin/user)
- `exp` - Token expiration timestamp
- `iat` - Token issued at timestamp
- `jti` - Unique token identifier

### Token Configuration
- **Algorithm**: HS256
- **Expiration**: 24 hours (configurable)
- **Secret**: Configurable via `JWT_SECRET` environment variable

## Security Features

### Password Requirements
- Minimum 6 characters
- Stored with bcrypt hashing
- Salt rounds: 10 (default)

### Input Validation
- Username: 3-50 characters, alphanumeric
- Email: Valid email format
- Password: Minimum 6 characters
- Role: Must be "admin" or "user"

### Role-Based Access Control
- **User Role**: Access to own profile and basic features
- **Admin Role**: Full access to user management and statistics

### JWT Security
- Signed with HMAC SHA-256
- Configurable secret key
- Token expiration
- Unique token IDs (JTI)

## Error Responses

### Validation Errors
```json
{
  "success": false,
  "error": "Invalid request body"
}
```

### Authentication Errors
```json
{
  "success": false,
  "error": "Invalid credentials"
}
```

### Authorization Errors
```json
{
  "success": false,
  "error": "Access denied - admin role required"
}
```

### Server Errors
```json
{
  "success": false,
  "error": "Internal server error"
}
```

## Database Schema

### Users Table (ClickHouse)
```sql
CREATE TABLE users (
    id UInt64,
    username String,
    email String,
    password String,
    role String,
    is_active UInt8,
    created_at DateTime,
    updated_at DateTime,
    last_login Nullable(DateTime)
) ENGINE = MergeTree()
ORDER BY (id)
SETTINGS index_granularity = 8192
```

## Configuration

### Environment Variables
```bash
# JWT Configuration
HEPIC_JWT_SECRET=your-super-secret-jwt-key-here
HEPIC_JWT_EXPIRE_HOURS=24

# Database Configuration
HEPIC_DATABASE_HOST=localhost
HEPIC_DATABASE_PORT=9000
HEPIC_DATABASE_DATABASE=hepic_analytics
HEPIC_DATABASE_USER=default
HEPIC_DATABASE_PASSWORD=
```

### Configuration File
```json
{
  "jwt": {
    "secret": "your-super-secret-jwt-key-here",
    "expire_hours": 24
  },
  "database": {
    "host": "localhost",
    "port": 9000,
    "database": "hepic_analytics",
    "user": "default",
    "password": ""
  }
}
```

## Usage Examples

### Complete Authentication Flow

```bash
#!/bin/bash

# 1. Register a new user
echo "Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')

echo "Registration response: $REGISTER_RESPONSE"

# 2. Login
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo "Login response: $LOGIN_RESPONSE"

# 3. Extract token (using jq if available)
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')

# 4. Use token for protected endpoint
echo "Getting user info..."
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/auth/me
```

### Admin Operations

```bash
# Create admin user (first user can be created manually in database)
# Then use admin token for user management

# Get all users
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  "http://localhost:8080/api/v1/auth/users?page=1&per_page=10"

# Get user statistics
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/v1/auth/stats
```

## Integration with Analytics API

The authentication system integrates seamlessly with the analytics API:

```bash
# Use JWT token for analytics endpoints
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/analytics/stats
```

## Troubleshooting

### Common Issues

1. **"Invalid token" error**
   - Check if token is properly formatted
   - Verify token hasn't expired
   - Ensure correct Authorization header format

2. **"Access denied" error**
   - Verify user has required role
   - Check if user account is active
   - Ensure JWT token is valid

3. **"User already exists" error**
   - Username or email already taken
   - Use different credentials

4. **"Invalid credentials" error**
   - Check username/password combination
   - Verify user account is active

### Debug Mode

Enable debug logging for detailed authentication information:

```bash
./hepic-app-server-v2 --log-level debug --log-format text
```

This will provide detailed logs of:
- JWT token validation
- User authentication attempts
- Database operations
- Error details

## Security Best Practices

1. **Use strong JWT secrets** - Generate cryptographically secure secrets
2. **Set appropriate token expiration** - Balance security vs usability
3. **Validate all inputs** - Server-side validation for all requests
4. **Use HTTPS in production** - Encrypt all communications
5. **Monitor authentication logs** - Track failed login attempts
6. **Regular password updates** - Encourage users to change passwords
7. **Role-based permissions** - Implement least privilege access

This authentication system provides a robust foundation for secure user management in the HEPIC App Server v2.

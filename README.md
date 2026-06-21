# NYEH-BACK

A secure, production-grade Go web API backend with Google SSO authentication, JWT token management, and Redis-based session caching.

## 📋 Overview

NYEH-BACK is a RESTful API server built with Go that provides:

- **Google SSO Authentication** - OAuth 2.0 integration for seamless login
- **JWT Token Management** - Access tokens with configurable TTL, refresh token rotation
- **Redis Session Cache** - Secure session storage with TTL management
- **API Documentation** - Swagger/OpenAPI 2.0 documentation
- **Production-Ready** - Configurable timeouts, middleware, and environment-based settings

## 🛠️ Tech Stack

| Category | Technology |
|----------|-----------|
| **Framework** | Gochi (chi/v5) - Lightweight HTTP router |
| **Authentication** | Google OAuth 2.0, JWT tokens |
| **Caching** | Redis (go-redis/v9) |
| **Documentation** | Swagger UI (swaggo) |
| **Security** | HTTPS, HTTP-only cookies, CSRF protection |

## 📁 Project Structure

```
nyeh-back/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── api/                 # API routers and setup
│   │   ├── api.go
│   │   └── v1/
│   │       └── v1.go        # API v1 endpoints
│   ├── handler/             # HTTP request handlers
│   │   ├── auth/            # Authentication handlers
│   │   ├── me/              # User profile handlers
│   │   └── health.go        # Health check handler
│   ├── domain/              # Business logic & domain models
│   │   ├── auth/            # Auth domain models
│   │   ├── me/              # User profile domain models
│   │   └── health_domain.go
│   ├── middleware/          # HTTP middleware (auth, logging)
│   ├── repository/          # Data access layer
│   ├── cache/               # Cache implementations
│   │   └── redis/           # Redis cache wrappers
│   ├── infra/               # Infrastructure dependencies
│   └── core/                # Core utilities (config, tokens, hashing)
├── docs/                    # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── .env                     # Environment configuration
├── go.mod
├── go.sum
└── Makefile
```

## 🔑 Key Features

### Authentication Flow
1. User redirects to `/auth/google/login`
2. Google OAuth redirects back to `/auth/google/callback`
3. Backend validates OAuth state, exchanges code for user info
4. Access token (JWT) and refresh token are generated
5. Refresh token is stored in Redis with TTL
6. JWT is sent as an HTTP-only cookie

### API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/healthCheck` | GET | No | API health status |
| `/auth/google/login` | GET | No | Redirect to Google OAuth |
| `/auth/google/callback` | GET | No | OAuth callback handler |
| `/auth/refresh` | POST | No | Refresh token rotation |
| `/me` | GET | Yes | Get current user info |
| `/me` | POST | Yes | Update user profile |
| `/swagger/*` | GET | No | Swagger UI documentation |

### Security Features
- **JWT Access Tokens** - Configurable TTL (default: minutes)
- **Refresh Token Rotation** - Stored in Redis with configurable TTL (default: hours)
- **HTTP-only Cookies** - Prevents XSS token theft
- **CSRF Protection** - OAuth state parameter validation
- **Email Whitelisting** - Only approved emails can authenticate

## 🚀 Getting Started

### Prerequisites
- Go 1.26.4+
- Redis server
- Google OAuth 2.0 credentials

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd nyeh-back
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables in `.env`:
```env
ENV=development
PORT=8080
DEBUG=true

# JWT Settings
JWT_AUTH_TOKEN=your-super-secret-jwt-key
TOKEN_DIGEST=sha256
TTL_ACCESS_TOKEN_MINUTES=30
TTL_REFRESH_TOKEN_HOURS=7

# Google OAuth (Get from Google Cloud Console)
GOOGLE_OAUTH_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret
GOOGLE_ALLOWED_EMAIL=user@example.com

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

4. Run the server:
```bash
make dev
```

The server will start on `http://localhost:8080`

### API Documentation
Once the server is running, visit:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Swagger JSON: `http://localhost:8080/swagger/doc.json`

## ⚙️ Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ENV` | No | Environment mode (development/production) |
| `PORT` | No | Server port (default: 8080) |
| `DEBUG` | No | Enable debug mode (default: false) |
| `JWT_AUTH_TOKEN` | Yes | Secret key for JWT signing |
| `TOKEN_DIGEST` | No | Hash algorithm (default: sha256) |
| `TTL_ACCESS_TOKEN_MINUTES` | No | Access token expiry (default: 30) |
| `TTL_REFRESH_TOKEN_HOURS` | No | Refresh token expiry (default: 7) |
| `GOOGLE_OAUTH_CLIENT_ID` | Yes | Google OAuth client ID |
| `GOOGLE_OAUTH_CLIENT_SECRET` | Yes | Google OAuth client secret |
| `GOOGLE_ALLOWED_EMAIL` | Yes | Whitelisted email for authentication |
| `REDIS_ADDR` | Yes | Redis connection string |
| `REDIS_PASSWORD` | No | Redis password (if required) |

## 📖 API Documentation

### Health Check
```http
GET /healthCheck
```
```json
{
  "status": "ok",
  "port": "Running on port 8080",
  "message": "API is fully operational"
}
```

### Google OAuth Login
```http
GET /auth/google/login
Status: 307 Temporary Redirect
```
Redirects to Google's OAuth consent screen.

### Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

## 🔧 Development

### Update Swagger Documentation
```bash
make dev  # Runs swag init and swag fmt automatically
```

### Run Tests
```bash
go test ./...
```

## 🔐 Security Considerations

1. **Never commit `.env`** - The file is in `.gitignore`
2. **Rotate JWT_SECRET** - Change `JWT_AUTH_TOKEN` for production
3. **Use HTTPS in production** - OAuth requires secure connections
4. **Rate limiting** - Consider adding rate limiting for production
5. **CSRF tokens** - State parameter validation prevents CSRF attacks

## 🤝 Contributing

1. Create a feature branch: `git checkout -b feature/AmazingFeature`
2. Commit your changes: `git commit -m 'Add some AmazingFeature'`
3. Push to the branch: `git push origin feature/AmazingFeature`
4. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built with [Chi](https://github.com/go-chi/chi) - Lightweight, stylish HTTP router
- Authentication powered by [Google OAuth 2.0](https://developers.google.com/identity/protocols/oauth2)
- Redis caching with [go-redis](https://github.com/redis/go-redis)

---

*For questions or support, please open an issue in the repository.*

# Sport Booking Backend API

A modern, high-performance REST API for sports venue booking system built with Go, Fiber, and GORM.

## 🚀 Features

- **Authentication & Authorization**: JWT-based auth with role-based access control
- **Venue Management**: CRUD operations for sports venues with categories
- **Booking System**: Real-time booking with conflict detection and validation
- **Search & Filtering**: Advanced search with pagination and sorting
- **Performance Optimized**: Connection pooling, efficient queries, and caching
- **Security**: Rate limiting, CORS, helmet security headers
- **Validation**: Comprehensive input validation with detailed error messages
- **Documentation**: Auto-generated API documentation
- **Monitoring**: Structured logging and health checks

## 📋 Prerequisites

- Go 1.19 or higher
- MySQL 8.0 or higher
- Git

## 🛠️ Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/adityaibrar/sport-booking-backend.git
   cd sport-booking-backend
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Set up database**

   ```sql
   CREATE DATABASE sport_booking;
   ```

5. **Run migrations**
   ```bash
   go run main.go
   # Migrations will run automatically on startup
   ```

## 🏃‍♂️ Running the Application

### Development Mode

```bash
go run main.go
```

### Production Mode

```bash
# Build the application
go build -o sport-booking-api

# Run the binary
./sport-booking-api
```

The API will be available at `http://localhost:8000`

## 📚 API Documentation

### Base URL

```
http://localhost:8000/api/v1
```

### Authentication Endpoints

| Method | Endpoint         | Description         |
| ------ | ---------------- | ------------------- |
| POST   | `/auth/register` | Register a new user |
| POST   | `/auth/login`    | Login user          |
| POST   | `/auth/logout`   | Logout user         |

### User Endpoints

| Method | Endpoint                     | Description              | Auth Required |
| ------ | ---------------------------- | ------------------------ | ------------- |
| GET    | `/user/venues`               | Get all venues           | ✅            |
| GET    | `/user/venues/:id`           | Get venue details        | ✅            |
| POST   | `/user/venues/availability`  | Check venue availability | ✅            |
| POST   | `/user/bookings`             | Create new booking       | ✅            |
| GET    | `/user/bookings/history/:id` | Get user booking history | ✅            |
| POST   | `/user/bookings/:id/cancel`  | Cancel booking           | ✅            |

### Admin Endpoints

| Method | Endpoint                     | Description                   | Auth Required | Admin Only |
| ------ | ---------------------------- | ----------------------------- | ------------- | ---------- |
| GET    | `/admin/dashboard`           | Get dashboard analytics       | ✅            | ✅         |
| GET    | `/admin/venues`              | Get all venues (admin view)   | ✅            | ✅         |
| POST   | `/admin/venues`              | Create new venue              | ✅            | ✅         |
| PUT    | `/admin/venues/:id`          | Update venue                  | ✅            | ✅         |
| DELETE | `/admin/venues/:id`          | Delete venue                  | ✅            | ✅         |
| GET    | `/admin/bookings`            | Get all bookings with filters | ✅            | ✅         |
| PUT    | `/admin/bookings/:id/status` | Update booking status         | ✅            | ✅         |

## 📊 Request & Response Examples

### Create Booking

**Request:**

```bash
POST /api/v1/user/bookings
Authorization: Bearer <token>
Content-Type: application/json

{
  "venue_id": 1,
  "start_time": "2024-08-15T10:00:00Z",
  "duration": 2,
  "notes": "Birthday party booking"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Booking created successfully",
  "data": {
    "id": 1,
    "start_time": "2024-08-15T10:00:00Z",
    "end_time": "2024-08-15T12:00:00Z",
    "duration": 2,
    "total_price": 200000,
    "status": "pending",
    "notes": "Birthday party booking",
    "created_at": "2024-08-08T15:30:00Z",
    "updated_at": "2024-08-08T15:30:00Z",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    },
    "venue": {
      "id": 1,
      "name": "Football Field A",
      "category": "football",
      "price_per_hour": 100000
    }
  },
  "timestamp": "2024-08-08T15:30:00Z"
}
```

### Get Bookings with Filters

**Request:**

```bash
GET /api/v1/admin/bookings?status=confirmed&start_date=2024-08-01&end_date=2024-08-31&page=1&limit=10
Authorization: Bearer <admin-token>
```

**Response:**

```json
{
  "success": true,
  "message": "Bookings retrieved successfully",
  "data": [...],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 45,
    "total_pages": 5,
    "has_next": true,
    "has_previous": false
  },
  "timestamp": "2024-08-08T15:30:00Z"
}
```

## 🏗️ Database Schema

### Users Table

```sql
CREATE TABLE users (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  phone VARCHAR(20),
  address VARCHAR(500),
  role ENUM('user', 'admin') DEFAULT 'user',
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);
```

### Venues Table

```sql
CREATE TABLE venues (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  category ENUM('football', 'basketball', 'badminton', 'tennis', 'volleyball', 'swimming', 'futsal', 'other') NOT NULL,
  price_per_hour DECIMAL(10,2) NOT NULL,
  description TEXT NOT NULL,
  location VARCHAR(500),
  facilities TEXT,
  images TEXT,
  is_active BOOLEAN DEFAULT TRUE,
  open_time VARCHAR(5) NOT NULL,
  close_time VARCHAR(5) NOT NULL,
  capacity INT DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);
```

### Bookings Table

```sql
CREATE TABLE bookings (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT NOT NULL,
  venue_id BIGINT NOT NULL,
  start_time TIMESTAMP NOT NULL,
  duration INT NOT NULL,
  total_price DECIMAL(10,2) NOT NULL,
  status ENUM('pending', 'confirmed', 'completed', 'cancelled') DEFAULT 'pending',
  payment_qris VARCHAR(500),
  notes VARCHAR(1000),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (venue_id) REFERENCES venues(id)
);
```

## 🔧 Configuration

### Environment Variables

| Variable              | Description             | Default       |
| --------------------- | ----------------------- | ------------- |
| `HOST`                | Server host             | `0.0.0.0`     |
| `PORT`                | Server port             | `8000`        |
| `DATABASE_URL`        | MySQL connection string | Required      |
| `JWT_SECRET`          | JWT signing secret      | Required      |
| `JWT_EXPIRATION`      | JWT token expiration    | `24h`         |
| `ENVIRONMENT`         | Environment mode        | `development` |
| `DEBUG`               | Enable debug logging    | `true`        |
| `RATE_LIMIT_MAX`      | Rate limit max requests | `100`         |
| `RATE_LIMIT_DURATION` | Rate limit duration     | `1m`          |

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./controllers -v
```

## 📈 Performance Features

- **Connection Pooling**: Optimized database connections
- **Query Optimization**: Efficient GORM queries with proper indexing
- **Rate Limiting**: Prevents API abuse
- **Caching Ready**: Structure supports Redis integration
- **Pagination**: All list endpoints support pagination
- **Lazy Loading**: Associations loaded only when needed

## 🔐 Security Features

- **JWT Authentication**: Stateless authentication
- **Password Hashing**: Bcrypt password hashing
- **Rate Limiting**: IP-based rate limiting
- **CORS Protection**: Configurable CORS settings
- **Security Headers**: Helmet middleware for security headers
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: GORM provides built-in protection

## 🚀 Deployment

### Docker Deployment

```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Build & Run

```bash
docker build -t sport-booking-api .
docker run -p 8000:8000 sport-booking-api
```

## 📝 API Response Format

All API responses follow a consistent format:

```json
{
  "success": boolean,
  "message": "string",
  "data": "object|array|null",
  "error": "object|null",
  "meta": "object|null",
  "timestamp": "ISO8601 datetime"
}
```

### Error Response Example

```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "errors": [
      {
        "code": "required",
        "field": "venue_id",
        "message": "venue_id is required"
      }
    ]
  },
  "timestamp": "2024-08-08T15:30:00Z"
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you have any questions or need help with setup, please create an issue in the GitHub repository.

## 🔄 Version History

- **v1.0.0** - Initial release with core booking functionality
  - User authentication and authorization
  - Venue management
  - Booking system with conflict detection
  - Admin dashboard
  - Comprehensive API documentation

---

**Developed with ❤️ by Aditya Ibrar**

# Services MicroService

A Service Go MicroService for managing services and their versions in an organization.

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set environment variables (optional):
   ```bash
   export DATABASE_URL="postgres://user:password@localhost/services_db"
   export PORT=8080
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

---

**Alternatively, you can use Docker Compose to start both the API and the Postgres database:**

```bash
docker-compose up
```

This will start both the API service and the Postgres database as defined in `docker-compose.yaml`.

**To load sample data into Postgres for testing the application:**

```bash
make db-load
```

This will execute the `sample_data.sql` script and populate the database with sample services and versions.

## API Endpoints

### Services

#### List Services

```
GET /api/v1/services?name=search&page=1&limit=10&sort=name&order=asc
```

Response:

```json
{
  "services": [
    {
      "id": 1,
      "name": "User Service",
      "description": "Manages user authentication and profiles",
      "versionCount": 3,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 50,
    "items_per_page": 10
  }
}
```

#### Get Service

```
GET /api/v1/services/:sid
```

Response:

```json
{
  "id": 1,
  "name": "User Service",
  "description": "Manages user authentication and profiles",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "versions": [
    {
      "id": 1,
      "service_id": 1,
      "version": "1.0.0",
      "description": "Initial release",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### Create Service

```
POST /api/v1/services
```

Request:

```json
{
  "name": "User Service",
  "description": "Manages user authentication and profiles"
}
```

Response: `201 Created`

```json
{
  "id": 1,
  "name": "User Service",
  "description": "Manages user authentication and profiles",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Update Service

```
PATCH /api/v1/services/:sid
```

Request:

```json
{
  "name": "User Authentication Service",
  "description": "Updated description"
}
```

Response: `200 OK`

```json
{
  "id": 1,
  "name": "User Authentication Service",
  "description": "Updated description",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Delete Service

```
DELETE /api/v1/services/:sid
```

Response: `204 No Content`

### Versions

#### Create Version

```
POST /api/v1/services/:sid/versions
```

Request:

```json
{
  "version": "1.0.0",
  "description": "Initial release"
}
```

Response: `201 Created`

```json
{
  "id": 1,
  "service_id": 1,
  "version": "1.0.0",
  "description": "Initial release",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Get Version

```
GET /api/v1/services/:sid/versions/:vid
```

Response: `200 OK`

```json
{
  "id": 1,
  "service_id": 1,
  "version": "1.0.0",
  "description": "Initial release",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Update Version

```
PUT /api/v1/services/:sid/versions/:vid
```

Request:

```json
{
  "version": "1.0.1",
  "description": "Bug fixes"
}
```

Response: `200 OK`

```json
{
  "id": 1,
  "service_id": 1,
  "version": "1.0.1",
  "description": "Bug fixes",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Delete Version

```
DELETE /api/v1/services/:sid/versions/:vid
```

Response: `204 No Content`

### Error Responses

All endpoints may return the following error responses:

#### 400 Bad Request

```json
{
  "code": "invalid_request_body",
  "message": "Invalid request body",
  "details": "error details"
}
```

#### 404 Not Found

```json
{
  "code": "service_not_found",
  "message": "Service with ID 123 could not be found"
}
```

```json
{
  "code": "version_not_found",
  "message": "Version not found",
  "details": "The requested version does not exist"
}
```

#### 500 Internal Server Error

```json
{
  "code": "internal_error",
  "message": "An error occurred while processing the request",
  "details": "error details"
}
```

### Health Check

```
GET /health
```

Response: `200 OK`

```json
{
  "status": "ok",
  "version": "1.0.0",
  "components": {
    "database": "ok"
  }
}
```

## Database

The API uses PostgreSQL for storing and persisting the user service configuration.

- If `DATABASE_URL` is set, it will connect to the specified PostgreSQL database

## Testing

To run tests:

```bash
go test ./...
```

## Required Future Enhancements

- **Authentication and Authorization**: Implement Auth Middleware to support the authentication and authorization for all apis
- **Caching Strategies**: Add Redis or in-memory caching for frequently accessed services and versions based on userId from authMiddleware
- **Rate Limiting**: Implement request throttling to prevent DDOS attacks or API misuse strategy based on the requirements/micro-service usecase
- **Request Validation**: Add input validation using a validation for handling various scenarios
- **Observability & Monitoring**: Integrate with Prometheus and Grafana for metrics collection, alerting, and performance tracking
- **Integration Tests**: Create end-to-end tests that verify the complete API functionality by spinning up a mock db and clean up after

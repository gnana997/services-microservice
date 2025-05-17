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

### List Services

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

### Get Service

```
GET /api/v1/services/:id
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

### Get Service Version

```
GET api/v1/service/{id}/versions/{id}
```

Response:

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

### Health Check

```
GET /health
```

Response:

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

## Future Enhancements

- Authentication and authorization
- CRUD operations for services and versions
- Rate limiting
- Request validation
- Integration tests
- Docker support
- CI/CD pipeline

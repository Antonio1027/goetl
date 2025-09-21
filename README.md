
# GoETL: Ads & CRM ETL Backend

GoETL is a modular backend service written in Go for extracting, transforming, and loading (ETL) Ads and CRM data. It exposes a REST API for metrics, supports MongoDB persistence.

## Features
- Extracts data from Ads and CRM public APIs
- Transforms and aggregates metrics (clicks, cost, leads, revenue, ROAS, etc.)
- REST API endpoints for ETL and metrics
- MongoDB integration for persistence
- Modular, testable codebase
- Makefile automation
- Docker & Docker Compose support

---

## Requirements

### Local Development
- **Go** (1.25)
- **MongoDB** (7) Configuired with docker compose
- **Make** (for automation commands)
- **Git** (To download the repository)
- **Docker** (You can install docker desktop)

---

## Installation & Setup

### 1. Clone the Repository
```sh
git clone git@github.com:Antonio1027/goetl.git
cd goetl
```

### 2. Environment Configuration
Copy the example environment file and edit as needed:
```sh
cp .env.example .env
```

**Environment variables:**
- `ADS_API_URL`
- `CRM_API_URL`
- `SINK_URL`
- `SINK_SECRET`
- `PORT`
- `MONGO_HOST`
- `MONGO_PORT`
- `MONGO_USERNAME`
- `MONGO_PASSWORD`
- `MONGO_DATABASE`


### 3. Start Locally

Command to build and run the service:

```sh
make run
```

The API will be available at `http://localhost:8080` (default).

---

## Testing

Command to run all tests:
```sh
make tests
```

---

## Project Structure

```
cmd/                # Main entrypoint
internal/
	api/              # API routes and server
	clients/          # Ads & CRM API clients
	db/               # MongoDB helpers
	etl/              # ETL logic
	models/           # Data models
	utils/            # Utility functions
Makefile            # Automation commands
Dockerfile          # Container build
Dockerfile.multistage # Multi-stage Docker build
.env.example        # Example environment config
docker-compose.yml  # Service orchestration
```

---


## OpenAPI Documentation

This project have swagger documentation in this URL:

```
http://localhost:8081/#
```

## Examples to execute api calls with curl


### Endpoint to start the ETL
```
curl --location --request POST 'http://localhost:8080/ingest/run?since=2025-01-01'
```

### Endpoint to get metrics by channel

```
curl --location 'http://localhost:8080/metrics/channel?from=2025-08-08&to=2025-08-08&channel=facebook_ads&limit=10&offset=0'
```

### Endpoint to get metrics by campaign

```
curl --location 'http://localhost:8080/metrics/campaign?from=2025-08-08&to=2025-08-08&utm_campaign=back_to_school&limit=2&offset=0'
```

## Contributing
Pull requests and issues are welcome!

---

## License
MIT

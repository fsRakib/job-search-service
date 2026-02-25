# Job Search Service

A microservice-based job search application built with **Go**, **gRPC**, and **Elasticsearch**.

## 📋 Overview

This service provides a high-performance job search system with the following operations:

- **CreateJob** - Add new job listings
- **SearchJobs** - Search jobs by title, description, company, location, and skills
- **GetJob** - Retrieve job details by ID
- **DeleteJob** - Remove job listings

## 🏗️ Architecture

```
Client → gRPC → Job Search Service → Elasticsearch
```

### Layers:

- **gRPC Layer** - Handles client requests
- **Service Layer** - Business logic
- **Repository Layer** - Elasticsearch operations
- **Elasticsearch** - Data storage and search engine

## 📁 Project Structure

```
job-search-service/
├── cmd/
│   ├── server/           # Main server application
│   │   └── main.go
│   └── client/           # Test client
│       └── test_client.go
├── internal/
│   ├── grpc/            # gRPC handlers
│   │   └── job_handler.go
│   ├── service/         # Business logic
│   │   └── job_service.go
│   ├── repository/      # Elasticsearch operations
│   │   └── elastic_repository.go
│   ├── elastic/         # Elasticsearch client
│   │   └── elastic_client.go
│   └── models/          # Data models
│       └── job.go
├── proto/               # Protocol buffers
│   ├── job.proto
│   ├── job.pb.go
│   └── job_grpc.pb.go
├── configs/
│   └── config.yaml      # Configuration file
├── bin/
│   └── server          # Compiled binary
├── go.mod
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Elasticsearch 8.x
- protoc (Protocol Buffer Compiler)

### Installation

1. **Clone and navigate to the project:**

   ```bash
   cd job-search-service
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Start Elasticsearch:**

   ```bash
   sudo systemctl start elasticsearch
   ```

4. **Verify Elasticsearch is running:**
   ```bash
   curl http://localhost:9200
   ```

### Building

```bash
go build -o bin/server cmd/server/main.go
```

### Running the Server

```bash
./bin/server
```

Expected output:

```
2026/02/25 11:16:00 Starting Job Search Service...
2026/02/25 11:16:00 Successfully connected to Elasticsearch
2026/02/25 11:16:00 Index 'jobs' created successfully
2026/02/25 11:16:00 gRPC server listening on port 50051
```

## 🧪 Testing

### Run the Test Client

```bash
go run cmd/client/test_client.go
```

This will:

1. Create 2 job listings
2. Search for jobs by query
3. Get job details by ID
4. Search by location
5. Search by skills
6. Delete a job
7. Verify deletion

### Manual Testing with grpcurl

Install grpcurl:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

**Create a Job:**

```bash
grpcurl -plaintext -d '{
  "title": "DevOps Engineer",
  "description": "Manage cloud infrastructure",
  "company": "CloudTech",
  "location": "Dhaka",
  "skills": ["Docker", "Kubernetes", "AWS"],
  "salary": 90000
}' localhost:50051 job.JobService/CreateJob
```

**Search Jobs:**

```bash
grpcurl -plaintext -d '{
  "query": "Engineer"
}' localhost:50051 job.JobService/SearchJobs
```

**Get Job by ID:**

```bash
grpcurl -plaintext -d '{
  "id": "YOUR_JOB_ID"
}' localhost:50051 job.JobService/GetJob
```

**Delete Job:**

```bash
grpcurl -plaintext -d '{
  "id": "YOUR_JOB_ID"
}' localhost:50051 job.JobService/DeleteJob
```

## 📊 Data Model

### Job Structure

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "company": "string",
  "location": "string",
  "skills": ["string"],
  "salary": 0.0,
  "created_at": "2026-02-25T00:00:00Z"
}
```

## 🔍 Elasticsearch Integration

### Index Name

`jobs`

### Mapping

- **id**: keyword (exact matching)
- **title**: text (full-text search)
- **description**: text (full-text search)
- **company**: text (full-text search)
- **location**: text (full-text search)
- **skills**: keyword (exact matching)
- **salary**: double
- **created_at**: date

### Search Capabilities

1. **Full-text search** across title, description, and company
2. **Location filtering**
3. **Skills matching** (exact match on keyword field)
4. **Relevance scoring** - Results ranked by relevance

## ⚙️ Configuration

Edit `configs/config.yaml`:

```yaml
elasticsearch:
  url: http://localhost:9200
  index: jobs

server:
  port: 50051
```

## 🛠️ Development

### Regenerate Proto Files

After modifying `proto/job.proto`:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/job.proto
```

### Add Dependencies

```bash
go get <package>
go mod tidy
```

## 📝 API Reference

### CreateJob

Creates a new job listing.

**Request:**

```protobuf
message CreateJobRequest {
  string title = 1;
  string description = 2;
  string company = 3;
  string location = 4;
  repeated string skills = 5;
  double salary = 6;
}
```

**Response:**

```protobuf
message CreateJobResponse {
  string id = 1;
  string message = 2;
}
```

### SearchJobs

Searches jobs with optional filters.

**Request:**

```protobuf
message SearchJobsRequest {
  string query = 1;      // Search in title, description, company
  string location = 2;   // Filter by location
  repeated string skills = 3;  // Filter by skills
}
```

**Response:**

```protobuf
message SearchJobsResponse {
  repeated Job jobs = 1;
  int32 total = 2;
}
```

### GetJob

Retrieves a job by ID.

**Request:**

```protobuf
message GetJobRequest {
  string id = 1;
}
```

**Response:**

```protobuf
message GetJobResponse {
  Job job = 1;
}
```

### DeleteJob

Deletes a job by ID.

**Request:**

```protobuf
message DeleteJobRequest {
  string id = 1;
}
```

**Response:**

```protobuf
message DeleteJobResponse {
  string message = 1;
}
```

## 🔥 Features

- ✅ Fast full-text search using Elasticsearch
- ✅ gRPC for efficient communication
- ✅ RESTful-like operations (CRUD)
- ✅ Automatic UUID generation for jobs
- ✅ Real-time indexing
- ✅ Graceful shutdown
- ✅ Structured logging

## 🚧 Future Improvements

- [ ] Pagination support
- [ ] Advanced filtering (salary range, date range)
- [ ] Sorting options
- [ ] Autocomplete functionality
- [ ] Authentication & Authorization
- [ ] Rate limiting
- [ ] Metrics and monitoring
- [ ] Docker containerization
- [ ] Kubernetes deployment
- [ ] Unit and integration tests

## 📚 Resources

- [Go gRPC Documentation](https://grpc.io/docs/languages/go/)
- [Elasticsearch Go Client](https://github.com/elastic/go-elasticsearch)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)

## 📄 License

This project is for educational purposes.

## 👨‍💻 Author

Built as a learning project for microservices architecture with Go, gRPC, and Elasticsearch.

---

**Happy Searching! 🔍**

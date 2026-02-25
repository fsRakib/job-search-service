# Elasticsearch Workflow - Job Search Service

## 📊 Complete Data Flow

### 1. Service Initialization

```
Application Start
    ↓
Create Elasticsearch Client (elastic_client.go)
    ↓
Connect to Elasticsearch (http://localhost:9200)
    ↓
Create/Verify Index "jobs" exists
    ↓
Ready to accept requests
```

---

## 🔄 Operation Workflows

### CREATE JOB Operation

```
Client Request (gRPC CreateJob)
    ↓
job_handler.go (receives request)
    ↓
job_service.go (generate UUID, create Job model)
    ↓
elastic_repository.go (Create method)
    ↓
Elasticsearch IndexRequest API
    ↓
Job Document indexed in "jobs" index
    ↓
Response with Job ID returned to client
```

**Elasticsearch Operation:**

- **API**: `IndexRequest`
- **Index**: `jobs`
- **Document ID**: UUID
- **Action**: Store JSON document
- **Refresh**: `true` (immediate visibility)

**Example Document Stored:**

```json
{
  "id": "d6a2662a-0755-42bb-b242-2a9862d3e5bc",
  "title": "Backend Developer",
  "description": "Develop backend services using Go",
  "company": "TechNext",
  "location": "Dhaka",
  "skills": ["Go", "gRPC", "Elasticsearch"],
  "salary": 80000,
  "created_at": "2026-02-25T11:18:11Z"
}
```

---

### SEARCH JOBS Operation

```
Client Request (gRPC SearchJobs with query/location/skills)
    ↓
job_handler.go (receives search parameters)
    ↓
job_service.go (pass to repository)
    ↓
elastic_repository.go (Search method)
    ↓
Build Elasticsearch Query DSL
    ↓
Execute Search Request
    ↓
Elasticsearch analyzes query and searches index
    ↓
Returns matching documents with scores
    ↓
Parse results into Job models
    ↓
Return jobs array to client
```

**Elasticsearch Query Types:**

1. **Full-text Search (multi_match)**

   ```json
   {
     "multi_match": {
       "query": "Developer",
       "fields": ["title", "description", "company"]
     }
   }
   ```

   - Searches across multiple text fields
   - Uses text analysis (tokenization, stemming)
   - Returns relevance-scored results

2. **Location Filter (match)**

   ```json
   {
     "match": {
       "location": "Dhaka"
     }
   }
   ```

   - Matches location field
   - Text analysis applied

3. **Skills Filter (terms)**
   ```json
   {
     "terms": {
       "skills": ["Go", "Docker"]
     }
   }
   ```

   - Exact match on keyword field
   - No text analysis (case-sensitive)

**Combined Query Example:**

```json
{
  "query": {
    "bool": {
      "must": [
        {
          "multi_match": {
            "query": "Developer",
            "fields": ["title", "description", "company"]
          }
        },
        {
          "match": {
            "location": "Dhaka"
          }
        },
        {
          "terms": {
            "skills": ["Go"]
          }
        }
      ]
    }
  }
}
```

---

### GET JOB BY ID Operation

```
Client Request (gRPC GetJob with ID)
    ↓
job_handler.go
    ↓
job_service.go
    ↓
elastic_repository.go (GetByID method)
    ↓
Elasticsearch Get API (by document ID)
    ↓
Direct document retrieval (very fast, no search)
    ↓
Return Job to client
```

**Elasticsearch Operation:**

- **API**: `Get`
- **Lookup**: By document ID (primary key)
- **Speed**: O(1) - Direct retrieval, no indexing scan
- **Returns**: Single document or 404 error

---

### DELETE JOB Operation

```
Client Request (gRPC DeleteJob with ID)
    ↓
job_handler.go
    ↓
job_service.go
    ↓
elastic_repository.go (Delete method)
    ↓
Elasticsearch DeleteRequest API
    ↓
Document removed from index
    ↓
Success message returned to client
```

**Elasticsearch Operation:**

- **API**: `DeleteRequest`
- **Refresh**: `true` (immediate removal)
- **Returns**: 404 if document doesn't exist

---

## 🧠 Elasticsearch Core Concepts

### Index

- Like a database in SQL
- Stores related documents
- Our index: `jobs`

### Document

- Like a row in SQL
- JSON format
- Our document: One job posting

### Mapping

- Like a schema in SQL
- Defines field types
- Example:
  ```json
  {
    "properties": {
      "id": { "type": "keyword" },
      "title": { "type": "text" },
      "skills": { "type": "keyword" },
      "salary": { "type": "double" }
    }
  }
  ```

### Field Types

1. **text** - Full-text search
   - Analyzed (tokenized, lowercased)
   - Used for: title, description, company, location
   - Example: "Backend Developer" → ["backend", "developer"]

2. **keyword** - Exact matching
   - Not analyzed
   - Used for: id, skills
   - Example: "Go" → "Go" (exact)

3. **double** - Numeric
   - Used for: salary
   - Supports range queries

4. **date** - Timestamp
   - Used for: created_at
   - Supports date range queries

---

## 🔍 Search Features

### Relevance Scoring

Elasticsearch ranks results by relevance:

- **TF-IDF**: Term frequency × Inverse document frequency
- **BM25**: Default ranking algorithm
- More relevant documents score higher

### Analyzers

Text fields are analyzed:

1. **Tokenization**: "Backend Developer" → ["Backend", "Developer"]
2. **Lowercasing**: ["Backend", "Developer"] → ["backend", "developer"]
3. **Stemming**: ["developer"] → ["develop"]

This enables:

- Case-insensitive search
- Partial word matching
- Synonym handling (if configured)

### Query Performance

- **Get by ID**: ~1ms (direct lookup)
- **Simple search**: ~10-50ms
- **Complex search**: ~50-200ms
- **Scales to millions of documents**

---

## 📦 Data Storage

### Physical Storage

- **Location**: `/var/lib/elasticsearch/`
- **Index data**: `/var/lib/elasticsearch/nodes/0/indices/`
- **Format**: Apache Lucene segments

### Persistence

- Data persists to disk
- Survives service restarts
- Transaction log for durability

### Refresh

- **refresh=true**: Makes changes immediately searchable
- Default: 1 second refresh interval
- Trade-off: Performance vs. real-time visibility

---

## 🚀 Why Elasticsearch for Job Search?

### Advantages

1. **Full-text Search**
   - Search "golang developer" finds "Go Developer"
   - Fuzzy matching handles typos
   - Relevance ranking

2. **Fast**
   - Inverted index structure
   - Sub-second search on millions of records
   - Distributed and scalable

3. **Flexible Queries**
   - Combine filters (location + skills + query)
   - Range queries (salary $50k-$80k)
   - Aggregations (count by location)

4. **No Complex Joins**
   - Document-oriented
   - Denormalized data = fast reads

### Comparison with SQL

| Feature           | Elasticsearch | SQL Database           |
| ----------------- | ------------- | ---------------------- |
| Full-text search  | ✅ Native     | ❌ LIKE queries (slow) |
| Fuzzy matching    | ✅ Built-in   | ❌ Complex             |
| Relevance scoring | ✅ Automatic  | ❌ Manual              |
| Scalability       | ✅ Horizontal | ⚠️ Vertical            |
| Complex joins     | ❌ Limited    | ✅ Strong              |
| ACID transactions | ⚠️ Eventual   | ✅ Immediate           |

---

## 🎯 Use Cases Perfect for Elasticsearch

1. ✅ **Job Search** (our project)
2. ✅ E-commerce product search
3. ✅ Log analysis
4. ✅ Content management
5. ✅ Autocomplete/suggestions
6. ❌ Banking transactions (need ACID)
7. ❌ Complex relational data

---

## 🔧 Repository Implementation Breakdown

### elastic_repository.go Functions

1. **`Create(ctx, job)`**
   - Marshals Job to JSON
   - Creates IndexRequest with job ID
   - Sends to Elasticsearch
   - Returns error if failed

2. **`Search(ctx, query, location, skills)`**
   - Builds bool query with must clauses
   - Constructs multi_match for text search
   - Adds match query for location
   - Adds terms query for skills
   - Executes search
   - Unmarshals results to Job array

3. **`GetByID(ctx, id)`**
   - Direct Get request by document ID
   - Returns 404 if not found
   - Unmarshals to Job model

4. **`Delete(ctx, id)`**
   - Creates DeleteRequest
   - Returns 404 if not found
   - Confirms deletion

---

## 📈 Performance Considerations

### Indexing Speed

- Single document: ~5-10ms
- Bulk indexing: 1000s per second
- Trade-off: refresh interval

### Search Speed

- Simple query: 10-50ms
- Complex aggregations: 100-500ms
- Caching helps repeated queries

### Optimization Tips

1. Use `refresh=false` for bulk imports
2. Use filters (cached) when possible
3. Limit result size
4. Use pagination
5. Denormalize data (avoid nested queries)

---

## ✅ Summary

**Elasticsearch in Job Search Service:**

1. **Stores** job documents in `jobs` index
2. **Indexes** text fields for fast full-text search
3. **Searches** across title, description, company
4. **Filters** by location and skills
5. **Ranks** results by relevance
6. **Returns** in milliseconds

**Key Operations:**

- Index: Add job document
- Search: Query with filters
- Get: Retrieve by ID
- Delete: Remove document

**Result:** Fast, scalable, intelligent job search! 🚀

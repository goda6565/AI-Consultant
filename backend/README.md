# Backend

## Overview

This is the backend for AI Consultant.

### Services

There are 3 services in the backend.

- Admin: Document Management (rest)
- Agent: Hearing Agent (rest)
- Vector: Sync Document to Vector Database (cloud tasks)
- Proposal Job: Generate Proposal Agent (cloud run job)

All services are running on the cloud run.

### Database

There are 2 databases in the backend.
- App: Document Management (postgres)
- Vector: Vector Database (postgres with pgvector extension)

## Run Application

```bash
go run main.go admin run
go run main.go agent run
go run main.go vector run
go run main.go proposal-job run
```
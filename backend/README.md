# Backend

## Overview

This is the backend for AI Consultant.

### Services

There are 3 services in the backend.

- Admin: Document Management (rest)
- Agent: AI Consultant Agent (websocket)
- Vector: Sync Document to Vector Database (pubsub)

## Run Application

```bash
go run main.go admin run
go run main.go agent run
go run main.go vector run
```
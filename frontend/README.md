# Frontend

## Overview

This is the frontend for AI Consultant.

### Stack
- Next.js 15 (App Router, Turbopack)
- TypeScript
- Tailwind CSS (+ typography)
- Firebase Auth
- Storybook
- Orval (API types)

### Architecture
- FSD (Feature-Sliced Design) を採用
- レイヤ: `shared → entities → features → widgets → pages → app`

### Run Application

```bash
pnpm install
pnpm dev
# open http://localhost:3000
```

### Scripts
- `pnpm dev`: Start dev server (Turbopack)
- `pnpm build`: Build production
- `pnpm start`: Start production server
- `pnpm storybook`: Start Storybook
- `pnpm lint`: Biome + Steiger
- `pnpm format`: Format code
- `pnpm orval`: Generate API types/clients

### Env Vars
- `NEXT_PUBLIC_ADMIN_API_URL`
- `NEXT_PUBLIC_AGENT_API_URL`
- Firebase project configs

### Docker
See `frontend/Dockerfile` for multi-stage build and runtime.

### Notes
- API calls attach Firebase ID token via axios interceptors
- SSE + SWR + local state for chat screen

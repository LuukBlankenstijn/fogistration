# Fogistration

Fogistration is a tool for synchronising laptops deployed with [FOG](https://fogproject.org/) to [DomJudge](https://www.domjudge.org/) for use in **ICPC contests**.
Its primary purpose is providing **live updates to the lockscreen wallpaper** of contestant laptops.

## Architecture

Fogistration is written fully in **Go** (except for the frontend). It is **modular** and built around clean communication through **Postgres**.
The system is designed so modules can run **independently**.
The different modules can even run replicated on for example k8s. This is not really a useful feature but a challenge I liked. It actually heavily contributed to the modular design.

### Modules

- **gRPC Server**  
  Handles streaming communication with deployed clients.

- **HTTP Server**  
  Provides management endpoints and frontend integration.

- **DomJudge Syncer**  
  Connects to the DomJudge API and keeps contest data (contests, teams, IPs, etc.) in sync.

- **Client**  
  A lightweight agent deployed on each laptop, connected via gRPC, that updates the lockscreen wallpaper and applies configuration changes in real time.

## Database Integration

Fogistration makes extensive use of **Postgres triggers** for efficient real-time updates.  
This allows live propagation of changes (e.g., IP updates, wallpaper changes) without polling.

## Key Features

- Real-time lockscreen wallpaper updates on contestant laptops.
- Lightweight client deployed to each laptop.
- Modular components (run standalone or together).
- Postgres-driven communication for reliability and simplicity.
- gRPC streaming for efficient client communication.

## Development

Requirements:

- Docker / Docker Compose (for local testing)
- yarn
- go 1.24+

```bash
git clone https://github.com/LuukBlankenstijn/fogistration
cd fogistration
```

Backends (HTTP, gRPC, domjudge syncer)

```bash
cd go
make swagger
docker compose up -d && docker compose logs -f
```

Frontend

```bash
cd frontend
yarn
yarn set version stable
corepack enable
yarn gen-client
yarn dev
```

Client

```bash
cd go
go run cmd/client/main.go -server=localhost:9090
```

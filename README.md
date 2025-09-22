# Knowledge Base IAM & Workflow Skeleton

This repository provides Go-based service skeletons for the account (IAM) and approval workflow capabilities required by the knowledge base project. Both services expose lightweight HTTP APIs with in-memory storage so the interfaces can be validated quickly before wiring real databases.

## Layout

- `cmd/iam`: entrypoint for the IAM service (users, roles, spaces, policies)
- `cmd/workflow`: entrypoint for the workflow service (flow definitions, instances)
- `internal/common`: shared utilities (config loading, logging, middleware, HTTP server helpers, generic in-memory storage)
- `internal/iam`: IAM domain models, repositories, services, handlers
- `internal/workflow`: workflow domain models, repositories, services, handlers
- `deploy`: Dockerfiles and compose file for offline packaging

## Getting Started

```bash
# build both binaries
make build

# run IAM service locally (requires postgres if IAM_DATABASE_DSN is set)
IAM_SERVER_PORT=8080 IAM_DATABASE_DSN="postgres://kb:kb@127.0.0.1:5432/kb_platform?sslmode=disable" go run ./cmd/iam

# run Workflow service locally
WORKFLOW_SERVER_PORT=8090 WORKFLOW_DATABASE_DSN="postgres://kb:kb@127.0.0.1:5432/kb_platform?sslmode=disable" go run ./cmd/workflow
```

Both services now require an `Authorization: Bearer <JWT>` header for protected endpoints; tokens are validated against `IAM_AUTH_JWT_SIGNING_KEY` (or the workflow equivalent). When `*_DATABASE_DSN` is omitted the services fall back to in-memory repositories (volatile). IAM user records use `phone` + hashed `password` as credentials, and `/api/auth/login` exchanges valid credentials for a JWT token that the Workflow service also trusts.

To prepare PostgreSQL tables run the bootstrap migration once:

```bash
psql "postgres://kb:kb@127.0.0.1:5432/kb_platform?sslmode=disable" -f deploy/migrations/001_init.sql
```

## Docker Packaging

1. Build images locally (ensures Go modules and binaries are baked into the container):
   ```bash
   make docker-iam
   make docker-workflow
   ```
2. Export for offline deployment:
   ```bash
   docker save iam-service:dev | gzip > iam-service.tar.gz
   docker save workflow-service:dev | gzip > workflow-service.tar.gz
   ```
3. Copy the tarballs together with `deploy/docker-compose.yml` to the target machine and load:
   ```bash
   gunzip -c iam-service.tar.gz | docker load
   gunzip -c workflow-service.tar.gz | docker load
   docker compose -f deploy/docker-compose.yml up -d
   ```

The compose file now provisions PostgreSQL (ports `5432`, credentials `kb/kb`) and wires DSNs for both services. The bootstrap SQL in `deploy/migrations/001_init.sql` is automatically mounted into the container and executed on the first start; repeat runs will be no-ops. For environments without Docker you can run the same SQL manually prior to launching the services.

## API Documentation

Swagger definitions are generated via `swag`. Regenerate after handler changes with:

```bash
(cd cmd/iam && swag init --dir . --generalInfo main.go --output ../docs/swagger/iam --instanceName iam --parseDependency --parseInternal)
(cd cmd/workflow && swag init --dir . --generalInfo main.go --output ../docs/swagger/workflow --instanceName workflow --parseDependency --parseInternal)
```

Generated artifacts:

- IAM service spec: `docs/swagger/iam/iam_swagger.yaml`
- Workflow service spec: `docs/swagger/workflow/workflow_swagger.yaml`

To browse the docs locally via Swagger UI:

```bash
docker run --rm -p 8088:8080 -e SWAGGER_JSON=/spec/iam.yaml \
  -v $(pwd)/docs/swagger/iam/iam_swagger.yaml:/spec/iam.yaml swaggerapi/swagger-ui
```

Import the YAML or JSON into Insomnia/Postman for quick client generation.

The compose file exposes IAM on port `8080` and Workflow on `8090` for intra-network access.

## Next Steps

- Replace the in-memory repositories with Postgres/Redis implementations.
- Swap the placeholder auth middleware with actual JWT validation.
- Extend workflow DSL (parallel branches, conditional routing) and persist audit logs.
- Harden error handling and add unit/integration tests per module.

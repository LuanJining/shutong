# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /workspace
COPY go.mod go.sum ./
COPY internal ./internal
COPY cmd/workflow ./cmd/workflow
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /workspace/out/workflow ./cmd/workflow

FROM alpine:3.20
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app
COPY --from=build /workspace/out/workflow /app/workflow
USER app
EXPOSE 8090
ENTRYPOINT ["/app/workflow"]

#!/bin/bash

echo "🚀 Starting Knowledge Base Platform (Spring Boot)"
echo "================================================"

# Check if PostgreSQL is running
echo "📦 Checking PostgreSQL..."
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "❌ PostgreSQL is not running. Please start PostgreSQL first."
    exit 1
fi
echo "✅ PostgreSQL is running"

# Check if MinIO is accessible
echo "📦 Checking MinIO..."
if ! curl -s http://localhost:9000/minio/health/live > /dev/null 2>&1; then
    echo "⚠️  MinIO is not accessible. Document upload may not work."
else
    echo "✅ MinIO is accessible"
fi

# Check if Qdrant is accessible
echo "📦 Checking Qdrant..."
if ! curl -s http://localhost:6333 > /dev/null 2>&1; then
    echo "⚠️  Qdrant is not accessible. Vector search may not work."
else
    echo "✅ Qdrant is accessible"
fi

echo ""
echo "🔨 Building and starting the application..."
./mvnw clean spring-boot:run


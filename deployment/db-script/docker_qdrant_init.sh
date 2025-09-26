docker pull qdrant/qdrant:latest

docker run -d --name qdrant -p 6333:6333 qdrant/qdrant

docker exec -it qdrant qdrant-client --host localhost --port 6333 --password qdrant

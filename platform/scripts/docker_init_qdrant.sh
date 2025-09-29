docker pull qdrant/qdrant:latest
docker run -d --name qdrant -p 6333:6333 qdrant/qdrant
docker exec -it qdrant qdrant-cli
docker exec -it qdrant qdrant-cli --host 0.0.0.0 --port 6333
docker exec -it qdrant qdrant-cli --host 0.0.0.0 --port 6333 --password qdrant
docker exec -it qdrant qdrant-cli --host 0.0.0.0 --port 6333 --password qdrant --create-collection
docker exec -it qdrant qdrant-cli --host 0.0.0.0 --port 6333 --password qdrant --create-collection --name kb_platform
docker exec -it qdrant qdrant-cli --host 0.0.0.0 --port 6333 --password qdrant --create-collection --name kb_platform --dimension 1536
docker exec -it qdrant qdrant-cli --host 0.0.0.0 --port 6333 --password qdrant --create-collection --name kb_platform --dimension 1536 --distance l2
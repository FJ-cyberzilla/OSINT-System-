# With environment variables
export ORCHESTRA_ADDR="tcp://localhost:5558"
export PORT="8080"
go run main.go

# Or with Docker
docker build -t osint-api .
docker run -p 8080:8080 -e ORCHESTRA_ADDR="tcp://host.docker.internal:5558" osint-api

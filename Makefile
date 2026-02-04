.PHONY: run build test clean docker-build docker-run docker-compose-up docker-compose-down

# Run the application locally
run:
	go run main.go

# Build the application
build:
	go build -o bin/weather-by-cep main.go

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t weather-by-cep .

# Run Docker container
docker-run:
	docker run -p 8080:8080 -e WEATHER_API_KEY=${WEATHER_API_KEY} weather-by-cep

# Start services with docker-compose
docker-compose-up:
	docker-compose up --build

# Stop services with docker-compose
docker-compose-down:
	docker-compose down

# Deploy to Google Cloud Run
deploy:
	gcloud run deploy weather-by-cep \
		--source . \
		--region us-central1 \
		--platform managed \
		--allow-unauthenticated \
		--set-env-vars "WEATHER_API_KEY=${WEATHER_API_KEY}"

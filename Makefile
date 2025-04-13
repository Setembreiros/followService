DEV-ENVIRONMENT=development
PROD-ENVIRONMENT=production

update:
	go mod tidy

build: update
	go build -o ./deployment/${PROD-ENVIRONMENT}/followService cmd/main.go cmd/startup.go

run:
	export ENVIRONMENT="${PROD-ENVIRONMENT}" && go run cmd/main.go cmd/startup.go

run-dev:
	export ENVIRONMENT="${DEV-ENVIRONMENT}" && go run ./cmd/main.go ./cmd/startup.go

run-dev-windows: 
	set ENVIRONMENT=${DEV-ENVIRONMENT} && go run ./cmd/main.go ./cmd/startup.go

test:
	go generate -v ./internal/... && go test ./internal/...
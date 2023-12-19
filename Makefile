ENV := $(PWD)/.env

# Environment variables for project
include $(ENV)

# Export all variable to sub-make
export

tests:
	go test ./...

build:
	go build ./cmd/app/main.go

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: gen
gen:
	mockgen -source=internal/services/message.go -destination internal/services/mocks/message_mock.go
	mockgen -source=internal/services/template.go -destination internal/services/mocks/template_mock.go
	mockgen -source=internal/services/user.go -destination internal/services/mocks/user_mock.go

migrate-create:
	 migrate create -ext sql -dir internal/databases/migrations $(name)

migrate-up:
	migrate -path ./internal/databases/migrations -database $(DATABASE_URL)?sslmode=disable up

migrate-fix:
	migrate -path ./internal/databases/migrations -database -database $(DATABASE_URL)?sslmode=disable force $(number)
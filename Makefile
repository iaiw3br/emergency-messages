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
	mockgen -source=internal/service/message.go -destination internal/store/postgres/mock/message_mock.go
	mockgen -source=internal/service/template.go -destination internal/store/postgres/mock/template_mock.go
	mockgen -source=internal/service/user.go -destination internal/store/postgres/mock/user_mock.go
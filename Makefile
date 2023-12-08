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
	mockgen -source=internal/store/message.go -destination internal/store/mock/message_mock.go
	mockgen -source=internal/store/template.go -destination internal/store/mock/template_mock.go
	mockgen -source=internal/store/user.go -destination internal/store/mock/user_mock.go
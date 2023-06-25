# .PHONY: lint test

all:
		make lint
		make test

lint:
		golangci-lint run --config .golangci.yml

test:
		go test -timeout 30s -cover ./...
mocks:
		mockery -all -dir dependency -case underscore
swagger:
	go mod vendor
	GO111MODULE=off swagger generate spec -m -o ./swagger.yml
	GO111MODULE=off swagger validate ./swagger.yml
	rm -rf vendor

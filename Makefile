
install:
	go mod download

build:
	go build

unit:
	@(go list ./...  | grep -v "vendor/" | xargs -n1 go test -race -v -cover)

docker:
	docker build -t aldor007/transformer-go -f Dockerfile . -t aldor007/transformer-go:latest

integrations:
	npm install
	./test-int/run-test.sh

format:
	@(go fmt ./...)
	@(go vet ./...)

tests: unit integrations

run-server:
	go run main.go


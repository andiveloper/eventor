build:
	go build -o eventor cmd/main.go

test:
	go test ./pkg ./cmd

run:
	go run cmd/main.go

run-sample:
	go run cmd/main.go -f sample_config.yaml

build:
	@go build -o bin/hailstorm cmd/hailstorm/main.go 

run: build
	@./bin/hailstorm

test:
	@go test ./...

clean:
	@rm -rf bin/hailstorm

goex:
	@./examples/go/build.sh
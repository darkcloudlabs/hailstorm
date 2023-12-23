wasm:
	tinygo build -o testdata/app.wasm -scheduler=none --no-debug -target wasi testdata/main.go

build:
	@go build -o bin/hailstorm cmd/hailstorm/main.go 

run: build
	@./bin/hailstorm

test:
	@go test ./...

clean:
	@rm -rf bin/hailstorm
	//@GOOS=wasip1 GOARCH=wasm go build -o app.wasm testdata/main.go
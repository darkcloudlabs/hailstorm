echo "building go application with tinygo..."
tinygo build -o examples/go/app.wasm -scheduler=none --no-debug -target wasi examples/go/main.go
echo "done!"
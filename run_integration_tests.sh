tinygo build -size short -no-debug -panic=trap -scheduler=none -gc=leaking -o main.wasm -target wasm-unknown ./ && ls -lh main.wasm && cd integration_tests && cargo run && cd ..
# ergo-lib-go
Go wrapper around C bindings for ErgoLib from [sigma-rust](https://github.com/ergoplatform/sigma-rust)

### Status
Package is still work in progress and will likely change in the future

### Library
The libraries under `lib` where compiled from `sigma-rust` with the following commands:
```
cross build -p ergo-lib-c --release --target x86_64-apple-darwin
cross build -p ergo-lib-c --release --target x86_64-unknown-linux-gnu
cross build -p ergo-lib-c --release --target aarch64-unknown-linux-gnu
cross build -p ergo-lib-c --release --target x86_64-pc-windows-gnu
```

### Credits
* [go-ergo](https://github.com/ross-weir/go-ergo) from [ross-weir](https://github.com/ross-weir) for initial code and examples
* [wasmer-go](https://github.com/wasmerio/wasmer-go) for package structure

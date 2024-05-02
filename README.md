# ergo-lib-go
Go wrapper around C bindings for ErgoLib from [sigma-rust](https://github.com/ergoplatform/sigma-rust)

### Install
```
go get -u github.com/sigmaspace-io/ergo-lib-go
```

### Supported Platforms
This library makes heavy use of cgo. A set of precompiled shared library objects are provided. For the time being the following platforms are supported and tested against: 

<table>
  <thead>
    <tr>
      <th>Platform</th>
      <th>Architecture</th>
      <th>Triple</th>
      <th>Supported</th>
      <th>Tested</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td rowspan="2">Linux</td>
      <td><code>amd64</code></td>
      <td><code>x86_64-unknown-linux-gnu</code></td>
      <td>✅</td>
      <td>✅</td>
    </tr>
    <tr>
      <td><code>aarch64</code></td>
      <td><code>aarch64-unknown-linux-gnu</code></td>
      <td>✅</td>
      <td>⏳</td>
    </tr>
    <tr>
      <td rowspan="2">Darwin</td>
      <td><code>amd64</code></td>
      <td><code>x86_64-apple-darwin</code></td>
      <td>✅</td>
      <td>✅</td>
    </tr>
    <tr>
      <td><code>aarch64</code></td>
      <td><code>aarch64-apple-darwin</code></td>
      <td>⏳</td>
      <td>⏳</td>
    </tr>
    <tr>
      <td>Windows</td>
      <td><code>amd64</code></td>
      <td><code>x86_64-pc-windows-gnu</code></td>
      <td>⏳</td>
      <td>⏳</td>
    </tr>
  </tbody>
</table>

### Supported sigma-rust versions
<table>
  <thead>
    <tr>
      <th>sigma-rust Version</th>
      <th>ergo-lib-go Version</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>v0.26.0</td>
      <td>v0.26.0</td>
    </tr>
  </tbody>
</table>

### Library
The libraries under `lib` were compiled from `sigma-rust` with the following commands:
```
cross build -p ergo-lib-c --release --target x86_64-apple-darwin
cross build -p ergo-lib-c --release --target x86_64-unknown-linux-gnu
cross build -p ergo-lib-c --release --target aarch64-unknown-linux-gnu
cross build -p ergo-lib-c --release --target x86_64-pc-windows-gnu
```

### Credits
* [go-ergo](https://github.com/ross-weir/go-ergo) from [ross-weir](https://github.com/ross-weir) for initial code and examples
* [wasmer-go](https://github.com/wasmerio/wasmer-go) for package structure

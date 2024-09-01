package ergo

// #cgo CFLAGS: -I${SRCDIR}/packaged/include
// #cgo LDFLAGS: -lergo
//
// #cgo linux,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/packaged/lib/linux-amd64 -L${SRCDIR}/packaged/lib/linux-amd64 -lm
// #cgo linux,arm64 LDFLAGS: -Wl,-rpath,${SRCDIR}/packaged/lib/linux-aarch64 -L${SRCDIR}/packaged/lib/linux-aarch64 -lm
// #cgo windows,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/packaged/lib/windows-amd64 -L${SRCDIR}/packaged/lib/windows-amd64 -lbcrypt -lntdll
// #cgo darwin,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/packaged/lib/darwin-amd64 -L${SRCDIR}/packaged/lib/darwin-amd64
// #cgo darwin,arm64 LDFLAGS: -Wl,-rpath,${SRCDIR}/packaged/lib/darwin-aarch64 -L${SRCDIR}/packaged/lib/darwin-aarch64
//
// #include <ergo.h>
import "C"

// See https://github.com/golang/go/issues/26366.
import (
	_ "github.com/sigmaspace-io/ergo-lib-go/packaged/include"
	_ "github.com/sigmaspace-io/ergo-lib-go/packaged/lib"
)

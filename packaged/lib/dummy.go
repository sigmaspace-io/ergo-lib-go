// See https://github.com/golang/go/issues/26366.
package lib

import (
	_ "github.com/sigmaspace-io/ergo-lib-go/packaged/lib/darwin-aarch64"
	_ "github.com/sigmaspace-io/ergo-lib-go/packaged/lib/darwin-amd64"
	_ "github.com/sigmaspace-io/ergo-lib-go/packaged/lib/linux-aarch64"
	_ "github.com/sigmaspace-io/ergo-lib-go/packaged/lib/linux-amd64"
	_ "github.com/sigmaspace-io/ergo-lib-go/packaged/lib/windows-amd64"
)

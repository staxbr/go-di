package di_test

import (
	"fmt"
	"github.com/staxbr/go-di"
	"io/ioutil"
	"testing"
)

func BenchmarkResolveContainer(b *testing.B) {
	ctx := di.NewContext()
	var c di.Container
	for n := 0; n < b.N; n++ {
		c = di.ResolveContainer(ctx) // holding the result to prevent compiler optimizations
	}
	_, _ = fmt.Fprint(ioutil.Discard, c)
}

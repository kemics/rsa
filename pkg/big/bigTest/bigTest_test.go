package bigTest_test

import (
	"fmt"
	"github.com/kemics/rsa/pkg/big/bigTest"
	"testing"
)

// max: 21267647932558653966460912964485513216

func TestName(t *testing.T) {
	i := bigTest.NewFrom10Text("21267647932558653966460912964485513215")
	fmt.Println(bigTest.StdBig(i).Text(10))
}

package rsa

import (
	"fmt"
	"github.com/kemics/rsa/pkg/bitBig"
	"testing"
)

func TestName(t *testing.T) {
	DO()
}

func TestEuklidGCD(t *testing.T) {
	d, x, y := euklidGCD(bitBig.FromUint(100), bitBig.FromUint(77))
	fmt.Println(d.Uint64(), x.Uint64(), y.Uint64())
	fmt.Println(x.Uint64(), y.Uint64())
}

func TestFindPublicKeyPart(t *testing.T) {
	d, m := bitBig.FromUint(100), bitBig.FromUint(77)
	fmt.Println(findPublicKeyPart(d, m))
}

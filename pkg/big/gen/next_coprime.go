package gen

import (
	"github.com/kemics/rsa/pkg/big"
)

func NextCoprime(a *big.Int) *big.Int {
	g := NewGenerator()
	for {
		b := g.NextBigInt()
		if a.GCD(b).IsUint(1) {
			return b
		}
	}
}

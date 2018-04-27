package gen

import (
	"github.com/kemics/rsa/pkg/bitBig"
)

func NextCoprime(a *bitBig.Int) *bitBig.Int {
	g := NewGenerator()
	for {
		b := g.NextBigInt()
		if a.GCD(b).Uint64() == 1 {
			return b
		}
	}
}

package gen

import (
	"fmt"
	"github.com/kemics/rsa/pkg/big"
	"testing"
)

func TestIsProbablyPrime(t *testing.T) {
	tr, f := 0, 0
	for i := 0; i < 10000; i++ {
		if next := NewGenerator().NextBigInt(); IsProbablyPrime(next) {
			fmt.Println(next)
			tr++
		} else {
			f++
		}
		if i%10 == 0 {
			fmt.Println(tr, f)
		}
	}

}

func TestIsProbablyPrimeCarmichael(t *testing.T) {
	for _, k := range []uint32{561, 41041, 825265, 321197185, 62745, 63973, 75361, 101101, 126217, 172081, 188461, 278545, 340561} {
		if IsProbablyPrime(big.FromUint(k)) {
			t.Fatal(k)
		}
	}
}

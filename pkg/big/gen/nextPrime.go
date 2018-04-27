package gen

import (
	"context"
	"fmt"
	"github.com/kemics/rsa/pkg/big"
)

const primeGoroutines = 10

func NextPrime() *big.Int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan *big.Int)
	for i := 0; i < primeGoroutines; i++ {
		go func(ctx context.Context) {
			for {
				g := NewGenerator()
				select {
				case <-ctx.Done():
					return
				default:
					p := g.NextBigInt()
					if IsProbablyPrime(p) {
						resultChan <- p
					}
				}
			}
		}(ctx)
	}
	p := <-resultChan
	return p
}

func IsProbablyPrime(i *big.Int) bool {
	for _, base := range []int{2, 3, 4} {
		fmt.Println(base)
		if !millerRabin(i, big.FromUint(uint32(base))) {
			return false
		}
	}
	return true
}

func millerRabin(i *big.Int, a *big.Int) bool {
	fmt.Println("1")
	if i.IsEven() {
		return false
	}
	fmt.Println("2")
	n := i.Copy()
	nMinus1 := n.Copy()
	nMinus1.Sub(big.One())
	fmt.Println("3")
	two := big.FromUint(2)
	k := uint64(0)
	var m *big.Int
	div := i.Copy()
	fmt.Println("4")
	for {
		fmt.Println("BEFORE", div, two)
		div, _ = div.DivMod(two)
		fmt.Println(div)
		if !div.IsEven() { // % 2 = 1
			m = div
			break
		}
		k++
	}
	fmt.Println("5")

	t := a.PowMod(m, n)
	fmt.Println("6")
	if t.IsUint(1) || t.Cmp(nMinus1) == 0 {
		return true
	}
	fmt.Println("7")
	defer fmt.Println("8")
	for c := uint64(0); c < k; c++ {
		t = t.PowMod(two, n)
		if t.IsUint(1) {
			return false
		}
		if t.Cmp(nMinus1) == 0 {
			return true
		}
	}
	return false
}

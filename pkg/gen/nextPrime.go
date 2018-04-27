package gen

import (
	"context"
	"github.com/kemics/rsa/pkg/bitBig"
)

const primeGoroutines = 10

func NextPrime() *bitBig.Int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan *bitBig.Int)
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

func IsProbablyPrime(i *bitBig.Int) bool {
	for _, base := range []int{2, 3, 4} { // можно было использовать случайный генератор, для скорости работы так
		if !isProbablyPrime(i, bitBig.FromUint(uint64(base))) {
			return false
		}
	}
	// если m пройдено тестов, то вероятноятность (1/4)**m
	return true
}

func IsProbablyPrimeWithRand(i *bitBig.Int, iter int) bool {
	g := NewGenerator()
	for c := 0; c < iter; c++ {
		if !isProbablyPrime(i, bitBig.FromUint(uint64(g.NextUint()))) { // меньше проверяемого
			return false
		}
	}
	return true
}

func isProbablyPrime(i *bitBig.Int, a *bitBig.Int) bool {
	if i.IsEven() {
		return false
	}
	n := i.Copy()
	nMinus1 := n.Copy()
	nMinus1.Bits[0] = bitBig.ZeroBit // i - 1

	two := bitBig.FromUint(2)
	k := bitBig.New()
	var m *bitBig.Int
	div := i.Copy()
	for {
		div, _ = div.DivMod(two)
		if !div.IsEven() { // % 2 = 1
			m = div
			break
		}
		k.Add(bitBig.FromUint(1))
	}

	t := a.PowMod(m, n)
	if t.IsOne() || t.Cmp(nMinus1) == 0 {
		return true
	}

	// проверяем до k, которое получено из i = представления 2**k*
	for c := 0; c < int(k.Uint64()); c++ {
		t = t.PowMod(two, n)
		if t.IsOne() {
			return false
		}
		if t.Cmp(nMinus1) == 0 {
			return true
		}
	}
	return false
}

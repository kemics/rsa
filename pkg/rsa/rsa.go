package rsa

import (
	"github.com/kemics/rsa/pkg/bitBig"
	"github.com/kemics/rsa/pkg/gen"
)

type RSA struct {
	e, n *bitBig.Int // public keys
	d    *bitBig.Int // private key
}

func New() *RSA {
	p, q := gen.NextPrime(), gen.NextPrime()
	one := bitBig.FromUint(1)

	n := p.Copy().Mult(q)
	m := p.Copy().Sub(one).Mult(q.Copy().Sub(one))
	d := gen.NextCoprime(m)

	e := findPublicKeyPart(d, m)

	return &RSA{
		e: e.Copy(),
		n: n.Copy(),
		d: d.Copy(),
	}
}

func (r *RSA) Encode(message string) []*bitBig.Int {
	result := make([]*bitBig.Int, len(message))
	for i, b := range message {
		bigChar := bitBig.FromUint(uint64(b))
		result[i] = bigChar.PowMod(r.e.Copy(), r.n)
	}
	return result
}

func (r *RSA) Decode(message ...*bitBig.Int) string {
	result := make([]rune, len(message))
	for i, b := range message {
		result[i] = rune(b.PowMod(r.d.Copy(), r.n).Uint64())
	}
	return string(result)
}

func findPublicKeyPart(d, m *bitBig.Int) *bitBig.Int {
	// e * d mod m = 1, e - ?
	g, x, _ := euklidGCD(d, m)
	if !g.IsOne() {
		panic("something wrong")
	}
	e := x.Copy().Add(m).Mod(m) // for the case if its negative
	return e
}

//// a*x + b*y
func euklidGCD(a, b *bitBig.Int) (d, x, y *bitBig.Int) {
	if b.IsZero() {
		return a, bitBig.FromUint(1), bitBig.FromUint(0)
	}
	d, x1, y1 := euklidGCD(b, a.Copy().Mod(b))
	return d, y1, x1.Copy().Sub((a.Copy().Div(b)).Mult(y1))
}

package gen

import (
	"github.com/kemics/rsa/pkg/bitBig"
	"time"
)

const BigIntLength = 4 //16//1//16//32 // 32*32+1 = 1025

type Generator struct {
	State uint32
}

func NewGenerator() *Generator {
	now := time.Now().UnixNano()
	return &Generator{
		State: uint32((now >> 32) ^ now),
	}
}

func (g *Generator) NextUint() uint32 {
	x := g.State
	// https://en.wikipedia.org/wiki/Xorshift
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	g.State = x
	return x
}

func (g *Generator) NextBigInt() *bitBig.Int {
	res := &bitBig.Int{}
	for i := 0; i < BigIntLength; i++ {
		res.AppendUint(g.NextUint())
	}
	res.AppendBit(bitBig.OneBit)

	return res
}

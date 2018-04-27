package gen

import (
	"log"
	"testing"
)

func TestGenerator_NextBigInt(t *testing.T) {
	g := NewGenerator()
	for i := 0; i < 100; i++ {
		b := g.NextBigInt()
		if b.Len != 32*BigIntLength+1 {
			t.Fatalf("Got lenth %d", b.Len)
		}
	}
}

func TestGenerator_NextUint(t *testing.T) {
	m := map[uint32]struct{}{}
	g := NewGenerator()
	for i := 1; i < 100000000; i++ {
		if i%10000000 == 0 {
			log.Printf("Period > %d\n", i)
		}
		next := g.NextUint()
		if _, exist := m[next]; exist {
			t.Fatalf("Period is '%d'", i)
		}
		m[next] = struct{}{}
	}
}

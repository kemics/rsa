package bigTest

import (
	"github.com/kemics/rsa/pkg/big"
	"math"
	stdbig "math/big"
)

func NewFrom10Text(s string) *big.Int {
	i := new(stdbig.Int)
	i.SetString(s, 10)

	maxUint32_1 := new(stdbig.Int).SetUint64(uint64(math.MaxUint32))
	maxUint32_2 := new(stdbig.Int).Mul(maxUint32_1, maxUint32_1)
	maxUint32_3 := new(stdbig.Int).Mul(maxUint32_2, maxUint32_1)

	d3Int, i := new(stdbig.Int).DivMod(i, maxUint32_3, new(stdbig.Int))
	d2Int, i := new(stdbig.Int).DivMod(i, maxUint32_2, new(stdbig.Int))
	d1Int, d0Int := new(stdbig.Int).DivMod(i, maxUint32_1, new(stdbig.Int))

	bigInt := big.New()
	bigInt.AppendDigits(uint32(d0Int.Uint64()), uint32(d1Int.Uint64()), uint32(d2Int.Uint64()), uint32(d3Int.Uint64()))
	return bigInt
}

func StdBig(b *big.Int) *stdbig.Int {
	i := new(stdbig.Int)

	maxUint32_1 := new(stdbig.Int).SetUint64(uint64(math.MaxUint32))
	maxUint32_2 := new(stdbig.Int).Mul(maxUint32_1, maxUint32_1)
	maxUint32_3 := new(stdbig.Int).Mul(maxUint32_2, maxUint32_1)

	if len(b.Digits) >= 4 {
		i.Add(i, maxUint32_3.Mul(maxUint32_3, new(stdbig.Int).SetUint64(uint64(b.Digits[3]))))
	}

	if len(b.Digits) >= 3 {
		i.Add(i, maxUint32_2.Mul(maxUint32_2, new(stdbig.Int).SetUint64(uint64(b.Digits[2]))))
	}

	if len(b.Digits) >= 2 {
		i.Add(i, maxUint32_1.Mul(maxUint32_1, new(stdbig.Int).SetUint64(uint64(b.Digits[1]))))
	}

	if len(b.Digits) >= 2 {
		i.Add(i, new(stdbig.Int).SetUint64(uint64(b.Digits[0])))
	}
	return i
}

func NewStd() *stdbig.Int {
	return new(stdbig.Int)
}

//
//import (
//	"math/big"
//	"fmt"
//)
//
//func ()  {
//
//}

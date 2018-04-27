package big

import (
	"math"
)

type sign int8

const (
	negative sign = -1
	positive sign = 1
)

const base = uint64(math.MaxUint32) // + 1

type Int struct {
	sign   sign
	Digits []uint32
}

func New() *Int {
	return &Int{
		sign:   positive,
		Digits: []uint32{},
	}
}

func (i *Int) AppendDigits(digits ...uint32) *Int {
	i.Digits = append(i.Digits, digits...)
	return i
}

func (i *Int) shrink() {
	for c := len(i.Digits) - 1; c >= 0; c-- {
		if i.Digits[c] == 0 {
			i.Digits = i.Digits[:len(i.Digits)-1] // delete last
		} else {
			break
		}
	}
}

func (i *Int) PushBack(nums ...uint32) *Int {
	res := i.Shift(len(nums))
	for c := 0; c < len(nums); c++ {
		res.Digits[c] = nums[c]
	}
	return res
}

func (i *Int) Shift(num int) *Int {
	return &Int{
		sign:   i.sign,
		Digits: append(make([]uint32, num), i.Digits...),
	}
}

func (i *Int) Cmp(other *Int) int {
	if i.sign != other.sign {
		if len(i.Digits) == 0 && len(other.Digits) == 0 {
			return 0
		}
		if i.sign == positive {
			return 1
		}
		return -1
	}
	if len(i.Digits) > len(other.Digits) {
		return 1
	} else if len(i.Digits) < len(other.Digits) {
		return -1
	}
	for c := len(i.Digits) - 1; c >= 0; c-- {
		if i.Digits[c] > other.Digits[c] {
			return 1
		} else if i.Digits[c] < other.Digits[c] {
			return -1
		}
	}
	return 0
}

func (i *Int) ChangeSign() *Int {
	if i.sign == negative {
		i.sign = positive
	} else {
		i.sign = negative
	}
	return i
}

func (i *Int) getOrZero(index int) uint32 {
	if index >= len(i.Digits) {
		return 0
	}
	return i.Digits[index]
}

func (i *Int) Copy() *Int {
	return &Int{
		sign:   i.sign,
		Digits: append([]uint32(nil), i.Digits...),
	}
}

func (i *Int) IsZero() bool {
	if len(i.Digits) == 0 {
		return true
	}
	return len(i.Digits) == 1 && i.Digits[0] == 0
}

func (i *Int) IsUint(num uint32) bool {
	return len(i.Digits) == 1 && i.Digits[0] == num
}

func (i *Int) IsEven() bool {
	if i.IsZero() {
		return true
	}
	return i.Digits[0]%2 == 0
}

func One() *Int {
	return FromUint(1)
}

func FromUint(num uint32) *Int {
	return New().AppendDigits(num)
}

func (i *Int) Mod(other *Int) *Int {
	_, mod := i.DivMod(other)
	return mod
}

func (i *Int) PowMod(n *Int, mod *Int) *Int {
	a := i.Mod(mod)
	if n.IsZero() {
		return One()
	}
	if !n.IsEven() {
		return a.PowMod(n.Sub(One()), mod).Mult(i).Mod(mod)
	}
	ndiv2, _ := n.DivMod(FromUint(2))
	b := a.PowMod(ndiv2, mod)
	return b.Mult(b).Mod(mod)
}

func (i *Int) GCD(other *Int) (gcd *Int) {
	if other.IsZero() {
		return i.Copy()
	}
	return other.GCD(i.Mod(other))
}

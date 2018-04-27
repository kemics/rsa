package bitBig

import (
	"fmt"
	"math/bits"
	"strings"
)

const MaxLen = 2048 // MaxLen must be big enough (f.e. > 100) to support some tests

// Bit is a single bit with two states.
type Bit bool

const (
	ZeroBit Bit = false
	OneBit  Bit = true
)

// getBit tries to retrieve bit from any supported value types.
func getBit(i interface{}) Bit {
	switch i.(type) {
	case int:
		if i.(int) == 0 {
			return ZeroBit
		}
		return OneBit
	case bool:
		return Bit(i.(bool))
	case Bit:
		return i.(Bit)
	default:
		panic(fmt.Sprintf("Unsupported type '%T'", i))
	}
}

// Int represents unsigned big integer with max length 2048 bits.
// Panics if overflow happens.
type Int struct {
	negative bool
	Len      uint16
	Bits     [MaxLen]Bit
}

// New creates big unsigned integer from bits.
// It removes leading zeros.
func New(bits ...interface{}) *Int {
	bigInt := &Int{}

	if len(bits) == 0 {
		bigInt.Len = 1
		return bigInt
	}

	onlyZeros := true
	// check if first bits are "zeros"
	for i, bit := range bits {
		b := getBit(bit)
		if b != ZeroBit { // first
			bits = bits[i:] // remove leading zero bits
			onlyZeros = false
			break
		}
	}
	if onlyZeros {
		bits = bits[:1]
	}
	if len(bits) > MaxLen {
		panic("too big value") // this not really idiomatic, but to simplify API it will not return an error
	}
	for i, b := range bits {
		bigInt.Bits[len(bits)-i-1] = getBit(b)
	}
	bigInt.Len = uint16(len(bits))
	return bigInt
}

func FromUint(u uint64) *Int {
	// TODO: rewrite, because it is written at 3am
	i := &Int{Len: 1}

	for c := 0; c < 64; c++ {
		before := bits.OnesCount64(u)
		u >>= 1
		after := bits.OnesCount64(u)

		if after-before != 0 {
			i.Bits[c] = OneBit
			i.Len = uint16(c) + 1
		}
	}
	return i
}

func (i *Int) AppendUint(u uint32) {
	for c := 0; c < 32; c++ {
		before := bits.OnesCount32(u)
		u >>= 1
		after := bits.OnesCount32(u)

		if after-before != 0 { // 1
			i.AppendBit(OneBit)
		} else {
			i.AppendBit(ZeroBit)
		}
	}
}

func (i *Int) AppendBit(b Bit) *Int {
	i.Bits[i.Len] = b
	i.Len++
	return i
}

func (i *Int) Shrink() {
	if i.Len <= 1 {
		i.Len = 1
		return
	}
	var c int
	for c = int(i.Len - 1); c > 0; c-- {
		if i.Bits[c] == OneBit {
			break
		}
	}
	i.Len = uint16(c + 1)
}

func (i *Int) Cmp(other *Int) int {
	if i.IsZero() && other.IsZero() {
		return 0
	}
	if i.negative != other.negative {
		if i.negative {
			return -1
		}
		return 1
	}
	if i.Len != other.Len {
		if i.Len < other.Len {
			if i.negative {
				return 1
			}
			return -1
		}
		if i.negative {
			return -1
		}
		return 1
	}
	for c := i.Len; c > 0; c-- {
		if i.Bits[c-1] != other.Bits[c-1] {
			if i.Bits[c-1] == ZeroBit {
				if i.negative {
					return 1
				}
				return -1
			}
			if i.negative {
				return -1
			}
			return 1
		}
	}
	return 0
}

func (i *Int) CmpAbs(other *Int) int {
	if i.IsZero() && other.IsZero() {
		return 0
	}
	if i.Len != other.Len {
		if i.Len < other.Len {
			return -1
		}
		return 1
	}
	for c := i.Len; c > 0; c-- {
		if i.Bits[c-1] != other.Bits[c-1] {
			if i.Bits[c-1] == ZeroBit {
				return -1
			}
			return 1
		}
	}
	return 0
}

// String returns debug string representation.
// It has following format: "{number of trailing zeros}1001...1010".
func (i Int) String() string {
	s := new(strings.Builder)
	if i.negative {
		s.WriteByte('-')
	}
	// it doesn't use integer length to ensure that it has right trailing zeros
	var trailngZeros int
	var stoppedZeros bool
	for c := MaxLen - 1; c >= 0; c-- { // c - cursor
		if !stoppedZeros {
			if i.Bits[c] == ZeroBit {
				trailngZeros++
			} else {
				stoppedZeros = true
				s.WriteString(fmt.Sprintf("{%d}1", trailngZeros))
			}
			continue
		}
		if i.Bits[c] == ZeroBit {
			s.WriteByte('0')
		} else {
			s.WriteByte('1')
		}
	}
	if !stoppedZeros {
		return fmt.Sprintf("{%d}", MaxLen)
	}
	return s.String()
}

// Add changes value of calling integer number to
// sum of that number and argument.
// Returns itself for better API.
func (i *Int) Add(another *Int) *Int {
	if i.negative != another.negative { // different sign
		return i.Sub(another.Copy().ChangeSign())
	}

	maxLen := i.Len
	if another.Len > maxLen {
		maxLen = another.Len
	}
	i.Len = maxLen
	var carry Bit
	for c := uint16(0); c < maxLen; c++ {
		b1, b2 := i.Bits[c], another.Bits[c]
		if b1 == ZeroBit && b2 == ZeroBit {
			i.Bits[c] = carry
			carry = ZeroBit
		} else if b1 != b2 { // b1 or b2 == 1, but not both
			if carry == OneBit { // then carry still one, but b1 -> 0
				i.Bits[c] = ZeroBit
			} else {
				i.Bits[c] = OneBit
			} // else nothing we need to do, 0+1=1 (+carry == 0)
		} else { // both 1
			if carry == ZeroBit { // 1+1=0, (carry==0 --> 1)
				i.Bits[c] = ZeroBit
				carry = OneBit
			} else { // 1+1=1, (carry==1 --> 1)
				i.Bits[c] = OneBit
			}
		}
	}
	if carry == OneBit {
		if i.Len == MaxLen { // overflow
			panic("overflow while sum operation")
		}
		i.AppendBit(OneBit)
	}
	i.Shrink()
	return i
}

func (i *Int) ChangeSign() *Int {
	i.negative = !i.negative
	return i
}

// Copy creates a new different equal number.
func (i Int) Copy() *Int {
	res := new(Int)
	res.Len = i.Len
	for c := 0; c < MaxLen; c++ {
		res.Bits[c] = i.Bits[c]
	}
	res.negative = i.negative
	return res
}

// Shift adds additional zeros in the end.
// F.e. 111 shift 2 equals 11100.
func (i *Int) Shift(n uint16) *Int {
	if i.Len+n > MaxLen {
		panic("too long shift")
	}
	for c := int(i.Len - 1); c >= 0; c-- {
		b := i.Bits[c]
		i.Bits[c] = ZeroBit
		i.Bits[uint16(c)+n] = b
	}
	i.Len += n
	return i
}

func (i *Int) PushBack(bit Bit) *Int {
	if i.Len == 0 {
		i.Len = 1
		i.Bits[0] = bit
		return i
	}
	if i.Len+1 > MaxLen {
		panic("too long push")
	}
	for c := int(i.Len - 1); c >= 0; c-- {
		i.Bits[c+1] = i.Bits[c]
	}
	i.Bits[0] = bit
	i.Len++
	return i
}

// Sub changes value of calling Int to difference between two Ints.
// another must be less than i.
func (i *Int) Sub(another *Int) *Int {
	if i.negative != another.negative { // different sign
		return i.Add(another.Copy().ChangeSign())
	}

	if cmp := i.CmpAbs(another); cmp == 0 {
		*i = Int{Len: 1}
		return i
	} else if cmp == -1 { // i < another
		*i = *another.Copy().Sub(i)
		i.ChangeSign()
		return i
	}
	var carry Bit
	l := i.Len
	i.Len = 1 // because zero value has length 1
	for c := uint16(0); c < l; c++ {
		b1, b2 := i.Bits[c], another.Bits[c]
		if b1 == b2 {
			// 0 - 0 or 1 - 1 = carry, carry doesn't change
			i.Bits[c] = carry
		} else if b1 == OneBit && b2 == ZeroBit {
			// 1 - 0 = 1 - carry, carry --> 0
			if carry == OneBit {
				carry = ZeroBit
				i.Bits[c] = ZeroBit
			}
		} else if b1 == ZeroBit && b2 == OneBit {
			// 0 - 1 = 1 - carry, carry --> 1
			if carry == OneBit {
				i.Bits[c] = ZeroBit
			} else {
				i.Bits[c] = OneBit
				carry = OneBit
			}
		}
		if i.Bits[c] == OneBit { // update length
			i.Len = c + 1
		}
	}
	if carry == OneBit {
		panic("Int is subtracted by bigger int")
	}
	i.Shrink()
	return i
}

// Uint64 tries to convert Int to uint64.
// Panics if Int is too big.
func (i Int) Uint64() uint64 {
	if i.Len > 64 {
		panic("Trying to convert to big number to uint64")
	}
	var result uint64
	for c := int(i.Len - 1); c >= 0; c-- {
		if i.Bits[c] == OneBit {
			result += 1 << uint(c)
		}
	}
	return result
}

func (i *Int) Mult(another *Int) *Int {
	res := New()
	for c := 0; c < int(another.Len); c++ {
		if another.Bits[c] == OneBit {
			res.Add(i.Copy().Shift(uint16(c)))
		}
	}
	res.negative = i.negative != another.negative
	*i = *res
	res.Shrink()
	return res
}

func (i *Int) DivMod(another *Int) (div *Int, mod *Int) {
	other := another.Copy()
	other.negative = i.negative
	d := &Int{}
	div = &Int{}
	for c := i.Len; c > 0; c-- {
		d.PushBack(i.Bits[c-1])
		d.Shrink()
		if d.Cmp(other) >= 0 { // >=
			div.PushBack(OneBit)
			d.Sub(other)
		} else {
			div.PushBack(ZeroBit)
		}
	}
	d.Shrink()
	d.negative = false // check it later
	div.negative = i.negative != another.negative
	return div, d
}

func (i *Int) GCD(another *Int) (gcd *Int) {
	if another.Len == 1 && another.Bits[0] == ZeroBit {
		return i.Copy()
	}
	_, mod := i.DivMod(another)
	return another.GCD(mod)
}

func (i *Int) IsZero() bool {
	return i.Len <= 1 && i.Bits[0] == ZeroBit
}

func (i *Int) IsOne() bool {
	return i.Len == 1 && i.Bits[0] == OneBit
}

func (i *Int) IsEven() bool {
	return i.Bits[0] == ZeroBit
}

func (i *Int) Pow(n *Int) *Int {
	a := i.Copy()
	if n.IsZero() {
		return FromUint(1)
	}
	if !n.IsEven() {
		return a.Pow(n.Sub(FromUint(1))).Mult(a)
	} else {
		ndiv2, _ := n.DivMod(FromUint(2))
		b := a.Pow(ndiv2)
		return b.Mult(b.Copy())
	}
}

func (i *Int) Mod(n *Int) *Int {
	_, i = i.DivMod(n)
	return i
}

func (i *Int) Div(n *Int) *Int {
	i, _ = i.DivMod(n)
	return i
}

func (i *Int) PowMod(n *Int, mod *Int) *Int {
	a := i.Copy().Mod(mod)
	if n.IsZero() {
		return FromUint(1)
	}
	if !n.IsEven() {
		return a.PowMod(n.Sub(FromUint(1)), mod).Mult(a).Mod(mod)
	} else {
		ndiv2, _ := n.DivMod(FromUint(2))
		b := a.PowMod(ndiv2, mod)
		return b.Mult(b.Copy()).Mod(mod)
	}
}

package big

func (i *Int) Add(other *Int) *Int {
	if i.sign != other.sign {
		if i.sign == negative {
			return other.Sub(i)
		}
		return i.Sub(other)
	}
	maxLen := len(i.Digits)
	if maxLen < len(other.Digits) {
		maxLen = len(other.Digits)
	}

	result := &Int{
		sign:   i.sign,
		Digits: make([]uint32, maxLen+1),
	}
	defer result.shrink()

	var carry uint64 = 0
	for c := 0; c < maxLen; c++ {
		d1, d2 := i.getOrZero(c), other.getOrZero(c)
		carry += uint64(d1) + uint64(d2)
		result.Digits[c] = uint32(carry)
		carry = carry / base
	}
	if carry != 0 {
		result.Digits[len(result.Digits)-1] = uint32(carry)
	}
	return result
}

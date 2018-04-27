package big

func (i *Int) Sub(other *Int) *Int {
	if i.sign != other.sign {
		return i.Add(other.Copy().ChangeSign())
	}
	maxLen := len(i.Digits)
	if maxLen < len(other.Digits) {
		maxLen = len(other.Digits)
	}
	result := &Int{
		sign:   i.sign,
		Digits: make([]uint32, maxLen),
	}
	cmp := i.Cmp(other)
	var num1, num2 *Int // num1 >= num2
	if cmp == 0 {
		return New() // zero
	} else if cmp == 1 {
		num1, num2 = i, other
		result.sign = positive
	} else {
		num1, num2 = other, i
		result.sign = negative
	}
	defer result.shrink()

	var carry uint64 = 0
	for c := 0; c < maxLen; c++ {
		d1, d2 := num1.getOrZero(c), num2.getOrZero(c)
		var diff = int64(d1) - int64(d2) - int64(carry)
		if diff >= 0 {
			carry = 0
		} else {
			diff += int64(base)
			carry = 1
		}
		result.Digits[c] = uint32(diff)
	}
	if carry != 0 {
		panic("something wrong")
	}
	return result
}

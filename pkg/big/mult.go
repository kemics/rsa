package big

func (i *Int) Mult(other *Int) *Int {
	res := New()
	res.sign = i.sign * other.sign
	for c := 0; c < len(other.Digits); c++ {
		t := other.Digits[c]
		toAdd := i.Copy().Shift(c)
		for t > 0 {
			t--
			res = res.Add(toAdd)
		}
	}
	return res
}

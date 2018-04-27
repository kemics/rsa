package big

import (
	"fmt"
)

func (i *Int) DivMod(other *Int) (div *Int, mod *Int) {
	d := New()
	div = New()
	for c := len(i.Digits); c > 0; c-- {
		d = d.PushBack(i.Digits[c-1])
		fmt.Println(d)
		if d.Cmp(other) >= 0 { // >=
			var digit uint32
			for d.Cmp(other) >= 0 {
				digit++
				d = d.Sub(other)
			}
			div = div.PushBack(digit)
		} else {
			div = div.PushBack(0)
		}
	}
	div.sign = i.sign * other.sign
	d.sign = positive
	return div, d
}

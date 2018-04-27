package bitBig

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type number = []interface{} // just to simplify table driven test

	input := []struct {
		input number
		want  func() Int // it is really hard to create such a struct
	}{
		{input: number{0, 1, 1, 1}, want: func() Int {
			i := Int{Len: 3}
			i.Bits[0], i.Bits[1], i.Bits[2] = true, true, true
			return i
		}},

		{input: number{0, 0, 0, 0}, want: func() Int {
			i := Int{Len: 1}
			i.Bits[0] = false
			return i
		}},

		{input: number{true, true, true}, want: func() Int {
			i := Int{Len: 3}
			i.Bits[0], i.Bits[1], i.Bits[2] = true, true, true
			return i
		}},

		{input: number{0, 0, 1, 0, 1}, want: func() Int {
			i := Int{Len: 3}
			i.Bits[0], i.Bits[1], i.Bits[2] = true, false, true
			return i
		}},

		{input: number{}, want: func() Int {
			return Int{Len: 1}
		}},
	}

	for index, i := range input {
		got := *New(i.input...)

		if want := i.want(); !reflect.DeepEqual(want, got) {
			t.Fatalf("#%d: want integer '%s' (len: %d), got '%s' (len: %d)", index, want, want.Len, got, got.Len)
		}
	}
}

func TestInt_String(t *testing.T) {
	input := []struct {
		num  *Int
		want string
	}{
		{num: New(1, 1, 1, 0), want: fmt.Sprintf("{%d}1110", MaxLen-4)},
		{num: New(1, 1, 1, 1), want: fmt.Sprintf("{%d}1111", MaxLen-4)},
		{num: New(0, 0, 0, 0), want: fmt.Sprintf("{%d}", MaxLen)},
		{num: New(1, 0, 0, 0, 0), want: fmt.Sprintf("{%d}10000", MaxLen-5)},
		{num: New(OneBit, ZeroBit, ZeroBit, ZeroBit, ZeroBit), want: fmt.Sprintf("{%d}10000", MaxLen-5)},
	}
	for index, i := range input {
		if got := i.num.String(); !reflect.DeepEqual(got, i.want) {
			t.Fatalf("#%d: want '%s', got '%s'", index, i.want, got)
		}
	}
}

func TestInt_Add(t *testing.T) {
	input := []struct {
		num1 *Int
		num2 *Int
		sum  *Int
	}{
		{num1: New(0), num2: New(1, 1, 1, 0), sum: New(1, 1, 1, 0)},
		{num1: New(1, 1, 1, 0), num2: New(1), sum: New(1, 1, 1, 1)},
		{num1: New(1, 1, 1, 0), num2: New(0), sum: New(1, 1, 1, 0)},
		{num1: New(1, 1, 1, 1), num2: New(1), sum: New(1, 0, 0, 0, 0)},
		{num1: New(1, 1, 1, 1), num2: New(1, 1, 1, 1), sum: New(1, 1, 1, 1, 0)},
		{num1: New(0), num2: New(0), sum: New(0)},
		{num1: New(1, 1, 1, 0), num2: New(1, 1), sum: New(1, 0, 0, 0, 1)},
	}

	for index, i := range input {
		got := i.num1.Copy()
		got.Add(i.num2)
		if want := i.sum; !reflect.DeepEqual(*want, *got) {
			t.Fatalf("#%d: want integer '%s' (len: %d), got '%s' (len: %d)\nDetails: '%s'+'%s'='%s' (want '%s')",
				index, want, want.Len, got, got.Len, i.num1, i.num2, got, want)
		}

		// after cheking A+B=C, lets check B+A=C

		got2 := i.num2.Copy()
		got2.Add(i.num1)
		if want := i.sum; !reflect.DeepEqual(*want, *got) {
			t.Fatalf("#%d.2: want integer '%s' (len: %d), got '%s' (len: %d)\nDetails: '%s'+'%s'='%s' (want '%s')",
				index, want, want.Len, got2, got2.Len, i.num2, i.num1, got2, want)
		}
	}
}

func TestInt_Sub(t *testing.T) {
	input := []struct {
		num1 *Int
		num2 *Int
		diff *Int
	}{
		{num1: New(1, 1, 1, 0), num2: New(1), diff: New(1, 1, 0, 1)},
		{num1: New(1, 1, 1, 0), num2: New(0), diff: New(1, 1, 1, 0)},
		{num1: New(1, 1, 1, 1), num2: New(1), diff: New(1, 1, 1, 0)},
		{num1: New(1, 1, 1, 1), num2: New(1, 1, 1, 1), diff: New(0)},
		{num1: New(1, 1, 1, 1), num2: New(1, 1, 1, 1), diff: New(0)},
		{num1: New(0), num2: New(0), diff: New(0)},
		{num1: New(0), num2: New(0), diff: New(0)},
		{num1: New(1, 1, 1, 0), num2: New(1, 1), diff: New(1, 0, 1, 1)},
	}

	for index, i := range input {
		got := i.num1.Copy()
		got.Sub(i.num2)
		if want := i.diff; !reflect.DeepEqual(*want, *got) {
			t.Fatalf("#%d: want integer '%s' (len: %d), got '%s' (len: %d)\nDetails: '%s'-'%s'='%s' (want '%s')",
				index, want, want.Len, got, got.Len, i.num1, i.num2, got, want)
		}
	}
}

func TestInt_Mult(t *testing.T) {
	input := []struct {
		num1 uint64
		num2 uint64
		want uint64
	}{
		{1, 1, 1},
		{10, 20, 200},
		{0, 10, 0},
		{7, 7, 49},
	}

	for index, i := range input {
		got := FromUint(i.num1).Mult(FromUint(i.num2)).Uint64()
		if want := i.want; want != got {
			t.Fatalf("#%d: '%d'*'%d'='%d' (want '%d')",
				index, i.num1, i.num2, got, want)
		}
	}
}

func TestInt_Uint64(t *testing.T) {
	input := []struct {
		num1 *Int
		want uint64
	}{
		{num1: New(1, 1, 1, 0), want: 14},
		{num1: New(1, 1, 1, 1), want: 15},
		{num1: New(1, 1, 1), want: 7},
		{num1: New(0), want: 0},
	}

	for index, i := range input {
		got := i.num1.Uint64()
		if want := i.want; want != got {
			t.Fatalf("#%d: want '%d', got '%d'\nDetails: '%s' --> %d (want '%d')",
				index, want, got, i.num1, got, want)
		}
	}
}

func TestInt_Cmp(t *testing.T) {
	for i := 0; i < 1000; i++ {
		for j := 0; j < 100; j++ {
			iBig, jBig := FromUint(uint64(i)), FromUint(uint64(j))
			if i == j {
				if iBig.Cmp(jBig) != 0 || jBig.Cmp(iBig) != 0 {
					t.Fatalf("'%d' cmp '%d' want 0", i, j)
				}
			} else if i > j {
				if iBig.Cmp(jBig) != 1 || jBig.Cmp(iBig) != -1 {
					t.Fatalf("'%d' cmp '%d'", i, j)
				}
			} else {
				if iBig.Cmp(jBig) != -1 || jBig.Cmp(iBig) != 1 {
					t.Fatalf("'%d' cmp '%d'", i, j)
				}
			}
		}
	}
}

func TestInt_PushBack(t *testing.T) {
	i := New()
	i.PushBack(OneBit)
	i.PushBack(ZeroBit)
	i.PushBack(OneBit)
	if i.Uint64() != 5 {
		t.Fail()
	}
}

func TestInt_DivMod(t *testing.T) {
	for i := 1; i < 1000; i++ {
		for j := 1; j < 100; j++ {
			iBig, jBig := FromUint(uint64(i)), FromUint(uint64(j))
			div, mod := iBig.DivMod(jBig)
			if div.Uint64() != uint64(i/j) {
				t.Fatalf("%d div %d = %d (want %d)", i, j, div.Uint64(), i/j)
			}
			if mod.Uint64() != uint64(i%j) {
				t.Fatalf("%d mod %d = %d (want %d)", i, j, mod.Uint64(), i%j)
			}
		}
	}
}

func TestInt_GCD(t *testing.T) {
	a, b := FromUint(180), FromUint(150)
	fmt.Println(a.GCD(b).Uint64())
}

func TestInt_AppendUint(t *testing.T) {
	input := []struct {
		num1 *Int
		uint uint32
		want *Int
	}{
		{num1: New(0), uint: 0, want: New(1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{num1: New(1, 1, 1, 0), uint: 14, want: New(1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0)},
		{num1: New(1, 1, 1, 1), uint: 15, want: New(1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1)},
	}

	for index, i := range input {
		got := i.num1.Copy()
		got.AppendUint(i.uint)
		got.AppendBit(OneBit)
		if want := i.want; !reflect.DeepEqual(*want, *got) {
			t.Fatalf("#%d: want integer '%s' (len: %d), got '%s' (len: %d)\nDetails: '%s' append %d (%s) ='%s' (want '%s')",
				index, want, want.Len, got, got.Len, i.num1, i.uint, FromUint(uint64(i.uint)), got, want)
		}
	}
}

func TestInt_FromUint(t *testing.T) {
	for i := 0; i < 1000; i++ {
		if got := FromUint(uint64(i)).Uint64(); got != uint64(i) {
			t.Fatalf("want '%d', got '%d'", i, got)
		}
	}
}

func TestInt_Shift(t *testing.T) {
	input := []struct {
		num1  *Int
		shift uint16
		want  *Int
	}{
		{num1: New(1, 1, 1, 0), shift: 10, want: New(1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{num1: New(1, 1, 1, 0), shift: 1, want: New(1, 1, 1, 0, 0)},
		{num1: New(1, 1, 1, 0), shift: 0, want: New(1, 1, 1, 0)},
	}

	for index, i := range input {
		got := i.num1.Copy()
		got.Shift(i.shift)
		if want := i.want; !reflect.DeepEqual(*want, *got) {
			t.Fatalf("#%d: want integer '%s' (len: %d), got '%s' (len: %d)\nDetails: '%s' shift %d='%s' (want '%s')",
				index, want, want.Len, got, got.Len, i.num1, i.shift, got, want)
		}
	}
}

func TestInt_Pow(t *testing.T) {

	fmt.Println(FromUint(25).PowMod(FromUint(10000), FromUint(10)).Uint64())

	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			want := uint64(math.Pow(float64(i), float64(j)))
			if want != FromUint(uint64(i)).Pow(FromUint(uint64(j))).Uint64() {
				t.Fatal(i, j)
			}
		}
	}
}

func TestName(t *testing.T) {
	// -{2044}1010 10 {2041}1001101 77 -{2044}1010 {2041}1001101
	i := FromUint(10).ChangeSign()
	fmt.Println(i, FromUint(77), i.Copy().Add(FromUint(77)).Uint64())
	//i := FromUint(90)
	//j := FromUint(100)
	//i := FromUint(3)
	//i.ChangeSign()
	//	fmt.Println(i)
	//	fmt.Println(i.PowMod(FromUint(100), FromUint(2)).Uint64())
	//j := FromUint(20)
	//fmt.Println(i.Mult(j))
}

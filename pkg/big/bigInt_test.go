package big_test

import (
	"fmt"
	"github.com/kemics/rsa/pkg/big"
	"github.com/kemics/rsa/pkg/big/bigTest"
	"testing"
)

func TestInt_Add(t *testing.T) {
	a := bigTest.NewFrom10Text("74538967894768973476")
	b := bigTest.NewFrom10Text("74538967868973476")

	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			fmt.Println(i, j)
			stdA := bigTest.StdBig(a)
			stdB := bigTest.StdBig(b)

			got := bigTest.StdBig(a.Add(b))
			if stdA.Add(stdA, stdB).Cmp(got) != 0 {
				t.Fail()
			}

			a = a.Add(bigTest.NewFrom10Text("1"))
			b = b.Add(bigTest.NewFrom10Text("6567"))
		}
	}
}

func TestInt_Sub(t *testing.T) {
	a := bigTest.NewFrom10Text("74538967894768973476")
	b := bigTest.NewFrom10Text("74538967868973476")

	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			fmt.Println(i, j)
			stdA := bigTest.StdBig(a)
			stdB := bigTest.StdBig(b)

			got := bigTest.StdBig(a.Sub(b))
			if stdA.Sub(stdA, stdB).Cmp(got) != 0 {
				t.Fail()
			}

			a = a.Add(bigTest.NewFrom10Text("1"))
			b = b.Add(bigTest.NewFrom10Text("17"))
		}
	}

	//fmt.Println(.String())
}

func TestInt_Mult(t *testing.T) {
	a := bigTest.NewFrom10Text("74538967894768973476")
	b := bigTest.NewFrom10Text("74538967868973476")

	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			fmt.Println(i, j)
			stdA := bigTest.StdBig(a)
			stdB := bigTest.StdBig(b)

			got := bigTest.StdBig(a.Mult(b))
			if stdA.Mul(stdA, stdB).Cmp(got) != 0 {
				t.Fail()
			}

			a = a.Add(bigTest.NewFrom10Text("1"))
			b = b.Add(bigTest.NewFrom10Text("17"))
		}
	}

	//fmt.Println(.String())
}

func TestInt_Div(t *testing.T) {
	a := bigTest.NewFrom10Text("74538967894768973476")
	b := bigTest.NewFrom10Text("74533248476")

	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			fmt.Println(i, j)
			stdA := bigTest.StdBig(a)
			stdB := bigTest.StdBig(b)

			div, mod := a.DivMod(b)

			gotDiv := bigTest.StdBig(div)
			gotMod := bigTest.StdBig(mod)
			wantDiv, wantMod := stdA.DivMod(stdA, stdB, bigTest.NewStd())
			if wantDiv.Cmp(gotDiv) != 0 {
				t.Fail()
			}
			if gotMod.Cmp(wantMod) != 0 {
				t.Fail()
			}

			a = a.Add(bigTest.NewFrom10Text("1"))
			b = b.Add(bigTest.NewFrom10Text("17"))
		}
	}

	//fmt.Println(.String())
}

func TestInt_Sub2(t *testing.T) {
	fmt.Println(big.New().AppendDigits(1234).Sub(big.New().AppendDigits(342234)))
	fmt.Println(big.New().AppendDigits(1234).Sub(big.New().AppendDigits(342234).ChangeSign()))

	fmt.Println(big.New().Add(big.New().AppendDigits(342234)))
	//fmt.Println(.String())
}

func TestInt_Mult2(t *testing.T) {
	a, b := big.FromUint(2), big.FromUint(10)
	fmt.Println("RESULT", a.Mult(b))
}

func TestInt_PowMod(t *testing.T) {
	a, pow, mod := big.FromUint(2), big.FromUint(10), big.FromUint(5)
	fmt.Println("RESULT", a.PowMod(pow, mod))

	//fmt.Println(.String())
}

func TestInt_DivMod(t *testing.T) {
	a, b := big.FromUint(1024), big.FromUint(2)
	fmt.Println(a.DivMod(b))

	fmt.Println(b.DivMod(a))
}

func TestInt_DivMod2(t *testing.T) {
	fmt.Println(big.FromUint(25).DivMod(big.FromUint(10)))

	//fmt.Println(.String())
}

package big

import (
	"fmt"
	"testing"
)

func TestBig_Div(t *testing.T) {
	fmt.Println(FromUint(321197185).AppendDigits(129384902).DivMod(FromUint(2)))
}

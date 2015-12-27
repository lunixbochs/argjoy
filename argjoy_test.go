package argjoy

import (
	"strconv"
	"testing"
)

func test(a, b, optC int) int {
	return a + b + optC
}

func TestArgjoy(t *testing.T) {
	aj := NewArgjoy(StrToInt)
	aj.Optional = true

	out, err := aj.Call(test, "1", "2")
	if err != nil {
		t.Fatal(err)
	}
	if out[0].(int) != 3 {
		t.Fatalf("Incorrect result: %v != 3\n", out[0])
	}
}

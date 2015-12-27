package argjoy

import (
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

func testv(a int, b ...int) int {
	for _, v := range b {
		a += v
	}
	return a
}

func TestVarargs(t *testing.T) {
	aj := NewArgjoy(StrToInt)

	// test with varargs filled
	out, err := aj.Call(testv, "1", "2", "3", "4")
	if err != nil {
		t.Error(err)
	} else {
		if out[0].(int) != 10 {
			t.Errorf("Incorrect result: %v != 10\n", out[0])
		}
	}

	// test with varargs empty
	out, err = aj.Call(testv, "1")
	if err != nil {
		t.Error(err)
	} else {
		if out[0].(int) != 1 {
			t.Errorf("Incorrect result: %v != 1\n", out[0])
		}
	}
}

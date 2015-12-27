package argjoy

import (
	"strconv"
	"testing"
)

var StrToIntArgjoy = NewArgjoy(StrToInt)

func strToIntTest(a int, b uint, c int32, d int64, e uint32, f uint64) int64 {
	return int64(a) + int64(b) + int64(c) + d + int64(e) + int64(f)
}

func TestStrToInt(t *testing.T) {
	out, err := StrToIntArgjoy.Call(strToIntTest, "-2", "1", "-4", "-6", "2", "3")
	if err != nil {
		t.Fatal(err)
	}
	if out[0].(int64) != -6 {
		t.Fatalf("invalid result: %v != -6\n", out[0])
	}
}

func TestStrToIntBounds(t *testing.T) {
	ofi32 := strconv.FormatInt((1<<62)-1, 10)
	if _, err := StrToIntArgjoy.Call(strToIntTest, "0", "0", ofi32, "0", "0", "0"); err == nil {
		t.Error("32-bit signed integer overflow check failed")
	}
	ofu32 := strconv.FormatUint(^uint64(0), 10)
	if _, err := StrToIntArgjoy.Call(strToIntTest, "0", "0", "0", "0", ofu32, "0"); err == nil {
		t.Error("32-bit unsigned integer overflow check failed")
	}
}

func TestStrToIntSign(t *testing.T) {
	StrToIntArgjoy.Optional = true
	if _, err := StrToIntArgjoy.Call(strToIntTest, "0", "-1"); err == nil {
		t.Error("failed to throw error when parsing negative for unsigned int")
	}
}

func TestStrToIntOptional(t *testing.T) {
	StrToIntArgjoy.Optional = true
	out, err := StrToIntArgjoy.Call(strToIntTest)
	if err != nil {
		t.Fatal(err)
	}
	if out[0].(int64) != 0 {
		t.Fatalf("weird return from StrToInt optional test: %v\n", out[0])
	}
}

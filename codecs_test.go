package argjoy

import (
	"strconv"
	"testing"
)

// StrToInt tests

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

// RadStrToInt tests

var RadStrToIntArgjoy = NewArgjoy(RadStrToInt)

func radStrToIntTest(args ...int) int {
	acc := 0
	for _, v := range args {
		acc += v
	}
	return acc
}

func TestRadStrToInt(t *testing.T) {
	out, err := RadStrToIntArgjoy.Call(radStrToIntTest, "10", "0x10", "010", "0b10", "-01", "-0x1")
	if err != nil {
		t.Fatal(err)
	}
	if out[0].(int) != 34 {
		t.Fatalf("invalid result: %v != 34\n", out[0])
	}
}

func TestRadStrToIntErr(t *testing.T) {
	_, err := RadStrToIntArgjoy.Call(radStrToIntTest, "-0xFFFFFFFFFFFFFFFFFFFFFF")
	if err == nil {
		t.Error("failed to throw error on invalid rad str input")
	}
}

// IntToInt tests

var IntToIntArgjoy = NewArgjoy(IntToInt)

func intToIntTest(a int, b int8, c int16, d int32, e int64, f uint, g uint8, h uint16, i uint32, j uint64) int64 {
	return int64(a) + int64(b) + int64(c) + int64(d) + int64(e) + int64(f) + int64(g) + int64(h) + int64(i) + int64(j)
}

func TestIntToInt(t *testing.T) {
	out, err := IntToIntArgjoy.Call(intToIntTest, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	if err != nil {
		t.Fatal(err)
	}
	if out[0].(int64) != 55 {
		t.Fatalf("invalid result: %v != 55\n", out[0])
	}
}

func TestIntToIntBounds(t *testing.T) {
	IntToIntArgjoy.Optional = true
	_, err := IntToIntArgjoy.Call(intToIntTest, 0, 1000)
	if err == nil {
		t.Error("failed to throw error on int8 overflow")
	}
	_, err = IntToIntArgjoy.Call(intToIntTest, 0, -1000)
	if err == nil {
		t.Error("failed to throw error on int8 underflow")
	}
	_, err = IntToIntArgjoy.Call(intToIntTest, 0, 0, 0, 0, 0, -1000)
	if err == nil {
		t.Error("failed to throw error on uint8 underflow")
	}
	_, err = IntToIntArgjoy.Call(intToIntTest, 0, 0, 0, 0, 0, 0, 1000)
	if err == nil {
		t.Error("failed to throw error on uint8 overflow")
	}
	_, err = IntToIntArgjoy.Call(intToIntTest, 0, 0, 0, 0, ^uint64(0))
	if err == nil {
		t.Error("failed to throw error on int64 overflow with uint64_max input")
	}
	_, err = IntToIntArgjoy.Call(intToIntTest, 0, 0, 0, 0, 0, 0, 0, 0, 0, ^int64((1<<63)-1))
	if err == nil {
		t.Error("failed to throw error on uint64 underflow with int64_min input")
	}
}

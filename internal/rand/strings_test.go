package rand

import "testing"

func TestString(t *testing.T) {
	value := String(5)
	if len(value) != 5 {
		t.Fatalf("unexpected string length")
	}
	value = String(10)
	if len(value) != 10 {
		t.Fatalf("unexpected string length")
	}
	value = String(50)
	if len(value) != 50 {
		t.Fatalf("unexpected string length")
	}
}

func TestStringWithCharset(t *testing.T) {
	value := StringWithCharset(5, "a")
	if value != "aaaaa" {
		t.Fatalf("unexpected string has been generated")
	}
}

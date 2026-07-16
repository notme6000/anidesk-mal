package utils

import (
	"testing"
)

func TestRandomHex(t *testing.T) {
	h1 := RandomHex(8)
	h2 := RandomHex(8)
	if h1 == h2 {
		t.Error("expected different values")
	}
	if len(h1) != 16 {
		t.Errorf("expected 16 chars, got %d", len(h1))
	}
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !Contains(slice, "a") {
		t.Error("expected true for 'a'")
	}
	if Contains(slice, "z") {
		t.Error("expected false for 'z'")
	}
	if Contains(nil, "a") {
		t.Error("expected false for nil slice")
	}
}

func TestMin(t *testing.T) {
	if Min(3, 5) != 3 {
		t.Error("min(3,5) should be 3")
	}
	if Min(5, 3) != 3 {
		t.Error("min(5,3) should be 3")
	}
}

func TestMax(t *testing.T) {
	if Max(3, 5) != 5 {
		t.Error("max(3,5) should be 5")
	}
	if Max(5, 3) != 5 {
		t.Error("max(5,3) should be 5")
	}
}

func TestClamp(t *testing.T) {
	if Clamp(5, 0, 10) != 5 {
		t.Error("clamp(5,0,10) should be 5")
	}
	if Clamp(-1, 0, 10) != 0 {
		t.Error("clamp(-1,0,10) should be 0")
	}
	if Clamp(15, 0, 10) != 10 {
		t.Error("clamp(15,0,10) should be 10")
	}
}

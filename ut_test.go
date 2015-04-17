package utee

import (
	"strings"
	"testing"
)

func TestSplitSlice(t *testing.T) {
	a := SplitSlice(nil, 5)
	if len(a) != 1 {
		t.Error("at least 1 element")
	}
	if a[0] != nil {
		t.Error("shoud be nil")
	}

	a = SplitSlice([]string{"a"}, 4)
	if len(a) != 1 {
		t.Error("should be 1 element")
	}
	if len(a[0]) != 1 {
		t.Error("shoud be 1 elements")
	}

	a = SplitSlice([]string{"a", "b"}, 5)
	if len(a) != 1 {
		t.Error("should be 1 element")
	}
	if len(a[0]) != 2 {
		t.Error("shoud be 2 elements")
	}

	a = SplitSlice([]string{"a", "b"}, 1)
	if len(a) != 1 {
		t.Error("should be 1 element")
	}
	if len(a[0]) != 2 {
		t.Error("shoud be 2 elements")
	}

	a = SplitSlice([]string{"a", "b", "c", "d"}, 3)
	if len(a) != 3 {
		t.Error("should be 3 element")
	}
	if len(a[0]) != 2 {
		t.Error("shoud be 2 elements")
	}
	if len(a[1]) != 1 || len(a[2]) != 1 {
		t.Error("shoud be 1 elements")
	}

	a = SplitSlice([]string{"a", "b", "c", "d", "e", "f"}, 3)
	if len(a) != 3 {
		t.Error("should be 3 element")
	}
	if a[2][1] != "f" {
		t.Error("shoud be f")
	}

	a = SplitSlice([]string{"a", "b", "c"}, 2)
	if len(a) != 2 {
		t.Error("should be 2 element")
	}

	if strings.Join(a[0], "") != "ac" {
		t.Error("shoud be ac")
	}
	if strings.Join(a[1], "") != "b" {
		t.Errorf("expect be b, got %v", a[1])
	}
}

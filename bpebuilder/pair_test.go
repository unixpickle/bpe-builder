package bpebuilder

import "testing"

func TestPairMap(t *testing.T) {
	p1 := Pair{Left: []byte("a"), Right: []byte("")}
	p2 := Pair{Left: []byte(""), Right: []byte("a")}
	p3 := Pair{Left: []byte("aaa"), Right: []byte("aaa")}
	p3Other := Pair{Left: []byte("aaaa")[:3], Right: []byte("aaa")}

	m := NewPairMap[int]()
	if _, ok := m.Get(p1); ok {
		t.Error("expected not found")
	}
	m.Set(p1, 0)
	if n, ok := m.Get(p1); !ok || n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
	m.Set(p1, 1)
	if n, ok := m.Get(p1); !ok || n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	m.Set(p2, 2)
	if n, ok := m.Get(p2); !ok || n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
	if _, ok := m.Get(p3); ok {
		t.Error("expected not found")
	}
	m.Set(p3, 3)
	if n, ok := m.Get(p3); !ok || n != 3 {
		t.Errorf("expected 3, got %d", n)
	}

	// Test getting with other reference to equivalent bytes
	if n, ok := m.Get(p3Other); !ok || n != 3 {
		t.Errorf("expected 3, got %d", n)
	}

	// Test deleting
	if n, ok := m.Delete(p1); !ok || n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if _, ok := m.Delete(p1); ok {
		t.Error("expected not found")
	}
	if _, ok := m.Get(p1); ok {
		t.Error("expected not found")
	}

	// Now re-insert
	m.Set(p1, 0)
	if n, ok := m.Get(p1); !ok || n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
	m.Set(p1, 1)
	if n, ok := m.Get(p1); !ok || n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
}

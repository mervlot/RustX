package object

import "testing"

func TestEnvironment(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Integer{Value: 10})

	inner := NewClosedEnvironment(outer)
	inner.Set("y", &Integer{Value: 20})

	// Get x from inner (recursively resolved to outer)
	val, ok := inner.Get("x")
	if !ok || val.(*Integer).Value != 10 {
		t.Errorf("expected x = 10, got %v", val)
	}

	// Update x from inner (should update in outer)
	updated := inner.Update("x", &Integer{Value: 15})
	if !updated {
		t.Errorf("expected update to succeed")
	}

	// Verify outer has updated value
	val, ok = outer.Get("x")
	if !ok || val.(*Integer).Value != 15 {
		t.Errorf("expected outer x = 15, got %v", val)
	}

	// Set local variable with same name in inner
	inner.Set("x", &Integer{Value: 30})
	val, ok = inner.Get("x")
	if !ok || val.(*Integer).Value != 30 {
		t.Errorf("expected local inner x = 30, got %v", val)
	}

	// Verify outer still has 15
	val, ok = outer.Get("x")
	if !ok || val.(*Integer).Value != 15 {
		t.Errorf("expected outer x = 15, got %v", val)
	}
}

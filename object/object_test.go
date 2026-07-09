package object

import "testing"

func TestEnvironment(t *testing.T) {
	env := NewEnvironment()
	env.Set("x", &Integer{Value: 5})

	val, ok := env.Get("x")
	if !ok {
		t.Errorf("x not found in environment")
	}
	if val.(*Integer).Value != 5 {
		t.Errorf("x value wrong. got=%d", val.(*Integer).Value)
	}

	outer := env
	inner := NewEnclosedEnvironment(outer)
	inner.Set("y", &Integer{Value: 10})

	val, ok = inner.Get("y")
	if !ok {
		t.Errorf("y not found in inner environment")
	}
	if val.(*Integer).Value != 10 {
		t.Errorf("y value wrong. got=%d", val.(*Integer).Value)
	}

	val, ok = inner.Get("x")
	if !ok {
		t.Errorf("x not found in inner environment")
	}
	if val.(*Integer).Value != 5 {
		t.Errorf("x value from outer wrong. got=%d", val.(*Integer).Value)
	}

	inner.Update("x", &Integer{Value: 7})
	val, ok = outer.Get("x")
	if !ok {
		t.Errorf("x not found in outer environment after update")
	}
	if val.(*Integer).Value != 7 {
		t.Errorf("x value in outer not updated. got=%d", val.(*Integer).Value)
	}
}

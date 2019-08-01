package close

import "testing"

type mockClose struct {
}

func (m *mockClose) Close() error {
	return nil
}

func TestRegister(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Error(e)
		}

	}()

	cl := new(mockClose)
	Register(cl)
	Register(nil)
	MultiRegister(new(mockClose), nil)
	Free()
}

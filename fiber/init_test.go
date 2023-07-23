package fiber

import (
	"testing"
)

func TestNewFiber(t *testing.T) {
	if got := NewFiber(); got == nil {
		t.Errorf("fiber не инициализирован, проверять дальше нечего")
		t.Fail()
	}
}

package wiremelt

import "testing"

func TestWiremeltAscii(t *testing.T) {
	execFunc, err := WiremeltAscii()
	if err != nil {
		t.Errorf("execFunc error: %v", err)
	}

	got := execFunc
	if execFunc != got {
		t.Errorf("Expected '%v', but got '%v'", execFunc, got)
	}
}

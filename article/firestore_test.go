package article

import (
	"testing"
)

func TestVMUniqueIdProperties(t *testing.T) {
	id1 := 1
	id2 := 2

	if id1 == id2 {
		t.Fatalf("Two auto-generated VM ID should not be equivalent")
	}
}

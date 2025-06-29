package object

import (
	"testing"
)

// GOFLAGS="-count=1" go test -run TestStringHashKey
func TestStringHashKey(t *testing.T) {
	// hello1 must have same hash as hello2
	hello1 := &String{Value: "1. same value means the hash must be equal"}
	hello2 := &String{Value: "1. same value means the hash must be equal"}
	diff1 := &String{Value: "2. same value means the hash must be equal"}
	diff2 := &String{Value: "2. same value means the hash must be equal"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

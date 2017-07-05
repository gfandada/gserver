package util

import (
	"fmt"
	"testing"
)

func Test_uuid(t *testing.T) {
	u1 := NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)
	u2, err := FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Successfully parsed: %s", u2)
}

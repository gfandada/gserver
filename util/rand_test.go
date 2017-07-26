package util

import (
	"fmt"
	"testing"
)

func Test_rand(t *testing.T) {
	fmt.Println(RandInterval(1, 8))
	fmt.Println(RandIntervalN(1, 8, 3))
}

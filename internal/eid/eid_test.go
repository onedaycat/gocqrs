package eid

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	eid := New("bh9vq3krtr34lt1oedc0", 1).String()
	fmt.Println(eid)
}

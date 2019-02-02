package gocqrs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithSeq(t *testing.T) {
	tt := time.Now().Unix()
	seq := WithSeq(tt, 10)
	require.Equal(t, (tt*100000)+10, seq)
}

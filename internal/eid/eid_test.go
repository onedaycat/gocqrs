package eid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEID(t *testing.T) {
	eid := New("bh9vq3krtr34lt1oedc0", 99)
	newEid, err := FromString(eid.String())

	require.NoError(t, err)
	require.Equal(t, eid, newEid)
}

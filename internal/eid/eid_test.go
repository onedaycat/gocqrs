package eid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateEID(t *testing.T) {
	id := CreateEID("bh9vq3krtr34lt1oedc0", 99)
	require.Len(t, id, 22)
}

func TestGenerateAggregateID(t *testing.T) {
	FreezeAggregateID("1")
	id := GenerateAggregateID()
	require.Equal(t, "1", id)

	UnFreezeAggregateID()
	id = GenerateAggregateID()
	require.Len(t, id, 20)
}

package gocqrs

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	FreezeID("1")
	id := generateID()
	require.Equal(t, "1", id)

	FreezeID("")
	id = generateID()
	require.NotEmpty(t, id)
}

func TestGenerateIDSubOneSec(t *testing.T) {
	subOne := time.Now().Unix() - 1
	idnow := generateID()
	idsubOne := generateIDSubOneSec(subOne)

	comp := strings.Compare(idnow, idsubOne)
	require.Equal(t, 1, comp)
}

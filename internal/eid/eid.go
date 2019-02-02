package eid

import (
	"strconv"

	"github.com/rs/xid"
)

var (
	emptyStr  = ""
	freezeaid = ""
)

func FreezeAggregateID(id string) {
	freezeaid = id
}

func UnFreezeAggregateID() {
	freezeaid = emptyStr
}

func CreateEID(aggID string, version int64) string {
	return aggID + strconv.FormatInt(version, 10)
}

func GenerateAggregateID() string {
	if freezeaid != emptyStr {
		return freezeaid
	}

	return xid.New().String()
}

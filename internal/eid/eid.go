package eid

import (
	"encoding/base64"
)

var (
	emptyStr = ""
	freezeid = ""
)

type EID struct {
	AggregateID string
	Version     int
}

func New(aggID string, version int) *EID {
	return &EID{
		AggregateID: aggID,
		Version:     version,
	}
}

func (o *EID) String() string {
	if freezeid != emptyStr {
		return freezeid
	}

	ef := EIDFields{o.AggregateID, o.Version}
	b, _ := ef.MarshalMsg(nil)

	return base64.RawURLEncoding.EncodeToString(b)
}

func Freeze(id string) {
	freezeid = id
}

func UnFreeze() {
	freezeid = emptyStr
}

func FromString(id string) (*EID, error) {
	idb, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return nil, err
	}

	ef := EIDFields{}
	if _, err = ef.UnmarshalMsg(idb); err != nil {
		return nil, err
	}

	return &EID{
		AggregateID: ef[0].(string),
		Version:     int(ef[1].(int64)),
	}, nil
}

package gocqrs

import (
	"time"

	"github.com/rs/xid"
)

var freezeid = ""

func generateID() string {
	if freezeid != "" {
		return freezeid
	}

	return xid.New().String()
}

func FreezeID(id string) {
	freezeid = id
}

func generateIDSubOneSec(unixTime int64) string {
	return xid.NewWithTime(time.Unix(unixTime-1, 0)).String()
}

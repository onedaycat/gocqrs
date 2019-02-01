package eid

import (
	"time"
)

var dec [256]byte

var (
	emptyStr = ""
	freezeid = ""
	encoding = "0123456789abcdefghijklmnopqrstuv"
)

//go:generate msgp
type EIDFields []interface{}

type EID struct {
	AggregateID string
	Version     int
	Timestamp   int64
}

func New(aggID string, version int) *EID {
	return &EID{
		AggregateID: aggID,
		Version:     version,
		Timestamp:   time.Now().Unix(),
	}
}

func NewWithTime(aggID string, version int, timestamp int64) *EID {
	return &EID{
		AggregateID: aggID,
		Version:     version,
		Timestamp:   timestamp,
	}
}

func (o *EID) String() string {
	if freezeid != emptyStr {
		return freezeid
	}

	ef := EIDFields{o.AggregateID, o.Version, o.Timestamp}
	b, _ := ef.MarshalMsg(nil)

	text := make([]byte, len(b))
	encode(text, b)

	return string(text)
}

func FreezeID(id string) {
	freezeid = id
}

func FromString(id string) *EID {

}

func encode(dst, id []byte) {
	dst[0] = encoding[id[0]>>3]
	dst[1] = encoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]
	dst[2] = encoding[(id[1]>>1)&0x1F]
	dst[3] = encoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]
	dst[4] = encoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[5] = encoding[(id[3]>>2)&0x1F]
	dst[6] = encoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[7] = encoding[id[4]&0x1F]
	dst[8] = encoding[id[5]>>3]
	dst[9] = encoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]
	dst[10] = encoding[(id[6]>>1)&0x1F]
	dst[11] = encoding[(id[7]>>4)&0x1F|(id[6]<<4)&0x1F]
	dst[12] = encoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[13] = encoding[(id[8]>>2)&0x1F]
	dst[14] = encoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[15] = encoding[id[9]&0x1F]
	dst[16] = encoding[id[10]>>3]
	dst[17] = encoding[(id[11]>>6)&0x1F|(id[10]<<2)&0x1F]
	dst[18] = encoding[(id[11]>>1)&0x1F]
	dst[19] = encoding[(id[11]<<4)&0x1F]
}

func decode(id *ID, src []byte) {
	id[0] = dec[src[0]]<<3 | dec[src[1]]>>2
	id[1] = dec[src[1]]<<6 | dec[src[2]]<<1 | dec[src[3]]>>4
	id[2] = dec[src[3]]<<4 | dec[src[4]]>>1
	id[3] = dec[src[4]]<<7 | dec[src[5]]<<2 | dec[src[6]]>>3
	id[4] = dec[src[6]]<<5 | dec[src[7]]
	id[5] = dec[src[8]]<<3 | dec[src[9]]>>2
	id[6] = dec[src[9]]<<6 | dec[src[10]]<<1 | dec[src[11]]>>4
	id[7] = dec[src[11]]<<4 | dec[src[12]]>>1
	id[8] = dec[src[12]]<<7 | dec[src[13]]<<2 | dec[src[14]]>>3
	id[9] = dec[src[14]]<<5 | dec[src[15]]
	id[10] = dec[src[16]]<<3 | dec[src[17]]>>2
	id[11] = dec[src[17]]<<6 | dec[src[18]]<<1 | dec[src[19]]>>4
}

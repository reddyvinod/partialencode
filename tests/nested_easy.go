package tests

import (
	"github.com/reddyvinod/partialencode"
	"github.com/reddyvinod/partialencode/jwriter"
)

//partialencode:json
type NestedInterfaces struct {
	Value interface{}
	Slice []interface{}
	Map   map[string]interface{}
}

type NestedEasyMarshaler struct {
	EasilyMarshaled bool
}

var _ partialencode.Marshaler = &NestedEasyMarshaler{}

func (i *NestedEasyMarshaler) MarshalPartialJSON(w *jwriter.Writer) {
	// We use this method only to indicate that partialencode.Marshaler
	// interface was really used while encoding.
	i.EasilyMarshaled = true
}
// +build none

package template

import (
	"fmt"

	"github.com/reddyvinod/partialencode/jlexer"
	"github.com/reddyvinod/partialencode/jwriter"
)

// template type Optional(A)
type A int

// A 'gotemplate'-based type for providing optional semantics without using pointers.
type Optional struct {
	Value A
	Valid bool
	Set   bool
}

func (v *Optional) SetValue(val A) {
	v.Value = val
	v.Set = true
	v.Valid = true
}

func (v *Optional) SetNull() {
	v.Set = true
	v.Valid = false
}

func (v *Optional) IsValid() bool {
	return v.Valid
}

func (v *Optional) IsSet() bool {
	return v.Set
}

func (v *Optional) Get() A {
	return v.Value
}

// MarshalPartialJSON does JSON marshaling using partialencode interface.
func (v Optional) MarshalPartialJSON(w *jwriter.Writer) {
	if v.Valid {
		w.Optional(v.Value)
	} else {
		w.RawString("null")
	}
}

// UnMarshalPartialJSON does JSON unmarshaling using partialencode interface.
func (v *Optional) UnMarshalPartialJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		l.Skip()
		v.SetNull()
	} else {
		v.SetValue(l.Optional())
	}
}

// MarshalJSON implements a standard json marshaler interface.
func (v Optional) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	v.MarshalPartialJSON(&w)
	return w.Buffer.BuildBytes(), w.Error
}

// UnmarshalJSON implements a standard json unmarshaler interface.
func (v *Optional) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{Data: data}
	v.UnMarshalPartialJSON(&l)
	return l.Error()
}

// String implements a stringer interface using fmt.Sprint for the value.
func (v Optional) String() string {
	if !v.Set {
		return "<undefined>"
	}
	if !v.Valid {
		return "null"
	}
	return fmt.Sprint(v.Value)
}

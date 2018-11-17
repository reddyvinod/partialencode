// generated by gotemplate

package basic

import (
	"fmt"

	"github.com/reddyvinod/partialencode/jlexer"
	"github.com/reddyvinod/partialencode/jwriter"
)

// template type Uint32(uint32)

// uint32 'gotemplate'-based type for providing optional semantics without using pointers.
type Uint32 struct {
	Value uint32
	Valid bool
	Set   bool
}

func (v *Uint32) SetValue(val uint32) {
	v.Value = val
	v.Set = true
	v.Valid = true
}

func (v *Uint32) SetNull() {
	v.Set = true
	v.Valid = false
}

func (v *Uint32) IsValid() bool {
	return v.Valid
}

func (v *Uint32) IsSet() bool {
	return v.Set
}

func (v *Uint32) Get() uint32 {
	return v.Value
}

// MarshalPartialJSON does JSON marshaling using partialencode interface.
func (v Uint32) MarshalPartialJSON(w *jwriter.Writer) {
	if v.Valid {
		w.Uint32(v.Value)
	} else {
		w.RawString("null")
	}
}

// UnMarshalPartialJSON does JSON unmarshaling using partialencode interface.
func (v *Uint32) UnMarshalPartialJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		l.Skip()
		v.SetNull()
	} else {
		v.SetValue(l.Uint32())
	}
}

// MarshalJSON implements a standard json marshaler interface.
func (v Uint32) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	v.MarshalPartialJSON(&w)
	return w.Buffer.BuildBytes(), w.Error
}

// UnmarshalJSON implements a standard json unmarshaler interface.
func (v *Uint32) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{Data: data}
	v.UnMarshalPartialJSON(&l)
	return l.Error()
}

// String implements a stringer interface using fmt.Sprint for the value.
func (v Uint32) String() string {
	if !v.Set {
		return "<undefined>"
	}
	if !v.Valid {
		return "null"
	}
	return fmt.Sprint(v.Value)
}

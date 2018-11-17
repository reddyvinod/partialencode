// generated by gotemplate

package opt

import (
	"fmt"

	"github.com/reddyvinod/partialencode/jlexer"
	"github.com/reddyvinod/partialencode/jwriter"
)

// template type Optional(A)

// A 'gotemplate'-based type for providing optional semantics without using pointers.
type String struct {
	V       string
	Defined bool
}

// Creates an optional type with a given value.
func OString(v string) String {
	return String{V: v, Defined: true}
}

// Get returns the value or given default in the case the value is undefined.
func (v String) Get(deflt string) string {
	if !v.Defined {
		return deflt
	}
	return v.V
}

// MarshalPartialJSON does JSON marshaling using partialencode interface.
func (v String) MarshalPartialJSON(w *jwriter.Writer) {
	if v.Defined {
		w.String(v.V)
	} else {
		w.RawString("null")
	}
}

// UnMarshalPartialJSON does JSON unmarshaling using partialencode interface.
func (v *String) UnMarshalPartialJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		l.Skip()
		*v = String{}
	} else {
		v.V = l.String()
		v.Defined = true
	}
}

// MarshalJSON implements a standard json marshaler interface.
func (v String) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	v.MarshalPartialJSON(&w)
	return w.Buffer.BuildBytes(), w.Error
}

// UnmarshalJSON implements a standard json unmarshaler interface.
func (v *String) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{Data: data}
	v.UnMarshalPartialJSON(&l)
	return l.Error()
}

// IsDefined returns whether the value is defined, a function is required so that it can
// be used in an interface.
func (v String) IsDefined() bool {
	return v.Defined
}

// String implements a stringer interface using fmt.Sprint for the value.
func (v String) String() string {
	if !v.Defined {
		return "<undefined>"
	}
	return fmt.Sprint(v.V)
}

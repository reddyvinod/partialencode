package gen

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var basicStructs = map[reflect.Kind]string{
	reflect.String:  "basic.String",
	reflect.Bool:    "basic.Bool",
	reflect.Int:     "basic.Int",
	reflect.Int8:    "basic.Int8",
	reflect.Int16:   "basic.Int16",
	reflect.Int32:   "basic.Int32",
	reflect.Int64:   "basic.Int64",
	reflect.Uint:    "basic.Uint",
	reflect.Uint8:   "basic.Uint8",
	reflect.Uint16:  "basic.Uint16",
	reflect.Uint32:  "basic.Uint32",
	reflect.Uint64:  "basic.Uint64",
	reflect.Float32: "basic.Float32",
	reflect.Float64: "basic.Float64",
}

func (g *PartialGenerator) getStructName(t reflect.Type) string {
	return g.structName("Partial", t)
}

func (g *PartialGenerator) genPartialStruct(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Struct:
		return g.genStructPartialStruct(t)
	default:
		return g.genNonStructPartialStruct(t)
	}
}

func (g *PartialGenerator) genNonStructPartialStruct(t reflect.Type) error {
	if t.Kind() == reflect.Struct {
		return fmt.Errorf("cannot generate encoder/decoder for %v, it is a struct type", t)
	}

	sname := g.getStructName(t)

	fmt.Fprintln(g.out, "type "+sname+" struct {")
	fmt.Fprint(g.out, "  Value ")
	g.genTypePartial(t, true, 2)
	fmt.Fprintln(g.out, "  Valid bool")
	fmt.Fprintln(g.out, "  Set bool")
	fmt.Fprintln(g.out, "}")

	return nil
}

func (g *PartialGenerator) genStructPartialStruct(t reflect.Type) error {
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a struct type", t)
	}

	sname := g.getStructName(t)

	fmt.Fprintln(g.out, "type "+sname+" struct {")
	fmt.Fprintln(g.out, "  Value struct {")
	for i := 0; i < t.NumField(); i++ {
		g.genFieldPartialStruct(t.Field(i), 2)
	}
	fmt.Fprintln(g.out, "  }")
	fmt.Fprintln(g.out, "  Valid bool")
	fmt.Fprintln(g.out, "  Set bool")
	fmt.Fprintln(g.out, "}")
	fmt.Fprintln(g.out, "")

	return nil
}

func (g *PartialGenerator) genFieldPartialStruct(f reflect.StructField, indent int) {
	if f.Anonymous {
		panic("anonymous fields are not supported")
	}
	ws := strings.Repeat("  ", indent)

	fmt.Fprint(g.out, ws+f.Name+" ")
	g.genTypePartial(f.Type, false, indent+1)
	if len(string(f.Tag)) > 0 {
		fmt.Fprintln(g.out, " `"+string(f.Tag)+"`")
	} else {
		fmt.Fprintln(g.out, "")
	}
}

// If the type/it's elem is not anonymous it will return the type name and true
// Else it will return an empty string and false
func (g *PartialGenerator) genTypePartial(t reflect.Type, basic bool, indent int) {
	ws := strings.Repeat("  ", indent)
	// non-defined and pre-defined types and anonymous types
	if t.PkgPath() == "" {
		// pre-defined types
		if t.Name() != "" {
			if basic {
				fmt.Fprint(g.out, t.Kind())
			} else {
				fmt.Fprint(g.out, basicStructs[t.Kind()])
			}
		}

		// composite/non-defined types
		switch t.Kind() {
		case reflect.Ptr:
			fmt.Fprint(g.out, "*")
			g.genTypePartial(t.Elem(), true, indent+1)
		case reflect.Slice:
			fmt.Fprintln(g.out, " struct {")
			fmt.Fprint(g.out, ws+"  Value []")
			g.genTypePartial(t.Elem(), true, indent+1)
			fmt.Fprintln(g.out, "")
			fmt.Fprintln(g.out, ws+"  Valid bool")
			fmt.Fprintln(g.out, ws+"  Set bool")
			fmt.Fprint(g.out, ws+"}")
		case reflect.Array:
			fmt.Fprintln(g.out, " struct {")
			fmt.Fprint(g.out, ws+"Value ["+strconv.Itoa(t.Len())+"]")
			g.genTypePartial(t.Elem(), true, indent+1)
			fmt.Fprintln(g.out, "")
			fmt.Fprintln(g.out, ws+"  Valid bool")
			fmt.Fprintln(g.out, ws+"  Set bool")
			fmt.Fprint(g.out, ws+"}")
		case reflect.Map:
			fmt.Fprintln(g.out, " struct {")
			fmt.Fprint(g.out, ws+" Value map[")
			g.genTypePartial(t.Key(), true, indent+1)
			fmt.Fprint(g.out, "]")
			g.genTypePartial(t.Elem(), true, indent+1)
			fmt.Fprintln(g.out, "")
			fmt.Fprintln(g.out, ws+"  Valid bool")
			fmt.Fprintln(g.out, ws+"  Set bool")
			fmt.Fprint(g.out, ws+"}")
		case reflect.Struct:
			fmt.Fprintln(g.out, " struct {")
			fmt.Fprint(g.out, ws+"  Value struct {")
			for i := 0; i < t.NumField(); i++ {
				g.genFieldPartialStruct(t.Field(i), indent+1)
			}
			fmt.Fprintln(g.out, ws+"  }")
			fmt.Fprintln(g.out, ws+"  Valid bool")
			fmt.Fprintln(g.out, ws+"  Set bool")
			fmt.Fprint(g.out, ws+"}")
		}
	} else if t.PkgPath() == g.pkgPath {
		g.addType(t)
		fmt.Fprint(g.out, g.getStructName(t))
	} else {
		g.addExternalType(t)
		g.pkgAlias(t.PkgPath())
		fmt.Fprint(g.out, g.getStructName(t))
	}
}

func (g *PartialGenerator) genExternalTypePartialStruct(t reflect.Type) error {
	sname := g.getStructName(t)

	fmt.Fprintln(g.out, "type "+sname+" struct {")
	fmt.Fprintln(g.out, "  Value ", g.pkgAlias(t.PkgPath())+"."+t.Name())
	fmt.Fprintln(g.out, "  Valid bool")
	fmt.Fprintln(g.out, "  Set bool")
	fmt.Fprintln(g.out, "}")

	return nil
}

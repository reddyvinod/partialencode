package gen

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const PartialValidKey = "PartialValid"
const PartialSetKey = "PartialSet"

func (g *PartialGenerator) getStructName(t reflect.Type) string {
	return g.structName("Partial", t)
}

func (g *PartialGenerator) getBoolStructName(t reflect.Type) string {
	return g.structName("PartialBool", t)
}

func (g *PartialGenerator) genPartialStruct(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Struct:
		return g.genStructPartialStruct(t)
	default:
		return fmt.Errorf("cannot generate partial struct for %v, it is not struct type", t)
	}
}

func (g *PartialGenerator) genStructPartialStruct(t reflect.Type) error {
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a struct type", t)
	}

	sname := g.getStructName(t)
	bname := g.getBoolStructName(t)

	fmt.Fprintln(g.out, "type "+sname+" struct {")
	for i := 0; i < t.NumField(); i++ {
		g.genFieldPartialStruct(t.Field(i), 1)
	}
	fmt.Fprintln(g.out, "  "+PartialValidKey+" "+bname+"`bson:\"-\" json:\"-\"`")
	fmt.Fprintln(g.out, "  "+PartialSetKey+" "+bname+"`bson:\"-\" json:\"-\"`")
	fmt.Fprintln(g.out, "}")
	fmt.Fprintln(g.out, "")

	return nil
}

func (g *PartialGenerator) genFieldPartialStruct(f reflect.StructField, indent int) {
	ws := strings.Repeat("  ", indent)

	fmt.Fprint(g.out, ws+f.Name+" ")
	if !f.Anonymous {
		g.genTypePartial(f.Type, indent+1)
	}
	if len(string(f.Tag)) > 0 {
		fmt.Fprintln(g.out, " `"+string(f.Tag)+"`")
	} else {
		fmt.Fprintln(g.out, "")
	}
}

// If the type/it's elem is not anonymous it will return the type name and true
// Else it will return an empty string and false
func (g *PartialGenerator) genTypePartial(t reflect.Type, indent int) {
	ws := strings.Repeat("  ", indent)
	// non-defined and pre-defined types and anonymous types
	if t.PkgPath() == "" {
		// pre-defined types
		if t.Name() != "" {
			fmt.Fprint(g.out, t.Kind())
		}

		// composite/non-defined types
		switch t.Kind() {
		case reflect.Ptr:
			fmt.Fprint(g.out, "*")
			g.genTypePartial(t.Elem(), indent+1)
		case reflect.Slice:
			fmt.Fprint(g.out, " []")
			g.genTypePartial(t.Elem(), indent+1)
		case reflect.Array:
			fmt.Fprint(g.out, " ["+strconv.Itoa(t.Len())+"]")
			g.genTypePartial(t.Elem(), indent+1)
		case reflect.Map:
			fmt.Fprint(g.out, " map[")
			g.genTypePartial(t.Key(), indent+1)
			fmt.Fprint(g.out, "]")
			g.genTypePartial(t.Elem(), indent+1)
		case reflect.Struct:
			fmt.Fprint(g.out, " struct {")
			for i := 0; i < t.NumField(); i++ {
				g.genFieldPartialStruct(t.Field(i), indent+1)
			}
			fmt.Fprintln(g.out, ws+"  "+PartialValidKey+" struct {")
			for i := 0; i < t.NumField(); i++ {
				fmt.Fprintln(g.out, ws+"    "+t.Field(i).Name+" bool")
			}
			fmt.Fprintln(g.out, ws+"  } `bson:\"-\" json:\"-\"`")
			fmt.Fprintln(g.out, ws+"  "+PartialSetKey+" struct {")
			for i := 0; i < t.NumField(); i++ {
				fmt.Fprintln(g.out, ws+"    "+t.Field(i).Name+" bool")
			}
			fmt.Fprintln(g.out, ws+"  } `bson:\"-\" json:\"-\"`")
			fmt.Fprint(g.out, ws+"}")
		}
	} else if t.PkgPath() == g.pkgPath {
		fmt.Fprint(g.out, g.getStructName(t))
	} else {
		fmt.Fprint(g.out, g.pkgAlias(t.PkgPath())+"."+t.Name())
	}
}

func (g *PartialGenerator) genPartialBoolStruct(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Struct:
		return g.genStructPartialBoolStruct(t)
	default:
		return fmt.Errorf("cannot generate partial bool struct for %v, it is not struct type", t)
	}
}

func (g *PartialGenerator) genStructPartialBoolStruct(t reflect.Type) error {
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a struct type", t)
	}

	bname := g.getBoolStructName(t)

	fmt.Fprintln(g.out, "type "+bname+" struct {")
	for i := 0; i < t.NumField(); i++ {
		fmt.Fprintln(g.out, "  "+t.Field(i).Name+" bool")
	}
	fmt.Fprintln(g.out, "}")
	fmt.Fprintln(g.out, "")

	return nil
}

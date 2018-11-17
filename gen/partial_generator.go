package gen

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"path"
	"reflect"
	"sort"
	"strings"
)

const basic = "github.com/reddyvinod/partialencode/basic"

// PartialGenerator generates the requested marshaler/unmarshalers.
type PartialGenerator struct {
	out *bytes.Buffer

	pkgName    string
	pkgPath    string
	buildTags  string
	hashString string

	varCounter int

	noStdMarshalers       bool
	omitEmpty             bool
	disallowUnknownFields bool
	fieldNamer            FieldNamer

	// package path to local alias map for tracking imports
	imports map[string]string

	// types that partials were requested for by user
	partials map[reflect.Type]bool

	// types that encoders were already generated for
	typesSeen map[reflect.Type]bool

	// types that encoders were requested for (e.g. by encoders of other types)
	typesUnseen []reflect.Type

	// types that encoders were already generated for
	externalTypesSeen map[reflect.Type]bool

	// types that encoders were requested for (e.g. by encoders of other types)
	externalTypesUnseen []reflect.Type

	// struct name to relevant type maps to track names of partial-struct in
	// case of a name clash or unnamed structs
	structNames map[string]reflect.Type
}

// SetPkg sets the name and path of output package.
func (g *PartialGenerator) SetPkg(name, path string) {
	g.pkgName = name
	g.pkgPath = path
}

// SetBuildTags sets build tags for the output file.
func (g *PartialGenerator) SetBuildTags(tags string) {
	g.buildTags = tags
}

// addTypes requests to generate encoding/decoding funcs for the given type.
func (g *PartialGenerator) addType(t reflect.Type) {
	if g.typesSeen[t] {
		return
	}
	for _, t1 := range g.typesUnseen {
		if t1 == t {
			return
		}
	}
	g.typesUnseen = append(g.typesUnseen, t)
}

// addTypes requests to generate encoding/decoding funcs for the given type.
func (g *PartialGenerator) addExternalType(t reflect.Type) {
	if g.externalTypesSeen[t] {
		return
	}
	for _, t1 := range g.externalTypesUnseen {
		if t1 == t {
			return
		}
	}
	g.externalTypesUnseen = append(g.externalTypesUnseen, t)
}

// Add requests to generate marshaler/unmarshalers and encoding/decoding
// funcs for the type of given object.
func (g *PartialGenerator) Add(obj interface{}) {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	g.addType(t)
}

func NewPartialGenerator(filename string) *PartialGenerator {
	ret := &PartialGenerator{
		imports: map[string]string{
			basic: "basic",
		},
		typesSeen:         make(map[reflect.Type]bool),
		externalTypesSeen: make(map[reflect.Type]bool),
		structNames:       make(map[string]reflect.Type),
	}

	// Use a file-unique prefix on all auxiliary funcs to avoid
	// name clashes.
	hash := fnv.New32()
	hash.Write([]byte(filename))
	ret.hashString = fmt.Sprintf("%x", hash.Sum32())

	return ret
}

// Run runs the generator and outputs generated code to out.
func (g *PartialGenerator) Run(out io.Writer) error {
	g.out = &bytes.Buffer{}

	for len(g.typesUnseen) > 0 {
		t := g.typesUnseen[len(g.typesUnseen)-1]
		g.typesUnseen = g.typesUnseen[:len(g.typesUnseen)-1]
		g.typesSeen[t] = true

		if err := g.genPartialStruct(t); err != nil {
			return err
		}
	}

	for len(g.externalTypesUnseen) > 0 {
		t := g.externalTypesUnseen[len(g.externalTypesUnseen)-1]
		g.externalTypesUnseen = g.externalTypesUnseen[:len(g.externalTypesUnseen)-1]
		g.externalTypesSeen[t] = true

		if err := g.genExternalTypePartialStruct(t); err != nil {
			return err
		}
	}
	g.printStructsHeader()
	_, err := out.Write(g.out.Bytes())
	return err
}

// printHeader prints package declaration and imports.
func (g *PartialGenerator) printStructsHeader() {
	if g.buildTags != "" {
		fmt.Println("// +build ", g.buildTags)
		fmt.Println()
	}
	fmt.Println("// Code generated by partial for partial-structs. DO NOT EDIT.")
	fmt.Println()
	fmt.Println("package ", g.pkgName)
	fmt.Println()

	byAlias := map[string]string{}
	var aliases []string
	for path, alias := range g.imports {
		aliases = append(aliases, alias)
		byAlias[alias] = path
	}

	sort.Strings(aliases)
	fmt.Println("import (")
	for _, alias := range aliases {
		fmt.Printf("  %s %q\n", alias, byAlias[alias])
	}

	fmt.Println(")")
	fmt.Println("")
	fmt.Println("// suppress unused package warning")
	fmt.Println("var (")
	fmt.Println("   _ basic.Int")
	fmt.Println(")")

	fmt.Println()
}

// functionName returns a function name for a given type with a given prefix. If a function
// with this prefix already exists for a type, it is returned.
//
// Method is used to track encoder/decoder names for the type.
func (g *PartialGenerator) structName(prefix string, t reflect.Type) string {
	prefix = joinFunctionNameParts(false, prefix, g.hashString)
	name := joinFunctionNameParts(false, prefix, t.Name())

	// Most of the names will be unique, try a shortcut first.
	if e, ok := g.structNames[name]; !ok || e == t {
		g.structNames[name] = t
		return name
	}

	// Search if the strut already exists.
	for name1, t1 := range g.structNames {
		if t1 == t && strings.HasPrefix(name1, prefix) {
			return name1
		}
	}

	// Create a new name in the case of a clash.
	for i := 1; ; i++ {
		nm := fmt.Sprint(name, i)
		if _, ok := g.structNames[nm]; ok {
			continue
		}
		g.structNames[nm] = t
		return nm
	}
}

// pkgAlias creates and returns and import alias for a given package.
func (g *PartialGenerator) pkgAlias(pkgPath string) string {
	pkgPath = fixPkgPathVendoring(pkgPath)
	if alias := g.imports[pkgPath]; alias != "" {
		return alias
	}

	for i := 0; ; i++ {
		alias := fixAliasName(path.Base(pkgPath))
		if i > 0 {
			alias += fmt.Sprint(i)
		}

		exists := false
		for _, v := range g.imports {
			if v == alias {
				exists = true
				break
			}
		}

		if !exists {
			g.imports[pkgPath] = alias
			return alias
		}
	}
}
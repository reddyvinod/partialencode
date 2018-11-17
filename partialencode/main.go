package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/reddyvinod/partialencode/bootstrap"
	// Reference the gen package to be friendly to vendoring tools,
	// as it is an indirect dependency.
	// (The temporary bootstrapping code uses it.)
	_ "github.com/reddyvinod/partialencode/gen"
	"github.com/reddyvinod/partialencode/parser"
)

var buildTags = flag.String("build_tags", "", "build tags to add to generated file")
var snakeCase = flag.Bool("snake_case", false, "use snake_case names instead of CamelCase by default")
var lowerCamelCase = flag.Bool("lower_camel_case", false, "use lowerCamelCase names instead of CamelCase by default")
var noStdMarshalers = flag.Bool("no_std_marshalers", false, "don't generate MarshalJSON/UnmarshalJSON funcs")
var omitEmpty = flag.Bool("omit_empty", false, "omit empty fields by default")
var allStructs = flag.Bool("all", false, "generate marshaler/unmarshalers for all structs in a file")
var leaveTemps = flag.Bool("leave_temps", false, "do not delete temporary files")
var stubs = flag.Bool("stubs", false, "only generate stubs for marshaler/unmarshaler funcs")
var noformat = flag.Bool("noformat", false, "do not run 'gofmt -w' on output file")
var partialSpecifiedName = flag.String("partial_filename", "", "specify the filename of the partial structs output")
var deencoderSpecifiedName = flag.String("de_encode_filename", "", "specify the filename of the de/encoders output")
var processPkg = flag.Bool("pkg", false, "process the whole package instead of just the given file")
var disallowUnknownFields = flag.Bool("disallow_unknown_fields", false, "return error if any unknown field in json appeared")

func generate(fname string) (err error) {
	var partialName string
	{
		fInfo, err := os.Stat(fname)
		if err != nil {
			return err
		}

		p := parser.Parser{AllStructs: *allStructs}
		if err := p.Parse(fname, fInfo.IsDir()); err != nil {
			return fmt.Errorf("Error parsing %v: %v", fname, err)
		}

		if fInfo.IsDir() {
			partialName = filepath.Join(fname, p.PkgName+"_partial.go")
		} else {
			if s := strings.TrimSuffix(fname, ".go"); s == fname {
				return errors.New("Filename must end in '.go'")
			} else {
				partialName = s + "_partial.go"
			}
		}
		if *partialSpecifiedName != "" {
			partialName = *partialSpecifiedName
		}

		var trimmedBuildTags string
		if *buildTags != "" {
			trimmedBuildTags = strings.TrimSpace(*buildTags)
		}

		g := bootstrap.Generator{
			BuildTags:             trimmedBuildTags,
			PkgPath:               p.PkgPath,
			PkgName:               p.PkgName,
			Types:                 p.StructNames,
			SnakeCase:             *snakeCase,
			LowerCamelCase:        *lowerCamelCase,
			NoStdMarshalers:       *noStdMarshalers,
			DisallowUnknownFields: *disallowUnknownFields,
			OmitEmpty:             *omitEmpty,
			LeaveTemps:            *leaveTemps,
			PartialName:           partialName,
			StubsOnly:             *stubs,
			NoFormat:              *noformat,
		}

		if err := g.RunPartial(); err != nil {
			return fmt.Errorf("Bootstrap failed: %v", err)
		}
	}
	{
		fInfo, err := os.Stat(partialName)
		if err != nil {
			return err
		}

		p := parser.Parser{AllStructs: *allStructs}
		if err := p.Parse(partialName, fInfo.IsDir()); err != nil {
			return fmt.Errorf("Error parsing %v: %v", partialName, err)
		}

		var deEncoderName string
		if fInfo.IsDir() {
			deEncoderName = filepath.Join(partialName, p.PkgName+"_de_encoder.go")
		} else {
			if s := strings.TrimSuffix(partialName, ".go"); s == partialName {
				return errors.New("Filename must end in '.go'")
			} else {
				s = strings.TrimSuffix(s, "_partial")
				deEncoderName = s + "_de_encoder.go"
			}
		}

		if *deencoderSpecifiedName != "" {
			deEncoderName = *deencoderSpecifiedName
		}

		var trimmedBuildTags string
		if *buildTags != "" {
			trimmedBuildTags = strings.TrimSpace(*buildTags)
		}
		g := bootstrap.Generator{
			BuildTags:             trimmedBuildTags,
			PkgPath:               p.PkgPath,
			PkgName:               p.PkgName,
			Types:                 p.StructNames,
			SnakeCase:             *snakeCase,
			LowerCamelCase:        *lowerCamelCase,
			NoStdMarshalers:       *noStdMarshalers,
			DisallowUnknownFields: *disallowUnknownFields,
			OmitEmpty:             *omitEmpty,
			LeaveTemps:            *leaveTemps,
			DeEncoderName:         deEncoderName,
			StubsOnly:             *stubs,
			NoFormat:              *noformat,
		}

		if err := g.RunDeEncode(); err != nil {
			return fmt.Errorf("Bootstrap failed: %v", err)
		}
	}

	return nil
}

func main() {
	flag.Parse()

	files := flag.Args()

	gofile := os.Getenv("GOFILE")
	if *processPkg {
		gofile = filepath.Dir(gofile)
	}

	if len(files) == 0 && gofile != "" {
		files = []string{gofile}
	} else if len(files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, fname := range files {
		if err := generate(fname); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

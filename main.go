package main

import (
	"errors"
	"flag"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	flagTargetPkg          string
	flagIncludeImportedPkg bool
)

var (
	errorIface = types.Universe.Lookup("error").(*types.TypeName)
)

func interfacesFromPackage(patterns ...string) ([]types.Object, error) {
	mode := packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedImports
	cfg := &packages.Config{Mode: mode}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	if len(pkgs) <= 0 {
		return nil, errors.New("very few packages")
	}

	ifaces := []types.Object{}

	// selected package
	for _, name := range pkgs[0].Types.Scope().Names() {
		obj := pkgs[0].Types.Scope().Lookup(name)
		if obj == nil {
			continue
		}

		if !types.IsInterface(obj.Type()) {
			continue
		}

		ifaces = append(ifaces, obj)
	}

	// imported pacakges
	for _, p := range pkgs[0].Imports {
		for _, name := range p.Types.Scope().Names() {
			obj := p.Types.Scope().Lookup(name)
			if obj == nil {
				continue
			}

			if !types.IsInterface(obj.Type()) {
				continue
			}

			ifaces = append(ifaces, obj)
		}
	}

	// Universe Scoop
	ifaces = append(ifaces, errorIface)

	return ifaces, nil
}

func typeObjFromName(pkg string, name string) (types.Object, error) {
	mode := packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedImports
	cfg := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	obj := pkgs[0].Types.Scope().Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("lookup: not found type %s.%s", pkg, name)
	}

	return obj, nil
}

func run(args []string) error {
	if len(args) < 1 {
		return errors.New("invalid arguments")
	}

	targetType := args[0]

	typ := strings.TrimLeft(filepath.Ext(targetType), ".")
	pkg := strings.TrimRight(strings.TrimSuffix(targetType, typ), ".")
	if pkg == "" || typ == "" {
		return errors.New("invalid type name")
	}

	ifs, err := interfacesFromPackage(pkg)
	if err != nil {
		return err
	}

	obj, err := typeObjFromName(pkg, typ)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", obj.Type())
	for _, iface := range ifs {
		i, ok := iface.Type().Underlying().(*types.Interface)
		if i == nil || !ok {
			return errors.New("invalid interface")
		}
		if types.Implements(obj.Type(), i) || types.Implements(types.NewPointer(obj.Type()), i) {
			fmt.Printf("\t%s\n", iface.Type())
		}
	}

	return nil
}

func init() {
	flag.StringVar(&flagTargetPkg, "t", "", "packages included in the search (',' separated list)")
	flag.BoolVar(&flagIncludeImportedPkg, "c", true, "")
	flag.Parse()
}

func main() {
	if err := run(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

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
	flagType string
)

func interfacesFromPackage(patterns ...string) ([]*types.TypeName, error) {
	mode := packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedImports
	cfg := &packages.Config{Mode: mode}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	if len(pkgs) <= 0 {
		return nil, errors.New("very few packages")
	}

	ifaces := []*types.TypeName{}

	spkgs := make([]*packages.Package, 0, len(pkgs[0].Imports)+1)
	spkgs = append(spkgs, pkgs[0])
	for _, imp := range pkgs[0].Imports {
		spkgs = append(spkgs, imp)
	}

	for _, p := range spkgs {
		for _, name := range p.Types.Scope().Names() {
			obj, ok := p.Types.Scope().Lookup(name).(*types.TypeName)
			if obj == nil || !ok {
				continue
			}

			if !types.IsInterface(obj.Type()) {
				continue
			}

			ifaces = append(ifaces, obj)
		}
	}

	// Universe Scoop
	obj, _ := types.Universe.Lookup("error").(*types.TypeName)
	ifaces = append(ifaces, obj)

	return ifaces, nil
}

func typeNameObj(pkg string, name string) (*types.TypeName, error) {
	mode := packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedImports
	cfg := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	obj, ok := pkgs[0].Types.Scope().Lookup(name).(*types.TypeName)
	if obj == nil || !ok {
		return nil, fmt.Errorf("lookup: not found type %s.%s", pkg, name)
	}

	return obj, nil
}

func run(args []string) error {
	if flagType == "" {
		return errors.New("must set -t flag")
	}

	typ := strings.TrimLeft(filepath.Ext(flagType), ".")
	pkg := strings.TrimRight(strings.TrimSuffix(flagType, typ), ".")
	if pkg == "" || typ == "" {
		return errors.New("invalid type name")
	}

	ifs, err := interfacesFromPackage(pkg)
	if err != nil {
		return err
	}

	obj, err := typeNameObj(pkg, typ)
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
	flag.StringVar(&flagType, "t", "", "type name")
	flag.Parse()
}

func main() {
	if err := run(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

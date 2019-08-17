package dsl

import (
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/dslengine"
)

type Root struct {
	Admins   []*AdminDefinition
	AdminMap map[string]bool
}

func (r *Root) DSLName() string {
	return "yamlgen"
}

func (r *Root) Context() string {
	return "yamlgen"
}

func (r *Root) DependsOn() []dslengine.Root {
	return []dslengine.Root{
		design.Design,
	}
}

func (r *Root) IterateSets(it dslengine.SetIterator) {
	for _, v := range r.Admins {
		it([]dslengine.Definition{v})
	}
}

func (r *Root) Reset() {
	r.Admins = []*AdminDefinition{}
	r.AdminMap = map[string]bool{}
}

var DslRoot = &Root{
	Admins:   []*AdminDefinition{},
	AdminMap: map[string]bool{},
}

func init() {
	dslengine.Register(DslRoot)
}

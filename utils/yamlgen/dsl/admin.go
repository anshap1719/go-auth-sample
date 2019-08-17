package dsl

import (
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/dslengine"
)

type AdminDefinition struct {
	Action *design.ActionDefinition
}

func AdminOnly() {
	action, ok := dslengine.CurrentDefinition().(*design.ActionDefinition)
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}
	DslRoot.Admins = append(DslRoot.Admins, &AdminDefinition{action})
}

func (a *AdminDefinition) Context() string {
	return "admin"
}

func (a *AdminDefinition) Finalize() {
	for _, v := range a.Action.Routes {
		DslRoot.AdminMap[v.FullPath()] = true
	}
}

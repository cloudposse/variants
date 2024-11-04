// https://forum.golangbridge.org/t/html-template-optional-argument-in-function/6080
// https://lkumarjain.blogspot.com/2020/11/deep-dive-into-go-template.html
// https://echorand.me/posts/golang-templates/
// https://www.practical-go-lessons.com/chap-32-templates
// https://docs.gofiber.io/template/next/html/TEMPLATES_CHEATSHEET/
// https://engineering.01cloud.com/2023/04/13/optional-function-parameter-pattern/

package exec

import (
	"context"
	"text/template"

	"github.com/cloudposse/atmos/pkg/schema"
	"github.com/hairyhenderson/gomplate/v3/data"
)

// FuncMap creates and returns a map of template functions
func FuncMap(cliConfig schema.CliConfiguration, ctx context.Context, gomplateData *data.Data) template.FuncMap {
	atmosFuncs := &AtmosFuncs{cliConfig, ctx, gomplateData}

	return map[string]any{
		"atmos": func() any { return atmosFuncs },
	}
}

type AtmosFuncs struct {
	cliConfig    schema.CliConfiguration
	ctx          context.Context
	gomplateData *data.Data
}

func (f AtmosFuncs) Component(component string, stack string) (any, error) {
	return componentFunc(f.cliConfig, component, stack)
}

func (f AtmosFuncs) GomplateDatasource(alias string, args ...string) (any, error) {
	return gomplateDatasourceFunc(f.cliConfig, alias, f.gomplateData, args...)
}

func (f AtmosFuncs) Vals(ref string) (interface{}, error) {
	return valsFunc(f.cliConfig, ref)
}

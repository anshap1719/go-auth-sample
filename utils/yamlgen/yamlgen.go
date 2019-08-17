package yamlgen

import (
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"gigglesearch.org/giggle-auth/utils/yamlgen/dsl"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
)

const apptemp = `
{{- define "yaml" -}}
runtime: go
service: {{ .Service }}
api_version: go1.9

handlers:
- url: /api/{{ .Service }}/swagger
  mime_type: text/plain
  static_files: "../swagger/swagger.yaml"
  upload: "../swagger/swagger.yaml"
  http_headers:
    Access-Control-Allow-Origin: "*"
{{- range .Endpoints }}
- url: {{ .Path }}
  secure: always
  redirect_http_response_code: 307
  {{- if .IsAdmin}}
  login: admin
  auth_fail_action: unauthorized
  {{- end}}
  script: _go_app
{{- end }}
{{ end }}
`

type urlPath struct {
	Path    string
	IsAdmin bool
}

type yaml struct {
	Service   string
	Endpoints []urlPath
}

func Generate() ([]string, error) {
	var (
		ver    string
		outDir string
	)
	set := flag.NewFlagSet("app", flag.PanicOnError)
	set.String("design", "", "") // Consume design argument so Parse doesn't complain
	set.StringVar(&ver, "version", "", "")
	set.StringVar(&outDir, "out", "", "")
	set.Parse(os.Args[2:])

	// First check compatibility
	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}

	g := Generator{
		OutDir: outDir,
	}
	return g.Generate()
}

type Generator struct {
	OutDir string
}

func (g *Generator) Generate() ([]string, error) {
	return WriteURLs(design.Design, g.OutDir)
}

func WriteURLs(api *design.APIDefinition, outDir string) ([]string, error) {
	endpoints := map[string]bool{}
	appYaml := yaml{
		Service: api.Name,
	}
	api.IterateResources(func(res *design.ResourceDefinition) error {
		res.IterateActions(func(a *design.ActionDefinition) error {
			for _, v := range a.Routes {
				path := v.FullPath()
				if i := strings.IndexRune(path, '*'); i != -1 {
					path = path[:i] + ".*"
				}
				if i := strings.IndexRune(path, ':'); i != -1 {
					path = path[:i] + ".*"
				}
				if !endpoints[path] {
					endpoints[path] = true
					appYaml.Endpoints = append(appYaml.Endpoints, urlPath{
						Path:    path,
						IsAdmin: dsl.DslRoot.AdminMap[path],
					})
				}
			}
			return nil
		})
		return nil
	})

	t, err := template.New("yaml").Parse(apptemp)
	if err != nil {
		return nil, err
	}
	outputFile := filepath.Join(outDir, "app.yaml")
	f, err := os.Create(outputFile)
	if err != nil {
		return nil, err
	}
	t.ExecuteTemplate(f, "yaml", appYaml)
	return []string{outputFile}, nil
}

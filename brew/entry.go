package brew

import (
	"bytes"
	"strings"
	"text/template"
)

type Entry struct {
	Name                 string
	RequiredDependencies []string
	RequiredBy           []string
	RestartService       string
	Args                 []string
	Info                 Info
}

func (e *Entry) FromInfo(i Info) {
	e.Name = i.FullName
	e.Info = i
	e.DetermineReqDeps()
}

func (e *Entry) DetermineReqDeps() {
	for _, dependency := range e.Info.Dependencies {
		e.RequiredDependencies = append(e.RequiredDependencies, dependency)
	}

	for _, optional := range e.Info.OptionalDependencies {
		e.RequiredDependencies = remove(e.RequiredDependencies, optional)
	}

	for _, build := range e.Info.BuildDependencies {
		e.RequiredDependencies = remove(e.RequiredDependencies, build)
	}

	for _, rec := range e.Info.RecommendedDependencies {
		e.RequiredDependencies = remove(e.RequiredDependencies, rec)
	}
}

func (e *Entry) Format() (string, error) {
	var bytes bytes.Buffer

	source := `brew '{{ .Name }}'

	{{- if .Args -}} , args: ['{{ StringsJoin .Args "', '" }}'] {{- end -}}

	{{- if .RestartService -}} , restart_service: {{ .RestartService }} {{- end -}}

	{{- if .RequiredBy }} # required by: {{ StringsJoin .RequiredBy ", " }} {{- end -}}`

	tmpl := template.Must(template.New("brew").Funcs(template.FuncMap{"StringsJoin": strings.Join}).Parse(source))
	if err := tmpl.Execute(&bytes, e); err != nil {
		return "", err
	}

	return bytes.String(), nil
}

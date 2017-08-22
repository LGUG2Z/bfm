package brew

import (
	"bytes"
	"strings"
	"text/template"
)

type Entry struct {
	Args                    []string
	BuildDependencies       []string
	Info                    Info
	Name                    string
	OptionalDependencies    []string
	RecommendedDependencies []string
	RequiredBy              []string
	RequiredDependencies    []string
	RestartService          string
}

func (e *Entry) FromInfo(i Info) {
	e.Name = i.FullName
	e.Info = i
	e.DetermineDependencies()
}

func (e *Entry) DetermineDependencies() {
	for _, dependency := range e.Info.Dependencies {
		e.RequiredDependencies = append(e.RequiredDependencies, dependency)
	}

	for _, optional := range e.Info.OptionalDependencies {
		e.RequiredDependencies = remove(e.RequiredDependencies, optional)
		e.OptionalDependencies = append(e.OptionalDependencies, optional)
	}

	for _, build := range e.Info.BuildDependencies {
		e.RequiredDependencies = remove(e.RequiredDependencies, build)
		e.BuildDependencies = append(e.BuildDependencies, build)
	}

	for _, recommended := range e.Info.RecommendedDependencies {
		e.RequiredDependencies = remove(e.RequiredDependencies, recommended)
		e.RecommendedDependencies = append(e.RecommendedDependencies, recommended)
	}
}

func (e *Entry) Format() (string, error) {
	var bytes bytes.Buffer

	source := `brew '{{ .Name }}'

	{{- if .Args -}} , args: ['{{ StringsJoin .Args "', '" }}'] {{- end -}}

	{{- if .RestartService -}} , restart_service: {{ .RestartService }} {{- end -}}

	{{- if .RequiredBy }} # required by: {{ StringsJoin .RequiredBy ", " }} {{- end -}}`

	funcMap := template.FuncMap{"StringsJoin": strings.Join}
	tmpl := template.Must(template.New("brew").Funcs(funcMap).Parse(source))
	if err := tmpl.Execute(&bytes, e); err != nil {
		return "", err
	}

	return bytes.String(), nil
}

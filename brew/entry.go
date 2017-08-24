package brew

import (
	"bytes"
	"strings"
	"text/template"
)

type Entry struct {
	Args                    []string
	BuildDependencies       []string
	Name                    string
	OptionalDependencies    []string
	RecommendedDependencies []string
	RequiredBy              []string
	RequiredDependencies    []string
	RestartService          string
}

func (e *Entry) FromInfo(i Info) {
	e.Name = i.FullName
	e.DetermineDependencies(i)
}

func (e *Entry) DetermineDependencies(i Info) {
	for _, dependency := range i.Dependencies {
		e.RequiredDependencies = append(e.RequiredDependencies, dependency)
	}

	for _, optional := range i.OptionalDependencies {
		e.RequiredDependencies = remove(e.RequiredDependencies, optional)
		e.OptionalDependencies = append(e.OptionalDependencies, optional)
	}

	for _, build := range i.BuildDependencies {
		e.RequiredDependencies = remove(e.RequiredDependencies, build)
		e.BuildDependencies = append(e.BuildDependencies, build)
	}

	for _, recommended := range i.RecommendedDependencies {
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

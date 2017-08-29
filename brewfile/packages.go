package brewfile

import (
	"bytes"
	"io/ioutil"
	"sort"
	"strings"
	"text/template"
)

type Packages struct {
	Tap, Brew, Cask, Mas []string
}

// Parses a Brewfile and separates the taps, brews, casks and mas apps.
func (p *Packages) FromBrewfile(brewfilePath string) error {
	bytes, err := ioutil.ReadFile(brewfilePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(bytes), "\n")

	p.Tap = separate("tap", lines)
	p.Brew = separate("brew", lines)
	p.Cask = separate("cask", lines)
	p.Mas = separate("mas", lines)

	sort.Strings(p.Tap)
	sort.Strings(p.Brew)
	sort.Strings(p.Cask)
	sort.Strings(p.Mas)

	return nil
}

// Creates the final output of an updated Brewfile as a byte array in the order taps ->
// primary brews -> dependent brews -> casks -> mas apps.
func (p *Packages) Bytes() ([]byte, error) {
	entries := `{{ range . }}
{{- . }}
{{ end -}}`

	var primaryBrews []string
	var dependentBrews []string
	for _, b := range p.Brew {
		if strings.Contains(b, "#") {
			dependentBrews = append(dependentBrews, b)
		} else {
			primaryBrews = append(primaryBrews, b)
		}
	}

	var tapBuffer, primaryBuffer, dependentBuffer, caskBuffer, masBuffer bytes.Buffer

	tmpl := template.Must(template.New("entries").Parse(entries))
	if err := tmpl.Execute(&tapBuffer, p.Tap); err != nil {
		return []byte{}, err
	}

	if err := tmpl.Execute(&primaryBuffer, primaryBrews); err != nil {
		return []byte{}, err
	}

	if err := tmpl.Execute(&dependentBuffer, dependentBrews); err != nil {
		return []byte{}, err
	}

	if err := tmpl.Execute(&caskBuffer, p.Cask); err != nil {
		return []byte{}, err
	}

	if err := tmpl.Execute(&masBuffer, p.Mas); err != nil {
		return []byte{}, err
	}

	var lines []string
	if len(tapBuffer.String()) > 0 {
		lines = append(lines, tapBuffer.String())
	}

	if len(primaryBuffer.String()) > 0 {
		lines = append(lines, primaryBuffer.String())
	}
	if len(dependentBuffer.String()) > 0 {
		lines = append(lines, dependentBuffer.String())
	}

	if len(caskBuffer.String()) > 0 {
		lines = append(lines, caskBuffer.String())
	}

	if len(masBuffer.String()) > 0 {
		lines = append(lines, masBuffer.String())
	}

	return []byte(strings.Join(lines, "\n")), nil
}

func separate(packageType string, lines []string) []string {
	var packages []string
	for _, line := range lines {
		if strings.HasPrefix(line, packageType) {
			packages = append(packages, line)
		}
	}

	return packages
}

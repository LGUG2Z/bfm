package brewfile

import (
	"io/ioutil"
	"sort"
	"strings"
)

type Packages struct {
	Tap, Brew, Cask, Mas []string
}

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

func (p *Packages) Bytes() []byte {
	lines := []string{}

	for _, line := range p.Tap {
		lines = append(lines, line)
	}

	lines = append(lines, "")

	for _, line := range p.Brew {
		lines = append(lines, line)
	}

	lines = append(lines, "")

	for _, line := range p.Cask {
		lines = append(lines, line)
	}

	lines = append(lines, "")

	for _, line := range p.Mas {
		lines = append(lines, line)
	}

	return []byte(strings.Join(lines, "\n"))
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

package brew

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
)

type CacheMap struct {
	m Map
	i *InfoCache
}

type Map map[string]Entry

func (m Map) FromBrewfile(entries []string, i *InfoCache) error {
	quotesRegexp := regexp.MustCompile(`'\S+'`)
	argsRegexp := regexp.MustCompile(`\[.*\]`)
	restartRegexp := regexp.MustCompile(`restart_service: (:changed|true)`)
	restartBehaviourRegexp := regexp.MustCompile(`(:changed|true)`)

	for _, e := range entries {
		match := quotesRegexp.FindString(e)
		pkg := match[1 : len(match)-1]

		info, err := i.Find(pkg)
		if err != nil {
			return err
		}

		b := Entry{}
		b.FromInfo(info)

		args := argsRegexp.FindString(e)
		if len(args) > 0 {
			matches := quotesRegexp.FindAllString(args, -1)
			for _, m := range matches {
				arg := m[1 : len(m)-1]
				b.Args = append(b.Args, arg)
			}
		}

		restartService := restartRegexp.FindString(e)
		if len(restartService) > 0 {
			b.RestartService = restartBehaviourRegexp.FindString(restartService)
		}

		m[info.Name] = b
	}

	return nil
}

func (m Map) ResolveDependencies(i *InfoCache) {
	for _, b := range m {
		if len(b.RequiredDependencies) > 0 {
			for _, d := range b.RequiredDependencies {
				m.AddRequiredBy(d, b.Name, i)
			}
		}
	}
}

func (m Map) RemoveRequiredBy(req string, by string, infoCache *InfoCache) {
	b := m[req]

	if contains(b.RequiredBy, by) {
		b.RequiredBy = remove(b.RequiredBy, by)
	}

	sort.Strings(b.RequiredBy)

	m[b.Name] = b
}

func (m Map) AddRequiredBy(req, by string, i *InfoCache) error {
	var b Entry

	if _, present := m[req]; !present {
		info, err := i.Find(req)
		if err != nil {
			return err
		}

		b = Entry{}
		b.FromInfo(info)
	} else {
		b = m[req]
	}

	if !contains(b.RequiredBy, by) {
		b.RequiredBy = append(b.RequiredBy, by)
	}

	sort.Strings(b.RequiredBy)

	m[b.Name] = b

	if len(b.RequiredDependencies) > 0 {
		for _, d := range b.RequiredDependencies {
			m.AddRequiredBy(d, b.Name, i)
		}
	}

	return nil
}

func (m Map) Remove(name string, i *InfoCache) error {
	if _, present := m[name]; !present {
		return errors.New("Nothing to remove.")
	}

	b := m[name]

	if len(b.RequiredDependencies) > 0 {
		for _, dep := range b.RequiredDependencies {
			m.RemoveRequiredBy(dep, name, i)
		}
		fmt.Printf("Removed %s from Brewfile and updated status of its required dependencies: %v \n", b.Name, b.RequiredDependencies)

		for _, dep := range b.RequiredDependencies {
			if len(m[dep].RequiredBy) < 1 {
				fmt.Printf("%s is not required by any other packages. It can be removed if desired.\n", dep)
			}
		}
	} else {
		fmt.Printf("Removed %s from Brewfile.\n", b.Name)
	}

	delete(m, name)

	return nil
}

func (m Map) Add(name, restart string, args []string, i *InfoCache) error {
	info, err := i.Find(name)
	if err != nil {
		return err
	}

	b := Entry{}
	b.FromInfo(info)

	if hasArgs(args) {
		for _, arg := range args {
			b.Args = append(b.Args, arg)
		}
	}

	if hasRestartService(restart) {
		if restart == "changed" {
			b.RestartService = ":changed"
		} else if restart == "always" {
			b.RestartService = "true"
		} else {
			errors.New("Valid options for --restart-services are 'always' and 'changed'. User input ignored.")
		}
	}

	m[info.Name] = b

	if len(b.RequiredDependencies) > 0 {
		for _, dep := range b.RequiredDependencies {
			m.AddRequiredBy(dep, name, i)
		}
		fmt.Printf("Added %s to Brewfile with required dependencies: %v \n", b.Name, b.RequiredDependencies)
	} else {
		fmt.Printf("Added %s to Brewfile.\n", b.Name)
	}

	return nil
}

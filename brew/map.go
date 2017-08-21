package brew

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
)

type Map map[string]Entry

type CacheMap struct {
	M Map
	C *InfoCache
}

func (c CacheMap) FromBrewfile(entries []string) error {
	quotesRegexp := regexp.MustCompile(`'\S+'`)
	argsRegexp := regexp.MustCompile(`\[.*\]`)
	restartRegexp := regexp.MustCompile(`restart_service: (:changed|true)`)
	restartBehaviourRegexp := regexp.MustCompile(`(:changed|true)`)

	for _, e := range entries {
		match := quotesRegexp.FindString(e)
		pkg := match[1 : len(match)-1]

		info, err := c.C.Find(pkg)
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

		c.M[info.Name] = b
	}

	return nil
}

func (c CacheMap) ResolveDependencies() {
	for _, b := range c.M {
		if len(b.RequiredDependencies) > 0 {
			for _, d := range b.RequiredDependencies {
				c.AddRequiredBy(d, b.Name)
			}
		}
	}
}

func (c CacheMap) RemoveRequiredBy(req, by string) {
	b := c.M[req]

	if contains(b.RequiredBy, by) {
		b.RequiredBy = remove(b.RequiredBy, by)
	}

	sort.Strings(b.RequiredBy)

	c.M[b.Name] = b
}

func (c CacheMap) AddRequiredBy(req, by string) error {
	var b Entry

	if _, present := c.M[req]; !present {
		info, err := c.C.Find(req)
		if err != nil {
			return err
		}

		b = Entry{}
		b.FromInfo(info)
	} else {
		b = c.M[req]
	}

	if !contains(b.RequiredBy, by) {
		b.RequiredBy = append(b.RequiredBy, by)
	}

	sort.Strings(b.RequiredBy)

	c.M[b.Name] = b

	if len(b.RequiredDependencies) > 0 {
		for _, d := range b.RequiredDependencies {
			c.AddRequiredBy(d, b.Name)
		}
	}

	return nil
}

func (c CacheMap) Remove(name string) error {
	if _, present := c.M[name]; !present {
		return errors.New("Nothing to remove.")
	}

	b := c.M[name]

	if len(b.RequiredDependencies) > 0 {
		for _, dep := range b.RequiredDependencies {
			c.RemoveRequiredBy(dep, name)
		}
		fmt.Printf("Removed %s from Brewfile and updated status of its required dependencies: %v \n", b.Name, b.RequiredDependencies)

		for _, dep := range b.RequiredDependencies {
			if len(c.M[dep].RequiredBy) < 1 {
				fmt.Printf("%s is not required by any other packages. It can be removed if desired.\n", dep)
			}
		}
	} else {
		fmt.Printf("Removed %s from Brewfile.\n", b.Name)
	}

	delete(c.M, name)

	return nil
}

func (c CacheMap) Add(name, restart string, args []string) error {
	info, err := c.C.Find(name)
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

	c.M[info.Name] = b

	if len(b.RequiredDependencies) > 0 {
		for _, dep := range b.RequiredDependencies {
			c.AddRequiredBy(dep, name)
		}
		fmt.Printf("Added %s to Brewfile with required dependencies: %v \n", b.Name, b.RequiredDependencies)
	} else {
		fmt.Printf("Added %s to Brewfile.\n", b.Name)
	}

	return nil
}

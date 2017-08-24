package brew

import (
	"errors"
	"regexp"
	"sort"
)

type Map map[string]Entry

type CacheMap struct {
	Map   Map
	Cache *Cache
}

func (c CacheMap) FromPackages(packages []string) error {
	quotesRegexp := regexp.MustCompile(`'\S+'`)
	argsRegexp := regexp.MustCompile(`\[.*\]`)
	restartRegexp := regexp.MustCompile(`restart_service: (:changed|true)`)
	restartBehaviourRegexp := regexp.MustCompile(`(:changed|true)`)

	for _, p := range packages {
		match := quotesRegexp.FindString(p)
		pkg := match[1 : len(match)-1]

		info, err := c.Cache.Find(pkg, c.Cache.DB)
		if err != nil {
			return err
		}

		e := Entry{}
		e.FromInfo(info)

		args := argsRegexp.FindString(p)
		if len(args) > 0 {
			matches := quotesRegexp.FindAllString(args, -1)
			for _, m := range matches {
				arg := m[1 : len(m)-1]
				e.Args = append(e.Args, arg)
			}
		}

		restartService := restartRegexp.FindString(p)
		if len(restartService) > 0 {
			e.RestartService = restartBehaviourRegexp.FindString(restartService)
		}

		c.Map[info.FullName] = e
	}

	return nil
}

func (c CacheMap) ResolveRequiredDependencyMap() error {
	for _, b := range c.Map {
		if len(b.RequiredDependencies) > 0 {
			for _, d := range b.RequiredDependencies {
				if err := c.addRequiredBy(d, b.Name); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c CacheMap) Add(entry Entry, opt int) error {
	info, err := c.Cache.Find(entry.Name, c.Cache.DB)
	if err != nil {
		return err
	}

	entry.FromInfo(info)

	switch opt {
	case AddAll:
		c.Map[entry.Name] = entry

		for _, dep := range entry.RequiredDependencies {
			c.addRequiredBy(dep, entry.Name)
		}
		for _, dep := range entry.RecommendedDependencies {
			c.Add(Entry{Name: dep}, opt)
		}
		for _, dep := range entry.OptionalDependencies {
			c.Add(Entry{Name: dep}, opt)
		}
		for _, dep := range entry.BuildDependencies {
			c.Add(Entry{Name: dep}, opt)
		}
	case AddPackageOnly:
		c.Map[entry.Name] = entry
	case AddPackageAndRequired:
		c.Map[entry.Name] = entry

		for _, dep := range entry.RequiredDependencies {
			c.addRequiredBy(dep, entry.Name)
		}
	}

	return nil
}

func (c CacheMap) Remove(name string, opt int) error {
	if _, present := c.Map[name]; !present {
		return errors.New("Nothing to remove.")
	}

	b := c.Map[name]

	switch opt {
	case RemoveAll:
		for _, dep := range b.RequiredDependencies {
			c.removeRequiredBy(dep, name)
		}

		for _, dep := range b.RequiredDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				c.Remove(c.Map[dep].Name, opt)
			}
		}

		for _, dep := range b.RecommendedDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				c.Remove(c.Map[dep].Name, opt)
			}
		}

		for _, dep := range b.OptionalDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				c.Remove(c.Map[dep].Name, opt)
			}
		}

		for _, dep := range b.BuildDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				c.Remove(c.Map[dep].Name, opt)
			}
		}
	case RemovePackageOnly:
		for _, dep := range b.RequiredDependencies {
			c.removeRequiredBy(dep, name)
		}
	case RemovePackageAndRequired:
		for _, dep := range b.RequiredDependencies {
			c.removeRequiredBy(dep, name)
		}

		for _, dep := range b.RequiredDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				c.Remove(c.Map[dep].Name, opt)
			}
		}
	}

	delete(c.Map, name)
	return nil
}

func (c CacheMap) addRequiredBy(req, by string) error {
	var e Entry

	if _, present := c.Map[req]; !present {
		info, err := c.Cache.Find(req, c.Cache.DB)
		if err != nil {
			return err
		}

		e = Entry{}
		e.FromInfo(info)
	} else {
		e = c.Map[req]
	}

	if !contains(e.RequiredBy, by) {
		e.RequiredBy = append(e.RequiredBy, by)
	}

	sort.Strings(e.RequiredBy)

	c.Map[e.Name] = e

	if len(e.RequiredDependencies) > 0 {
		for _, d := range e.RequiredDependencies {
			c.addRequiredBy(d, e.Name)
		}
	}

	return nil
}

func (c CacheMap) removeRequiredBy(req, by string) {
	b := c.Map[req]

	if contains(b.RequiredBy, by) {
		b.RequiredBy = remove(b.RequiredBy, by)
	}

	sort.Strings(b.RequiredBy)

	c.Map[b.Name] = b
}

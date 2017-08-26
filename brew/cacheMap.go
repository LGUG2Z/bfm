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

		info, err := c.Cache.Find(pkg)
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

// TODO: This needs to take recommended/all for clean
func (c CacheMap) ResolveDependencyMap(level int) error {
	if level >= Required {
		for _, b := range c.Map {
			if len(b.RequiredDependencies) > 0 {
				for _, d := range b.RequiredDependencies {
					if err := c.addDependency(d, b.Name, RequiredDependency); err != nil {
						return err
					}
				}
			}
		}
	}

	if level >= Recommended {
		for _, b := range c.Map {
			if len(b.RecommendedDependencies) > 0 {
				for _, d := range b.RecommendedDependencies {
					if err := c.addDependency(d, b.Name, RecommendedDependency); err != nil {
						return err
					}
				}
			}
		}
	}

	if level >= Optional {
		for _, b := range c.Map {
			if len(b.OptionalDependencies) > 0 {
				for _, d := range b.OptionalDependencies {
					if err := c.addDependency(d, b.Name, OptionalDependency); err != nil {
						return err
					}
				}
			}
		}
	}

	if level >= Build {
		for _, b := range c.Map {
			if len(b.BuildDependencies) > 0 {
				for _, d := range b.BuildDependencies {
					if err := c.addDependency(d, b.Name, BuildDependency); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (c CacheMap) Add(entry Entry, level int) error {
	info, err := c.Cache.Find(entry.Name)
	if err != nil {
		return err
	}

	entry.FromInfo(info)
	c.Map[entry.Name] = entry

	if level >= Required {
		for _, dep := range entry.RequiredDependencies {
			if err := c.addDependency(dep, entry.Name, RequiredDependency); err != nil {
				return err
			}
		}
	}
	if level >= Recommended {
		for _, dep := range entry.RecommendedDependencies {
			if err := c.addDependency(dep, entry.Name, RecommendedDependency); err != nil {
				return err
			}
		}
	}
	if level >= Optional {
		for _, dep := range entry.OptionalDependencies {
			if err := c.addDependency(dep, entry.Name, OptionalDependency); err != nil {
				return err
			}
		}
	}
	if level >= Build {
		for _, dep := range entry.BuildDependencies {
			if err := c.addDependency(dep, entry.Name, BuildDependency); err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO: Think about how to rework this
func (c CacheMap) Remove(name string, level int) error {
	if _, present := c.Map[name]; !present {
		return errors.New("Nothing to remove.")
	}

	entry := c.Map[name]

	if level >= Required {
		for _, dep := range entry.RequiredDependencies {
			c.removeDependency(dep, name, RequiredDependency)
		}

		for _, dep := range entry.RequiredDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				if err := c.Remove(c.Map[dep].Name, level); err != nil {
					return err
				}
			}
		}
	}
	if level >= Recommended {
		for _, dep := range entry.RecommendedDependencies {
			c.removeDependency(dep, name, RecommendedDependency)
		}

		for _, dep := range entry.RecommendedDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				if err := c.Remove(c.Map[dep].Name, level); err != nil {
					return err
				}
			}
		}
	}
	if level >= Optional {
		for _, dep := range entry.OptionalDependencies {
			c.removeDependency(dep, name, OptionalDependency)
		}

		for _, dep := range entry.OptionalDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				if err := c.Remove(c.Map[dep].Name, level); err != nil {
					return err
				}
			}
		}
	}
	if level >= Build {
		for _, dep := range entry.BuildDependencies {
			c.removeDependency(dep, name, BuildDependency)
		}

		for _, dep := range entry.BuildDependencies {
			if len(c.Map[dep].RequiredBy) < 1 {
				if err := c.Remove(c.Map[dep].Name, level); err != nil {
					return err
				}
			}
		}
	}

	delete(c.Map, name)
	return nil
}

func (c CacheMap) addDependency(req, by string, dependencyType int) error {
	var e Entry

	if _, present := c.Map[req]; !present {
		info, err := c.Cache.Find(req)
		if err != nil {
			return err
		}

		e = Entry{}
		e.FromInfo(info)
	} else {
		e = c.Map[req]
	}

	switch dependencyType {
	case RequiredDependency:
		if !contains(e.RequiredBy, by) {
			e.RequiredBy = append(e.RequiredBy, by)
			sort.Strings(e.RequiredBy)
		}

		c.Map[e.Name] = e

		if len(e.RequiredDependencies) > 0 {
			for _, d := range e.RequiredDependencies {
				if err := c.addDependency(d, e.Name, RequiredDependency); err != nil {
					return err
				}
			}
		}
	case RecommendedDependency:
		if !contains(e.RecommendedFor, by) {
			e.RecommendedFor = append(e.RecommendedFor, by)
			sort.Strings(e.RecommendedFor)
		}

		c.Map[e.Name] = e

		if len(e.RecommendedDependencies) > 0 {
			for _, d := range e.RecommendedDependencies {
				if err := c.addDependency(d, e.Name, RecommendedDependency); err != nil {
					return err
				}
			}
		}
	case OptionalDependency:
		if !contains(e.OptionalFor, by) {
			e.OptionalFor = append(e.OptionalFor, by)
			sort.Strings(e.OptionalFor)
		}

		c.Map[e.Name] = e

		if len(e.OptionalDependencies) > 0 {
			for _, d := range e.OptionalDependencies {
				if err := c.addDependency(d, e.Name, OptionalDependency); err != nil {
					return err
				}
			}
		}
	case BuildDependency:
		if !contains(e.BuildOf, by) {
			e.BuildOf = append(e.BuildOf, by)
			sort.Strings(e.BuildOf)
		}

		c.Map[e.Name] = e

		if len(e.BuildDependencies) > 0 {
			for _, d := range e.BuildDependencies {
				if err := c.addDependency(d, e.Name, BuildDependency); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c CacheMap) removeDependency(req, by string, dependencyType int) {
	b := c.Map[req]

	switch dependencyType {
	case RequiredDependency:
		if contains(b.RequiredBy, by) {
			b.RequiredBy = remove(b.RequiredBy, by)
			sort.Strings(b.RequiredBy)
		}
	case RecommendedDependency:
		if contains(b.RecommendedFor, by) {
			b.RecommendedFor = remove(b.RecommendedFor, by)
			sort.Strings(b.RecommendedFor)
		}
	case OptionalDependency:
		if contains(b.OptionalFor, by) {
			b.OptionalFor = remove(b.OptionalFor, by)
			sort.Strings(b.OptionalFor)
		}
	case BuildDependency:
		if contains(b.BuildOf, by) {
			b.BuildOf = remove(b.BuildOf, by)
			sort.Strings(b.BuildOf)
		}
	}

	c.Map[b.Name] = b
}

package brew

import (
	"errors"
	"regexp"
	"sort"

	"github.com/boltdb/bolt"
)

type Map map[string]Entry

type CacheMap struct {
	Map   Map
	Cache *InfoCache
}

func (c CacheMap) FromPackages(entries []string, db *bolt.DB) error {
	quotesRegexp := regexp.MustCompile(`'\S+'`)
	argsRegexp := regexp.MustCompile(`\[.*\]`)
	restartRegexp := regexp.MustCompile(`restart_service: (:changed|true)`)
	restartBehaviourRegexp := regexp.MustCompile(`(:changed|true)`)

	for _, e := range entries {
		match := quotesRegexp.FindString(e)
		pkg := match[1 : len(match)-1]

		info, err := c.Cache.Find(pkg, db)
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

		c.Map[info.FullName] = b
	}

	return nil
}

func (c CacheMap) ResolveRequiredDependencyMap(db *bolt.DB) {
	for _, b := range c.Map {
		if len(b.RequiredDependencies) > 0 {
			for _, d := range b.RequiredDependencies {
				c.addRequiredBy(d, b.Name, db)
			}
		}
	}
}

func (c CacheMap) Add(entry Entry, opt int, db *bolt.DB) error {
	info, err := c.Cache.Find(entry.Name, db)
	if err != nil {
		return err
	}

	entry.FromInfo(info)

	switch opt {
	case AddAll:
		c.Map[entry.Name] = entry

		for _, dep := range entry.RequiredDependencies {
			c.addRequiredBy(dep, entry.Name, db)
		}
		for _, dep := range entry.RecommendedDependencies {
			c.Add(Entry{Name: dep}, opt, db)
		}
		for _, dep := range entry.OptionalDependencies {
			c.Add(Entry{Name: dep}, opt, db)
		}
		for _, dep := range entry.BuildDependencies {
			c.Add(Entry{Name: dep}, opt, db)
		}
	case AddPackageOnly:
		c.Map[entry.Name] = entry
	case AddPackageAndRequired:
		c.Map[entry.Name] = entry

		for _, dep := range entry.RequiredDependencies {
			c.addRequiredBy(dep, entry.Name, db)
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

func (c CacheMap) addRequiredBy(req, by string, db *bolt.DB) error {
	var e Entry

	if _, present := c.Map[req]; !present {
		info, err := c.Cache.Find(req, db)
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
			c.addRequiredBy(d, e.Name, db)
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

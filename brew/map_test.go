package brew

import (
	"reflect"
	"testing"
)

func TestCacheMap_FromBrewfile(t *testing.T) {
	expected := []string{"vim", "emacs"}

	info := []Info{}

	for _, e := range expected {
		info = append(info, Info{Name: e})
	}

	i := InfoCache(info)

	brewfile := []string{
		"brew 'vim'",
		"brew 'emacs'",
	}

	c := CacheMap{
		C: &i,
		M: make(Map),
	}

	c.FromBrewfile(brewfile)

	for _, e := range expected {
		if _, present := c.M[e]; !present {
			t.Fatalf("Expected %s to be present", e)
		}
	}
}

func TestCacheMap_ResolveDependencies(t *testing.T) {
	expected := []string{"vim"}

	info := []Info{
		Info{
			Name:         "vim",
			FullName:     "vim",
			Dependencies: []string{"python"},
		},
		Info{
			Name:         "python",
			FullName:     "python",
			Dependencies: []string{"openssl"},
		},
	}

	brewfile := []string{"brew 'vim'", "brew 'python'"}

	i := InfoCache(info)

	c := CacheMap{
		C: &i,
		M: make(Map),
	}

	c.FromBrewfile(brewfile)
	c.ResolveDependencies()

	actual := c.M["python"].RequiredBy

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestCacheMap_AddRequiredBy(t *testing.T) {
	info := []Info{
		Info{
			Name:         "vim",
			FullName:     "vim",
			Dependencies: []string{"python"},
		},
		Info{
			Name:         "python",
			FullName:     "python",
			Dependencies: []string{"openssl"},
		},
	}

	brewfile := []string{
		"brew 'vim'",
		"brew 'python'",
	}

	i := InfoCache(info)

	c := CacheMap{
		C: &i,
		M: make(Map),
	}

	c.FromBrewfile(brewfile)
	c.ResolveDependencies()

	new := Info{
		Name:         "neovim",
		FullName:     "neovim",
		Dependencies: []string{"python"},
	}

	c.AddRequiredBy(new.Dependencies[0], new.Name)

	expected := []string{"neovim", "vim"}
	actual := c.M["python"].RequiredBy

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestCacheMap_RemoveRequiredBy(t *testing.T) {
	info := []Info{
		Info{
			Name:         "vim",
			FullName:     "vim",
			Dependencies: []string{"python"},
		},
		Info{
			Name:     "python",
			FullName: "python",
		},
	}

	brewfile := []string{
		"brew 'vim'",
		"brew 'python'",
	}

	i := InfoCache(info)

	c := CacheMap{
		C: &i,
		M: make(Map),
	}

	c.FromBrewfile(brewfile)
	c.ResolveDependencies()

	c.RemoveRequiredBy("python", "vim")

	expected := []string{}
	actual := c.M["python"].RequiredBy

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestCacheMap_Add(t *testing.T) {
	info := []Info{
		Info{
			Name:         "vim",
			FullName:     "vim",
			Dependencies: []string{"python"},
		},
	}

	i := InfoCache(info)

	actual := CacheMap{
		C: &i,
		M: make(Map),
	}

	actual.Add("vim", "always", []string{"with-override-system-vim"})

	expected := CacheMap{
		C: &i,
		M: Map{
			"vim": Entry{
				Name:                 "vim",
				RequiredDependencies: []string{"python"},
				RestartService:       "true",
				Args:                 []string{"with-override-system-vim"},
				Info:                 info[0],
			},
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestCacheMap_Remove(t *testing.T) {
	info := []Info{
		Info{
			Name:         "vim",
			FullName:     "vim",
			Dependencies: []string{"python"},
		},
	}

	i := InfoCache(info)
	c := CacheMap{
		C: &i,
		M: make(Map),
	}

	c.Add("vim", "always", []string{"with-override-system-vim"})
	c.Remove("vim")

	if _, present := c.M["vim"]; present {
		t.Fatalf("Expected vim to not be present.")
	}
}

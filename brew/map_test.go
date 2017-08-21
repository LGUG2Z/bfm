package brew

import (
	"reflect"
	"testing"
)

func TestMap_FromBrewfile(t *testing.T) {
	expected := []string{"vim", "emacs"}

	infos := []Info{}

	for _, e := range expected {
		infos = append(infos, Info{Name: e})
	}

	i := InfoCache(infos)

	e := []string{
		"brew 'vim'",
		"brew 'emacs'",
	}

	actual := make(Map)

	actual.FromBrewfile(e, &i)

	for _, e := range expected {
		if _, present := actual[e]; !present {
			t.Fatalf("Expected %s to be present", e)
		}
	}
}

func TestMap_ResolveDependencies(t *testing.T) {
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

	actual := make(Map)
	actual.FromBrewfile(brewfile, &i)
	actual.ResolveDependencies(&i)

	expected := []string{"vim"}

	python := actual["python"].RequiredBy

	if !reflect.DeepEqual(python, expected) {
		t.Fatalf("Expected %t but got %t", expected, python)
	}
}

func TestMap_AddRequiredBy(t *testing.T) {
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

	actual := make(Map)
	actual.FromBrewfile(brewfile, &i)
	actual.ResolveDependencies(&i)

	newInfo := Info{
		Name:         "neovim",
		FullName:     "neovim",
		Dependencies: []string{"python"},
	}

	actual.AddRequiredBy(newInfo.Dependencies[0], newInfo.Name, &i)

	python := actual["python"].RequiredBy
	expected := []string{"neovim", "vim"}

	if !reflect.DeepEqual(python, expected) {
		t.Fatalf("Expected %t but got %t", expected, python)
	}
}

func TestMap_RemoveRequiredBy(t *testing.T) {
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

	brewfile := []string{"brew 'vim'", "brew 'python'"}

	i := InfoCache(info)

	actual := make(Map)
	actual.FromBrewfile(brewfile, &i)
	actual.ResolveDependencies(&i)

	actual.RemoveRequiredBy("python", "vim", &i)

	python := actual["python"].RequiredBy
	expected := []string{}

	if !reflect.DeepEqual(python, expected) {
		t.Fatalf("Expected %t but got %t", expected, python)
	}

}

func TestMap_Add(t *testing.T) {
	info := []Info{
		Info{
			Name:         "vim",
			FullName:     "vim",
			Dependencies: []string{"python"},
		},
	}

	i := InfoCache(info)

	actual := make(Map)
	actual.Add("vim", "always", []string{"with-override-system-vim"}, &i)

	expected := Map{
		"vim": Entry{
			Name:                 "vim",
			RequiredDependencies: []string{"python"},
			RestartService:       "true",
			Args:                 []string{"with-override-system-vim"},
			Info:                 info[0],
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestMap_Remove(t *testing.T) {
	info := []Info{
		Info{
			Name:         "vim",
			FullName:     "vim",
			Dependencies: []string{"python"},
		},
	}

	i := InfoCache(info)

	actual := make(Map)
	actual.Add("vim", "always", []string{"with-override-system-vim"}, &i)
	actual.Remove("vim", &i)

	if _, present := actual["vim"]; present {
		t.Fatalf("Expected vim to not be present.")
	}
}

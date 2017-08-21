package brew

import (
	"reflect"
	"testing"
)

func TestEntry_FromInfo(t *testing.T) {
	i := Info{
		FullName:                "a",
		Dependencies:            []string{"b", "c", "d", "e", "f"},
		OptionalDependencies:    []string{"c"},
		BuildDependencies:       []string{"d"},
		RecommendedDependencies: []string{"e", "f"},
	}

	expected := Entry{
		Name:                 "a",
		RequiredDependencies: []string{"b"},
		Info: Info{
			FullName:                "a",
			Dependencies:            []string{"b", "c", "d", "e", "f"},
			OptionalDependencies:    []string{"c"},
			BuildDependencies:       []string{"d"},
			RecommendedDependencies: []string{"e", "f"},
		},
	}

	actual := Entry{}
	actual.FromInfo(i)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestEntry_DetermineReqDeps(t *testing.T) {
	i := Info{
		Dependencies: []string{"b"},
	}

	e := Entry{
		Info: i,
	}

	e.DetermineReqDeps()

	actual := e.RequiredDependencies
	expected := []string{"b"}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, e)
	}
}

func TestEntry_DetermineReqDepsMixed(t *testing.T) {
	i := Info{
		Dependencies:            []string{"b", "c", "d", "e", "f"},
		OptionalDependencies:    []string{"c"},
		BuildDependencies:       []string{"d"},
		RecommendedDependencies: []string{"e", "f"},
	}

	e := Entry{
		Info: i,
	}

	e.DetermineReqDeps()

	actual := e.RequiredDependencies
	expected := []string{"b"}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, e)
	}
}

func TestEntry_Format(t *testing.T) {
	e := Entry{
		Name: "a",
	}

	expected := `brew 'a'`
	actual, _ := e.Format()

	if expected != actual {
		t.Fatalf("Expected %t but got %t", expected, e)
	}
}

func TestEntry_FormatWithArgsRestartServiceAndRequiredBy(t *testing.T) {
	e := Entry{
		Name:           "a",
		Args:           []string{"with-tests"},
		RestartService: ":changed",
		RequiredBy:     []string{"z", "y"},
	}

	expected := `brew 'a', args: ['with-tests'], restart_service: :changed # required by: z, y`
	actual, _ := e.Format()

	if expected != actual {
		t.Fatalf("Expected %t but got %t", expected, e)
	}
}

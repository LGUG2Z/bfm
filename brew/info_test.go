package brew

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestInfoCache_Find(t *testing.T) {
	info := Info{
		FullName: "a",
	}

	infos := []Info{info}

	i := InfoCache(infos)

	actual, _ := i.Find("a")
	expected := info

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestInfoCache_FindError(t *testing.T) {
	i := InfoCache([]Info{})

	_, err := i.Find("a")

	if err == nil {
		t.Fatal("Expected an error but error was nil")
	}
}

func TestInfoCache_Read(t *testing.T) {
	info := Info{
		Name:     "a2ps",
		FullName: "a2ps",
		Desc:     "Any-to-PostScript filter",
	}

	expected := InfoCache([]Info{info})

	actual := InfoCache([]Info{})
	goPath := os.Getenv("GOPATH")
	actual.Read(fmt.Sprintf("%s/src/github.com/lgug2z/bfm/testData/test.json", goPath))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}

}

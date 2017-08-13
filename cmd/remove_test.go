package cmd

import (
	"reflect"
	"testing"
)

func TestRemovePackage(t *testing.T) {
	packageType := "brew"
	packageToRemove := "bfm"
	packages := []string{"brew 'vim'", "brew 'bfm'"}

	actual := removePackage(packageType, packageToRemove, packages)
	expected := []string{"brew 'vim'"}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

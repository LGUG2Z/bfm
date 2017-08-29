package cmd

import (
	"reflect"
	"testing"
)

func TestConstructBaseEntry(t *testing.T) {
	packageType := "brew"
	packageName := "bfm"

	actual := constructBaseEntry(packageType, packageName)
	expected := "brew 'bfm'"

	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

func TestEntryExistsTrue(t *testing.T) {
	contents := "brew 'vim'\nbrew 'neovim'"
	packageType := "brew"
	packageName := "neovim"

	actual := entryExists(contents, packageType, packageName)
	expected := true

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestEntryExistsFalse(t *testing.T) {
	contents := "brew 'vim'\nbrew 'neovim'"
	packageType := "cask"
	packageName := "macvim"

	actual := entryExists(contents, packageType, packageName)
	expected := false

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestGetPackages(t *testing.T) {
	lines := []string{"tap 'homebrew/bundle'", "brew 'vim'", "cask 'macvim'"}

	expected := []string{"brew 'vim'"}
	actual := getPackages("brew", lines)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

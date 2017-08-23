package cmd

import (
	"reflect"
	"testing"
)

//func TestFlagProvidedTrue(t *testing.T) {
//	tap := false
//	brew := true
//	cask := false
//	mas := false
//
//	actual := flagProvided(tap, brew, cask, mas)
//	expected := true
//
//	if actual != expected {
//		t.Fatalf("Expected %t but got %t", expected, actual)
//	}
//}
//
//func TestFlagProvidedFalse(t *testing.T) {
//	tap := false
//	brew := false
//	cask := false
//	mas := false
//
//	actual := flagProvided(tap, brew, cask, mas)
//	expected := false
//
//	if actual != expected {
//		t.Fatalf("Expected %t but got %t", expected, actual)
//	}
//}
//
//func TestPackageTypeTap(t *testing.T) {
//	tap := true
//	brew := false
//	cask := false
//	mas := false
//
//	actual := getPackageType(tap, brew, cask, mas)
//	expected := "tap"
//
//	if actual != expected {
//		t.Fatalf("Expected %s but got %s", expected, actual)
//	}
//}
//
//func TestPackageTypeBrew(t *testing.T) {
//	tap := false
//	brew := true
//	cask := false
//	mas := false
//
//	actual := getPackageType(tap, brew, cask, mas)
//	expected := "brew"
//
//	if actual != expected {
//		t.Fatalf("Expected %s but got %s", expected, actual)
//	}
//}
//
//func TestPackageTypeCask(t *testing.T) {
//	tap := false
//	brew := false
//	cask := true
//	mas := false
//
//	actual := getPackageType(tap, brew, cask, mas)
//	expected := "cask"
//
//	if actual != expected {
//		t.Fatalf("Expected %s but got %s", expected, actual)
//	}
//}
//
//func TestPackageTypeMas(t *testing.T) {
//	tap := false
//	brew := false
//	cask := false
//	mas := true
//
//	actual := getPackageType(tap, brew, cask, mas)
//	expected := "mas"
//
//	if actual != expected {
//		t.Fatalf("Expected %s but got %s", expected, actual)
//	}
//}

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

//func TestConstructFileContents(t *testing.T) {
//	tap := []string{"tap 'homebrew/bundle'"}
//	brew := []string{"brew 'vim', args: ['with-python3']"}
//	cask := []string{"cask 'macvim'"}
//	mas := []string{"mas 'Xcode', id: 49779983"}
//
//	expected := `tap 'homebrew/bundle'
//
//brew 'vim', args: ['with-python3']
//
//cask 'macvim'
//
//mas 'Xcode', id: 49779983`
//
//	actual := constructFileContents(tap, brew, cask, mas)
//
//	if actual != expected {
//		t.Fatalf("Expected %s but got %s", expected, actual)
//	}
//}

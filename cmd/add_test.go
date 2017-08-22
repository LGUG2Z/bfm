package cmd

import (
	"reflect"
	"testing"
)

func TestHasCorrectTapFormatTrue(t *testing.T) {
	actual := hasCorrectTapFormat("homebrew/correct")
	expected := true

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestHasCorrectTapFormatFalse(t *testing.T) {
	actual := hasCorrectTapFormat("incorrect")
	expected := false

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestHasMasIdWithArg(t *testing.T) {
	r := "with"
	actual := hasMasID(r)
	expected := true

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestHasMasIdWithNoArg(t *testing.T) {
	r := ""
	actual := hasMasID(r)
	expected := false

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestAppendMasId(t *testing.T) {
	i := "1235644"
	actual := appendMasID("mas 'bfm'", i)
	expected := "mas 'bfm', id: 1235644"

	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

func TestAddPackage(t *testing.T) {
	packageType := "brew"
	newPackage := "bfm"
	packages := []string{"brew 'vim'"}

	actual := addPackage(packageType, newPackage, packages)
	expected := []string{"brew 'vim'", "brew 'bfm'"}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

package brew

import "testing"

func TestHasArgsWithArgs(t *testing.T) {
	args := []string{"one", "two"}
	actual := hasArgs(args)
	expected := true

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestHasArgsWithNoArgs(t *testing.T) {
	a := []string{}
	actual := hasArgs(a)
	expected := false

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestHasRestartServiceWithArg(t *testing.T) {
	r := "with"
	actual := hasRestartService(r)
	expected := true

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

func TestHasRestartServiceWithNoArg(t *testing.T) {
	r := ""
	actual := hasRestartService(r)
	expected := false

	if actual != expected {
		t.Fatalf("Expected %t but got %t", expected, actual)
	}
}

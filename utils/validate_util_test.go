package utils

import (
	"testing"
)

func TestIsUrlValid_URLEmpty(t *testing.T) {
	expected := "Url cannot be empty"

	got := IsUrlValid("")
	if got != expected {
		t.Errorf("Url expected = %q; got = %q", expected, got)
	}
}

func TestIsUrlValid_URLContainWhiteSpace(t *testing.T) {
	expected := "Url cannot be empty"

	got := IsUrlValid("  ")
	if got != expected {
		t.Errorf("Url expected = %q; got = %q", expected, got)
	}
}

func TestIsUrlValid_InvalidUrl(t *testing.T) {
	expected := "Invalid url"

	got := IsUrlValid("htp:/invalid")
	if got != expected {
		t.Errorf("Url expected = %q; got = %q", expected, got)
	}
}

func TestIsUrlValid_MissingScheme(t *testing.T) {
	expected := "Invalid url"

	got := IsUrlValid("www.example.com")
	if got != expected {
		t.Errorf("Url expected = %q; got = %q", expected, got)
	}
}

func TestIsUrlValid_MissingHost(t *testing.T) {
	expected := "Invalid url"

	got := IsUrlValid("http://")
	if got != expected {
		t.Errorf("Url expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion(t *testing.T) {
	doctype := "<!DOCTYPE html>"
	expected := "HTML5"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion_HTML5(t *testing.T) {
	doctype := "<!DOCTYPE html>"
	expected := "HTML5"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion_HTML4(t *testing.T) {
	doctype := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`
	expected := "HTML 4.01 Strict"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}
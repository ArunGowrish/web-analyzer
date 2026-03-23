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

	got := IsUrlValid("www.mock.com")
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

func TestIdentifyHTMLVersion_HTML4Strict(t *testing.T) {
	doctype := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`
	expected := "HTML 4.01 Strict"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion_HTML4Transitional(t *testing.T) {
	doctype := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/transitional.dtd">`
	expected := "HTML 4.01 Transitional"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion_HTML4Frameset(t *testing.T) {
	doctype := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/frameset.dtd">`
	expected := "HTML 4.01 Frameset"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion_XHTMLStrict(t *testing.T) {
	doctype := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`
	expected := "XHTML 1.0 Strict"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion_XHTMLTransitional(t *testing.T) {
	doctype := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`
	expected := "XHTML 1.0 Transitional"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion_XHTMLFrameset(t *testing.T) {
	doctype := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd">`
	expected := "XHTML 1.0 Frameset"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

func TestIdentifyHTMLVersion_XHTML(t *testing.T) {
	doctype := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">`
	expected := "XHTML 1.1"

	got := IdentifyHTMLVersion(doctype)

	if got != expected {
		t.Errorf("HTML version expected = %q; got = %q", expected, got)
	}
}

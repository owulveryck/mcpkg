package mcp

import "testing"

func TestExtractGraphPathFromURI(t *testing.T) {
	uri := "graph://mygraph?from=subject1&to=subject2"
	expected := "mygraph"
	result := extractGraphPathFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractFromFromURI(t *testing.T) {
	uri := "graph://mygraph?from=subject1&to=subject2"
	expected := "subject1"
	result := extractFromFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractToFromURI(t *testing.T) {
	uri := "graph://mygraph?from=subject1&to=subject2"
	expected := "subject2"
	result := extractToFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractGraphPathFromURI_InvalidURI(t *testing.T) {
	uri := "invalid-uri"
	expected := ""
	result := extractGraphPathFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractFromFromURI_NoFrom(t *testing.T) {
	uri := "graph://mygraph?to=subject2"
	expected := ""
	result := extractFromFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractToFromURI_NoTo(t *testing.T) {
	uri := "graph://mygraph?from=subject1"
	expected := ""
	result := extractToFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractGraphPathFromURI_EmptyURI(t *testing.T) {
	uri := ""
	expected := ""
	result := extractGraphPathFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractFromFromURI_EmptyURI(t *testing.T) {
	uri := ""
	expected := ""
	result := extractFromFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractToFromURI_EmptyURI(t *testing.T) {
	uri := ""
	expected := ""
	result := extractToFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractFromFromURI_InvalidURI2(t *testing.T) {
	uri := "graph://mygraph?from=subject1&to=subject2&from=subject3"
	expected := "subject1"
	result := extractFromFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractFromFromURI_ComplexSubject(t *testing.T) {
	uri := "graph://mygraph?from=http://example.com/subject1&to=subject2"
	expected := "http://example.com/subject1"
	result := extractFromFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractToFromURI_ComplexSubject(t *testing.T) {
	uri := "graph://mygraph?from=subject1&to=http://example.com/subject2"
	expected := "http://example.com/subject2"
	result := extractToFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractGraphPathFromURI_ComplexPath(t *testing.T) {
	uri := "graph://my-complex-graph?from=subject1&to=subject2"
	expected := "my-complex-graph"
	result := extractGraphPathFromURI(uri)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

package mcp

import "net/url"

// extractGraphPathFromURI extracts the {knowledge_graph_path} parameter from an URI of type graph://{knowledge_graph_path}?from={from_subject}&to={to_subject}.
func extractGraphPathFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return u.Host + u.Path
}

// extractFromFromURI extracts the {from_subject} parameter from an URI of type graph://{knowledge_graph_path}?from={from_subject}&to={to_subject}.
func extractFromFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	q := u.Query()
	return q.Get("from")
}

// extractToFromURI extracts the {to_subject} parameter from an URI of type graph://{knowledge_graph_path}?from={from_subject}&to={to_subject}.
func extractToFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	q := u.Query()
	return q.Get("to")
}

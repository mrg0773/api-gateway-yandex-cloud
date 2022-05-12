package main

import (
	"fmt"
	"net/http"
	"strings"
)

func inSlice(searchElement string, slice []string) bool {
	for _, el := range slice {
		if el == searchElement {
			return true
		}
	}
	return false
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}

	return strings.Join(request, "\n")
}

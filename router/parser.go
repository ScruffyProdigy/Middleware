package router

import (
	"github.com/ScruffyProdigy/TheRack/httper"
	"strings"
)

const (
	currentSectionIndex = "currentSection"
	parsedRouteIndex    = "parsedRoute"
	caseSensitiveIndex  = "caseSensitive"
	RouteEnd            = "/"
)

/*
	parser breaks down the request's URL path into a slice of strings
	later middleware will use it to direct control
*/

// V is a type you can cast your vars to in order to access the following functions
type V map[string]interface{}

// Parse will parse the URL into sections and store each section
func (vars V) Parse() []string {
	r := (httper.V)(vars).GetRequest()
	parsedRoute := strings.Split(r.URL.Path, "/")
	newParsedRoute := make([]string, 0, len(parsedRoute)+1)
	for _, section := range parsedRoute {
		if section != "" {
			l := len(newParsedRoute)
			newParsedRoute = newParsedRoute[0 : l+1]
			newParsedRoute[l] = section
		}
	}

	vars[parsedRouteIndex] = newParsedRoute
	return newParsedRoute
}

// CurrentSection() will return the section of the URL that we are currently looking at to try to route properly
func (vars V) CurrentSection() string {
	parsedRoute, ok := vars[parsedRouteIndex].([]string)
	if !ok {
		parsedRoute = V(vars).Parse()
	}

	index, ok := vars[currentSectionIndex].(int)
	if !ok {
		index = 0
	}

	if index < 0 || index >= len(parsedRoute) {
		return RouteEnd
	}

	result := parsedRoute[index]

	if !V(vars).IsCaseSensitive() {
		result = strings.ToLower(result)
	}

	return result
}

// nextSection() indicates that we have found a subrouter for the current section of the URL, and that it is time to move onto the next
func (vars V) nextSection() {
	index, ok := vars[currentSectionIndex].(int)
	if !ok {
		index = 0
	}
	vars[currentSectionIndex] = index + 1
}

//IsCaseSensitive() returns whether or not to assume that the route should be case sensitive
func (vars V) IsCaseSensitive() bool {
	result, ok := vars[caseSensitiveIndex].(bool)

	if ok && result {
		return true
	}
	return false
}

//SetCaseSensitive() sets the router to be case sensitive
func (vars V) SetCaseSensitive() {
	vars[caseSensitiveIndex] = true
}

//SetCaseInsensitive() is the default and sets the router to be case insensitive
func (vars V) SetCaseInsensitive() {
	vars[caseSensitiveIndex] = false
}

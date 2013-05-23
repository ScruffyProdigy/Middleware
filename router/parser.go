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

type V map[string]interface{}

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

func (vars V) nextSection() {
	index, ok := vars[currentSectionIndex].(int)
	if !ok {
		index = 0
	}
	vars[currentSectionIndex] = index + 1
}

func (vars V) IsCaseSensitive() bool {
	result, ok := vars[caseSensitiveIndex].(bool)

	if ok && result {
		return true
	}
	return false
}

func (vars V) SetCaseSensitive() {
	vars[caseSensitiveIndex] = true
}

func (vars V) SetCaseInsensitive() {
	vars[caseSensitiveIndex] = false
}

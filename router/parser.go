package router

import (
	"github.com/HairyMezican/TheRack/rack"
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

func Parse(vars rack.Vars) []string {
	r := rack.GetRequest(vars)
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

func CurrentSection(vars rack.Vars) string {
	parsedRoute, ok := vars[parsedRouteIndex].([]string)
	if !ok {
		parsedRoute = Parse(vars)
	}

	index, ok := vars[currentSectionIndex].(int)
	if !ok {
		index = 0
	}

	if index < 0 || index >= len(parsedRoute) {
		return RouteEnd
	}

	result := parsedRoute[index]

	if !IsCaseSensitive(vars) {
		result = strings.ToLower(result)
	}

	return result
}

func nextSection(vars rack.Vars) {
	index, ok := vars[currentSectionIndex].(int)
	if !ok {
		index = 0
	}
	vars[currentSectionIndex] = index + 1
}

func IsCaseSensitive(vars rack.Vars) bool {
	result, ok := vars[caseSensitiveIndex].(bool)

	if ok && result {
		return true
	}
	return false
}

func SetCaseSensitive(vars rack.Vars) {
	vars[caseSensitiveIndex] = true
}

func SetCaseInsensitive(vars rack.Vars) {
	vars[caseSensitiveIndex] = false
}

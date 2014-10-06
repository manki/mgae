package mgae

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const returnAllMatches = -1

var (
	variablePattern = regexp.MustCompile(`^%(.+)%$`)
  emptyMap = map[string]string{}
)

// Parses REST-style URL paths.
// For example, when called with
//     format="/cusotmer/%id%/%action%" and url path="/customer/44/edit",
// returns map
//     {"id": "44", "action": "edit"}
func ParsePath(format string, inputUrl *url.URL) (map[string]string, error) {
	formatParts := strings.Split(format, "/")
	pathParts := strings.Split(inputUrl.Path, "/")

	if len(formatParts) != len(pathParts) {
		return emptyMap, fmt.Errorf(
			"Number of components in path and format do not match. format=%q, path=%q.",
			format, inputUrl.Path)
	}

	variables := map[string]string{}
	for i := 0; i < len(formatParts); i++ {
		matches := variablePattern.FindStringSubmatch(formatParts[i])
		if len(matches) > 1 {
			varName := matches[1]
			variables[varName] = pathParts[i]
		} else if formatParts[i] != pathParts[i] {
			return emptyMap, fmt.Errorf(
				"Path and format do not match. format=%q, path=%q.", format, inputUrl.Path)
		}
	}

	return variables, nil
}

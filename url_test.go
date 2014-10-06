package mgae

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePath_noVariables(t *testing.T) {
	format := "/customer"
	path := "/customer/44/edit"
	_, err := ParsePath(format, urlForPath(path))
	assert.Error(t, err)
}

func TestParsePath_mismatchingVariables(t *testing.T) {
	format := "/customer/%id%"
	path := "/customer/44/edit"
	_, err := ParsePath(format, urlForPath(path))
	assert.Error(t, err)
}

func TestParsePath_mismatchingPath(t *testing.T) {
	format := "/customer/%id%/add"
	path := "/customer/44/edit"
	_, err := ParsePath(format, urlForPath(path))
	assert.Error(t, err)
}

func TestParsePath(t *testing.T) {
	format := "/customer/%id%/%action%"
	path := "/customer/44/edit"
	expected := map[string]string{
		"id": "44",
		"action": "edit",
	}
	vars, err := ParsePath(format, urlForPath(path))
	assert.NoError(t, err)
	assert.Equal(t, expected, vars)
}

func urlForPath(path string) *url.URL {
	inputUrl, err := url.Parse(path)
	if err != nil {
		panic("URL parsing failed.")
	}
	return inputUrl
}

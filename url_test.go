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

func TestParsePath_malformedFormat(t *testing.T) {
	path := "/customer/44"
	_, err := ParsePath("/customer/%id", urlForPath(path))
	assert.Error(t, err)

	_, err = ParsePath("/customer/%id%%", urlForPath(path))
	assert.Error(t, err)
}

func TestParsePath_onePartPath(t *testing.T) {
	format := "%id%"
	path := "44"
	expected := map[string]string{
		"id": "44",
	}
	vars, err := ParsePath(format, urlForPath(path))
	assert.NoError(t, err)
	assert.Equal(t, expected, vars)
}

func TestParsePath_multiPartsPath(t *testing.T) {
	format := "/customer/%id%/%action%"
	path := "/customer/44/edit"
	expected := map[string]string{
		"id":     "44",
		"action": "edit",
	}
	vars, err := ParsePath(format, urlForPath(path))
	assert.NoError(t, err)
	assert.Equal(t, expected, vars)
}

func TestParsePath_trailingSlashInPath(t *testing.T) {
	format := "/customer/%id%"
	path := "/customer/44/"
	expected := map[string]string{
		"id": "44",
	}
	vars, err := ParsePath(format, urlForPath(path))
	assert.NoError(t, err)
	assert.Equal(t, expected, vars)
}

func TestParsePath_trailingSlashInFormat(t *testing.T) {
	format := "/customer/%id%/"
	path := "/customer/44"
	_, err := ParsePath(format, urlForPath(path))
	assert.Error(t, err)
}

func urlForPath(path string) *url.URL {
	inputUrl, err := url.Parse(path)
	if err != nil {
		panic("URL parsing failed.")
	}
	return inputUrl
}

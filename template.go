// +build appengine

package mgae

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
)

// Encodes data into a 'data:' URI specified at
// https://developer.mozilla.org/en/data_URIs.
func dataUrl(data []byte, contentType string) template.URL {
	return template.URL(fmt.Sprintf(
		"data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(data)))
}

func inline(fileName, contentType string) (uri template.URL, err error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return dataUrl(data, contentType), nil
}

func NewTemplate(glob string) *template.Template {
	funcs := template.FuncMap(map[string]interface{}{
		"inline": inline,
	})

	tmpl := template.New("templates")
	tmpl.Funcs(funcs)
	return template.Must(tmpl.ParseGlob(glob))
}

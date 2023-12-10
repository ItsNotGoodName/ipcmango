//go:build dev

package webserver

import (
	"html/template"
)

func (t Renderer) Template(name string) (*template.Template, error) {
	nameHTML := name + ".html"
	tmpl, err := parseTemplate(nameHTML)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

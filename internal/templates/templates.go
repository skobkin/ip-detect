// Package templates embeds and exposes HTML templates.
package templates

import (
	"embed"
	"fmt"
	"html/template"
	"sync"
)

//go:embed client.gohtml
var templateFS embed.FS

var (
	clientTemplate     *template.Template
	clientTemplateOnce sync.Once
	errClientTemplate  error
)

// Client returns the compiled client data template.
func Client() (*template.Template, error) {
	clientTemplateOnce.Do(func() {
		var err error

		clientTemplate, err = template.ParseFS(templateFS, "client.gohtml")
		if err != nil {
			errClientTemplate = fmt.Errorf("parse client template: %w", err)

			return
		}

		errClientTemplate = nil
	})

	return clientTemplate, errClientTemplate
}

package app

import (
	"embed"
	"html/template"
	"strings"
	"sync"
)

// embeddedTemplates contains all HTML templates under app/templates.
//
//go:embed templates/*.html
var embeddedTemplates embed.FS

// tplCache caches parsed template sets keyed by a pipe-joined list of file names.
var tplCache sync.Map // map[string]*template.Template

// loadTemplates parses (or returns cached) templates from the embedded FS.
// Each name must include the "templates/" prefix (as produced by the //go:embed pattern).
func loadTemplates(names ...string) (*template.Template, error) {
	key := strings.Join(names, "|")
	if v, ok := tplCache.Load(key); ok {
		return v.(*template.Template), nil
	}
	t, err := template.ParseFS(embeddedTemplates, names...)
	if err != nil {
		return nil, err
	}
	tplCache.Store(key, t)
	return t, nil
}

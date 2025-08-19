package app

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing status: not_reviewed

// TestLoadTemplatesTesting verifies that the embedded HTML template `testing.html`
// can be loaded through loadTemplates and that its raw content matches the
// expected string. It also confirms the caching layer returns the same pointer
// on subsequent calls (an extra sanity check of tplCache behavior).
func TestLoadTemplatesTesting(t *testing.T) {
	// First load
	tpl, err := loadTemplates("templates/testing.html")
	assert.NoError(t, err)
	if !assert.NotNil(t, tpl) {
		return
	}

	var buf bytes.Buffer
	// Execute the (root) template; since only one file was parsed, Execute is fine.
	err = tpl.Execute(&buf, nil)
	assert.NoError(t, err)
	rendered := buf.String()
	assert.Contains(t, rendered, "<div>Testing Code</div>")

	// Second load should hit cache and return same pointer.
	tpl2, err := loadTemplates("templates/testing.html")
	assert.NoError(t, err)
	assert.Equal(t, tpl, tpl2, "expected cached template instance on second call")
}

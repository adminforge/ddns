package templates

import (
	"embed"
)

// Embed all templates in this folder.
// The pattern is relative to THIS file's directory.
//
//go:embed *.tmpl
var FS embed.FS

package ui

import (
	"embed"
)

// comment directive to instruct Go to store files from ui/html and ui/static folders in an embed.FS filesystem referenced by the global variable Files
//
//go:embed "html" "static"
var Files embed.FS

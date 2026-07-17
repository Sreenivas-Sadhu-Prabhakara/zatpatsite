// Package web embeds the builder UI assets.
package web

import "embed"

// FS holds the builder single-page app: index.html + css/js.
//
//go:embed index.html app.css app.js
var FS embed.FS

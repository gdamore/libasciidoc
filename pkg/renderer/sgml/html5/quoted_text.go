package html5

const (
	boldTextTmpl        = `<strong{{ if .ID }} id="{{ .ID }}"{{ end }}{{ if .Role }} class="{{ .Role }}"{{ end }}>{{ .Content }}</strong>`
	italicTextTmpl      = `<em{{ if .ID }} id="{{ .ID }}"{{ end }}{{ if .Role }} class="{{ .Role }}"{{ end }}>{{ .Content }}</em>`
	monospaceTextTmpl   = `<code{{ if .ID }} id="{{ .ID }}"{{ end }}{{ if .Role }} class="{{ .Role }}"{{ end }}>{{ .Content }}</code>`
	subscriptTextTmpl   = `<sub{{ if .ID }} id="{{ .ID }}"{{ end }}{{ if .Role }} class="{{ .Role }}"{{ end }}>{{ .Content }}</sub>`
	superscriptTextTmpl = `<sup{{ if .ID }} id="{{ .ID }}"{{ end }}{{ if .Role }} class="{{ .Role }}"{{ end }}>{{ .Content }}</sup>`
)

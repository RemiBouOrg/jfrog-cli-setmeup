package assets

import (
	_ "embed"
	"text/template"
)

//go:embed maven_settings.xml
var mavenSettingsTemplateStr string

func MavenSettingsTemplate() (*template.Template, error) {
	return template.New("maven_settings").Parse(mavenSettingsTemplateStr)
}

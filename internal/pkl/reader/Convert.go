package reader

import (
	"bytes"
	"html/template"
)

type YamlManifestToPklFile struct {
	ConvertUri string

	YamlManifests string `json:"yamlManifests"`

	// Kind containing apiVersion, containing the Resource
	CustomResourceTemplates map[string]map[string]string `json:"customResourceTemplates"`
}

func NewYamlManifestToPklFile(manifests string, crds map[string]map[string]string) YamlManifestToPklFile {
	return YamlManifestToPklFile{
		ConvertUri:              "package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib@1.0.1#/convert.pkl",
		YamlManifests:           manifests,
		CustomResourceTemplates: crds,
	}
}

// TODO replace this. It's silly...
func (y2p YamlManifestToPklFile) SillyHack() (string, error) {
	tmpl, err := template.New("").Parse(`
extends "package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib@1.0.1#/convert.pkl"
import "pkl:yaml"

resourcesToConvert: List<Mapping> =
	new yaml.Parser { useMapping = true }
	.parseAll(yamlManifests)
	.filterNonNull() as List<Mapping>

yamlManifests = """
{{ .YamlManifests }}
"""

{{ if .CustomResourceTemplates }}
customResourceTemplates {
{{- range $kind, $apiVersionCRD := .CustomResourceTemplates }}
	["{{ $kind }}"] {
		{{- range $apiVersion, $crd := $apiVersionCRD }}
		["{{ $apiVersion }}"] = import("{{ $crd }}")
		{{- end }}
	}
{{- end }}
}
{{- end }}
`)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, y2p); err != nil {
		return "", err
	}
	return out.String(), nil
}

// Convert

package luno

import "github.com/luno/luno-go/decimal"

{{range $i, $t := .API.AdditionalTypes -}}
{{if $t.Kind.IsEnum -}}
type {{$t.Name}} string

const (
	{{range $t.EnumProps.Values}}
	{{enumvalue $t.Name .}} {{$t.Name}} = "{{.}}"
	{{- end}}
)
{{end -}}
{{if $t.Kind.IsStruct -}}
type {{$t.Name}} struct {
	{{- range $i, $p := $t.StructProps.Properties}}
	{{if ne $i 0}}{{if $p.Description}}
	{{end}}{{end}}{{range $p.Description -}}
	// {{.}}
	{{end -}}
	{{propname $p.Name}} {{typename $p.Type}} `json:"{{$p.Name}}"`
	{{- end}}
}
{{end}}
{{end -}}
// vi: ft=go

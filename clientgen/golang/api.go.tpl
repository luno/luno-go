package luno

import (
	"context"

	"github.com/luno/luno-go/decimal"
)

{{range $e := .API.AllEndpoints}}
{{with $r := $e.Request}}
// {{$r.Name}} is the request struct for {{$e.Name}}.
type {{$r.Name}} struct {
	{{- range $i, $p := $r.StructProps.Properties}}
	{{if ne $i 0}}{{if $p.Description}}
	{{end}}{{end}}{{range $p.Description -}}
	// {{.}}
	{{end -}}{{if .Required -}}
	//
	// required: true
	{{end -}}
	{{propname $p.Name}} {{typename $p.Type}} `json:"{{$p.Name}}" url:"{{$p.Name}}"`
	{{- end}}
}

{{end -}}
{{with $r := $e.Response -}}
// {{$r.Name}} is the response struct for {{$e.Name}}.
type {{$r.Name}} struct {
	{{- range $i, $p := $r.StructProps.Properties}}
	{{if ne $i 0}}{{if $p.Description}}
	{{end}}{{end}}{{range $p.Description -}}
	// {{.}}
	{{end -}}
	{{propname $p.Name}} {{typename $p.Type}} `json:"{{$p.Name}}"`
	{{- end}}
}

{{end -}}
// {{.Name}} makes a call to {{.Method}} {{.Path}}.
//
{{range .Description -}}
// {{.}}
{{end -}}
func (cl *Client) {{$e.Name}}(ctx context.Context, req *{{$e.Request.Name}}) (*{{$e.Response.Name}}, error) {
	var res {{$e.Response.Name}}
	err := cl.do(ctx, "{{$e.Method}}", "{{$e.Path}}", req, &res, {{$e.RequiresAuth}})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
{{- end}}

// vi: ft=go

<?php declare(strict_types=1);

namespace Luno\Request;

{{- $i := importtypes .Request}}{{if $i}}
{{range $i}}
use Luno\Types\{{.Name}};
{{- end}}{{end}}

class {{.Name}} extends AbstractRequest
{
  {{- range $i, $p := .Request.StructProps.Properties}}
  {{- if $p.Description}}{{if gt $i 0}}
{{end}}
  /**
  {{- range $p.Description}}
   * {{.}}
  {{- end}}
   */ 
  {{- end}}
  protected ${{$p.Name}};
  {{- end}}
  {{range $i, $p := .Request.StructProps.Properties}}
  {{- template "getter" $p}}
  {{- template "setter" $p}}
  {{- end -}}
}

// vi: ft=php

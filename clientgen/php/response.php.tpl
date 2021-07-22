<?php declare(strict_types=1);

namespace Luno\Response;

{{- $i := importtypes .Response}}{{if $i}}
{{range $i}}
use Luno\Types\{{.Name}};
{{- end}}{{end}}

class {{.Name}} extends AbstractResponse
{
  {{- range $i, $p := .Response.StructProps.Properties}}
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
  {{range $i, $p := .Response.StructProps.Properties}}
  {{- template "getter" $p}}
  {{- template "setter" $p}}
  {{- end -}}
}

// vi: ft=php

<?php declare(strict_types=1);

namespace Luno\Types;

class {{.Name}}
{
  {{- if .Kind.IsEnum}}
  {{- range .EnumProps.Values}}
  const {{.}} = "{{.}}";
  {{- end}}
  {{- else if .Kind.IsStruct}}
  {{- range $i, $p := .StructProps.Properties}}
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
  {{range $i, $p := .StructProps.Properties}}
  {{- template "getter" $p}}
  {{- template "setter" $p}}
  {{- end -}}
{{end}}
}

// vi: ft=php

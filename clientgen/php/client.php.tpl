<?php declare(strict_types=1);

namespace Luno;

class Client extends AbstractClient
{
  {{- range $e := .AllEndpoints}}
  /**
   * {{$e.Name}} makes a call to {{$e.Method}} {{$e.Path}}.
   *
{{- range $e.Description}}
   * {{.}}
{{- end}}
   */ 
  public function {{$e.Name}}(Request\{{$e.Name}} $req): Response\{{$e.Name}}
  {
    $res = $this->do("{{$e.Method}}", "{{$e.Path}}", $req, {{if $e.RequiresAuth}}true{{else}}false{{end}});
    $mapper = new \JsonMapper();
    return $mapper->map($res, new Response\{{$e.Name}});
  }
{{end -}}
}

// vi: ft=php

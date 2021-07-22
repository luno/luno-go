{{define "getter"}}
  /**
   * @return {{typename .Type true}}
   */
  public function {{gettername .Name}}(): {{typename .Type false}}
  {
    if (!isset($this->{{.Name}})) {
      return {{zeroval .Type.Kind}};
    }
    return $this->{{.Name}};
  }
{{end}}

{{define "setter"}}
  /**
   * @param {{typename .Type true}} ${{varname .Name}}
   */
  public function {{settername .Name}}({{typename .Type false}} ${{varname .Name}})
  {
    $this->{{.Name}} = ${{varname .Name}};
  }
{{end}}

// vi: ft=php

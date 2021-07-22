from .base_client import BaseClient


class Client(BaseClient):
    """
    Python SDK for the Luno API.

    Example usage:

      from luno_python.client import Client


      c = Client(api_key_id='key_id', api_key_secret='key_secret')
      try:
        res = c.get_ticker(pair='XBTZAR')
        print res
      except Exception as e:
        print e
    """
{{range .AllEndpoints}}
    def {{funcname .Name}}(self{{range .Request.StructProps.Properties}}, {{.Name}}{{if not .Required}}=None{{end}}{{end}}):
        """Makes a call to {{.Method}} {{.Path}}.
{{range .Description}}
{{if .}}        {{.}}{{end}}
{{- end}}
{{range .Request.StructProps.Properties}}
        :param {{.Name}}: {{$name := .Name}}{{range $i, $l := .Description}}{{if $l}}{{if gt $i 0}}{{paramlen $name}}        {{end}}{{end}}{{$l}}
{{end}}        :type {{.Name}}: {{typename .Type}}
{{- end}}
        """
{{- if .Request.StructProps.Properties}}
        req = {
{{- range .Request.StructProps.Properties}}
            '{{.Name}}': {{.Name}}, 
{{- end}}
        }
{{- end}}
        return self.do('{{.Method}}', '{{.Path}}', req={{if .Request.StructProps.Properties}}req{{else}}None{{end}}, auth={{boolval .RequiresAuth}})
{{end}}

# vi: ft=python

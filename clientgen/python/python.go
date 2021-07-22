package python

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/luno/luno-go/clientgen/clientgen"
)

// Generate generates a Python API client.
func Generate(api clientgen.API) ([]clientgen.File, error) {
	var fl []clientgen.File

	// Generate client.py
	f, err := generatePythonFile("luno_python/client.py", "client.py.tpl", api)
	if err != nil {
		return nil, err
	}
	fl = append(fl, f)

	return fl, nil
}

func generatePythonFile(fileName, tplName string, c interface{}) (
	clientgen.File, error) {

	funcMap := template.FuncMap{
		"funcname": funcname,
		"typename": typename,
		"paramlen": paramlen,
		"boolval":  boolval,
	}

	filenames := []string{
		filepath.Join(os.Getenv("GOPATH"),
			"src/bitx/fe/publicapi/tools/clientgen/python", tplName),
	}
	tpl, err := template.New(tplName).Funcs(funcMap).ParseFiles(filenames...)
	if err != nil {
		return clientgen.File{}, err
	}

	f := clientgen.NewFile(fileName)

	if err := tpl.Execute(f, c); err != nil {
		return clientgen.File{}, err
	}

	return f, nil
}

func funcname(in string) string {
	b := bytes.NewBuffer(nil)
	for i, r := range in {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteRune('_')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return string(b.Bytes())
}

func typename(t clientgen.Type) string {
	switch t.Kind {
	case clientgen.KindFloat:
		return "float"
	case clientgen.KindInteger:
		return "int"
	case clientgen.KindDecimal:
		return "float"
	case clientgen.KindString:
		return "str"
	case clientgen.KindBoolean:
		return "bool"
	case clientgen.KindEnum:
		return "str"
	case clientgen.KindStruct:
		return t.Name
	case clientgen.KindArray:
		return "list"
	case clientgen.KindTimestamp:
		return "int"
	}
	return ""
}

func paramlen(p string) string {
	return strings.Repeat(" ", len(p)+9)
}

func boolval(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

package php

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/luno/luno-go/clientgen/clientgen"
)

// Generate generates a PHP API client.
func Generate(api clientgen.API) ([]clientgen.File, error) {
	var fl []clientgen.File

	// In PHP, every class is in a separate file.

	// Generate Client.php
	f, err := generatePHPFile("src/Luno/Client.php", "client.php.tpl", api)
	if err != nil {
		return nil, err
	}
	fl = append(fl, f)

	// Generate requests + responses.
	for _, e := range api.AllEndpoints() {
		f, err := generatePHPFile(
			"src/Luno/Request/"+e.Name+".php", "request.php.tpl", e)
		if err != nil {
			return nil, err
		}
		fl = append(fl, f)

		f, err = generatePHPFile(
			"src/Luno/Response/"+e.Name+".php", "response.php.tpl", e)
		if err != nil {
			return nil, err
		}
		fl = append(fl, f)
	}

	// Generate addition structs.
	for _, t := range api.AdditionalTypes {
		f, err = generatePHPFile(
			"src/Luno/Types/"+t.Name+".php", "type.php.tpl", t)
		if err != nil {
			return nil, err
		}
		fl = append(fl, f)
	}

	return fl, nil
}

func generatePHPFile(fileName, tplName string, c interface{}) (
	clientgen.File, error) {

	funcMap := template.FuncMap{
		"gettername":  gettername,
		"settername":  settername,
		"typename":    typename,
		"varname":     varname,
		"importtypes": importtypes,
		"zeroval":     zeroval,
	}

	filenames := []string{
		filepath.Join(os.Getenv("GOPATH"),
			"src/bitx/fe/publicapi/tools/clientgen/php", tplName),
		filepath.Join(os.Getenv("GOPATH"),
			"src/bitx/fe/publicapi/tools/clientgen/php/partials.php.tpl"),
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

func propfunc(s string) string {
	b := bytes.NewBuffer(nil)
	for i := 0; i < len(s); i++ {
		r := s[i]
		if r == '_' {
			if i < len(s)-1 {
				fmt.Fprintf(b, strings.ToUpper(string(s[i+1])))
				i++
			}
			continue
		}
		fmt.Fprint(b, string(r))
	}
	return string(b.Bytes())
}

func gettername(s string) string {
	return "get" + strings.Title(propfunc(s))
}

func settername(s string) string {
	return "set" + strings.Title(propfunc(s))
}

func varname(s string) string {
	return propfunc(s)
}

func typename(t clientgen.Type, inComment bool) string {
	switch t.Kind {
	case clientgen.KindFloat:
		return "float"
	case clientgen.KindInteger:
		return "int"
	case clientgen.KindDecimal:
		return "float" // TODO
	case clientgen.KindString:
		return "string"
	case clientgen.KindBoolean:
		return "bool"
	case clientgen.KindEnum:
		return "string"
	case clientgen.KindStruct:
		return t.Name
	case clientgen.KindArray:
		if inComment {
			return "\\Luno\\Types\\" + typename(t.ArrayProps.Type, true) + "[]"
		}
		return "array"
	case clientgen.KindTimestamp:
		return "int"
	case clientgen.KindObject:
		return "array"
	}
	return ""
}

func importtypes(t clientgen.Type) []clientgen.Type {
	if t.Kind != clientgen.KindStruct {
		return nil
	}

	if t.StructProps == nil {
		return nil
	}

	tm := make(map[string]bool)
	var tl []clientgen.Type
	for _, p := range t.StructProps.Properties {
		if p.Type.Kind == clientgen.KindArray {
			ttl := importtypes(p.Type.ArrayProps.Type)
			for _, t := range ttl {
				if tm[t.Name] {
					continue
				}
				tm[t.Name] = true
				tl = append(tl, t)
			}
			continue
		}

		if p.Type.Kind != clientgen.KindStruct {
			continue
		}

		if tm[p.Type.Name] {
			continue
		}
		tm[p.Type.Name] = true
		tl = append(tl, p.Type)
	}

	sort.Slice(tl, func(i, j int) bool {
		return tl[i].Name < tl[j].Name
	})

	return tl
}

func zeroval(k clientgen.Kind) string {
	switch k {
	case clientgen.KindFloat:
		return "0"
	case clientgen.KindInteger:
		return "0"
	case clientgen.KindDecimal:
		return "0" // TODO
	case clientgen.KindString:
		return "\"\""
	case clientgen.KindBoolean:
		return "false"
	case clientgen.KindEnum:
		return "\"\""
	case clientgen.KindStruct:
		return "null"
	case clientgen.KindArray, clientgen.KindObject:
		return "[]"
	case clientgen.KindTimestamp:
		return "0"
	}
	return ""
}

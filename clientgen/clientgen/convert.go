package clientgen

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	openapi "github.com/go-openapi/spec"
)

type Generator func(API) ([]File, error)

func ConvertSpec(spec *openapi.Swagger) (api API) {
	api.Description = spec.Info.Description

	seen := make(map[string]bool)

	for _, t := range spec.Tags {
		var s Section
		s.Name = t.Name
		s.Description = t.Description
		api.Sections = append(api.Sections, s)
	}

	for path, items := range spec.Paths.Paths {
		convertFunc(spec, &api, http.MethodGet, path, items.Get, seen)
		convertFunc(spec, &api, http.MethodPost, path, items.Post, seen)
		convertFunc(spec, &api, http.MethodPut, path, items.Put, seen)
		convertFunc(spec, &api, http.MethodPatch, path, items.Patch, seen)
		convertFunc(spec, &api, http.MethodDelete, path, items.Delete, seen)
	}

	for _, s := range api.Sections {
		sort.Slice(s.Endpoints, func(i, j int) bool {
			return s.Endpoints[i].Name < s.Endpoints[j].Name
		})
	}

	sort.Slice(api.AdditionalTypes, func(i, j int) bool {
		return api.AdditionalTypes[i].Name < api.AdditionalTypes[j].Name
	})

	return
}

func convertFunc(spec *openapi.Swagger, api *API, method, path string,
	op *openapi.Operation, seen map[string]bool) {

	if op == nil {
		return
	}

	id := strings.Title(op.ID)

	e := Endpoint{
		Method:       method,
		Path:         path,
		Name:         id,
		Description:  convertDescription(op.Description),
		RequiresAuth: strings.Contains(op.Description, "Permissions required"),
	}

	_, requestSchema := lookupRef(spec, "#/definitions/"+op.ID+"Request")
	if requestSchema == nil {
		e.Request = convertParams(spec, api, id+"Request", op, seen)
	} else {
		e.Request = convertSchema(spec, api, id+"Request", requestSchema, seen)
	}

	for code, r := range op.Responses.StatusCodeResponses {
		if code != http.StatusOK {
			continue
		}
		responseSchema := r.Schema
		for responseSchema.Ref.String() != "" {
			_, responseSchema = lookupRef(spec, responseSchema.Ref.String())
		}
		if responseSchema == nil {
			return
		}
		e.Response = convertSchema(spec, api, id+"Response", responseSchema, seen)
	}

	for _, t := range op.Tags {
		for i, s := range api.Sections {
			if s.Name == t {
				api.Sections[i].Endpoints = append(api.Sections[i].Endpoints, e)
			}
		}
	}
}

func convertDescription(s string) []string {
	metaRx := regexp.MustCompile("^[a-z]+:")
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	var newLines []string
	for _, l := range lines {
		if metaRx.MatchString(l) {
			continue
		}
		newLines = append(newLines, l)
	}
	if len(newLines) == 0 {
		return newLines
	}
	if newLines[len(newLines)-1] == "" {
		newLines = newLines[:len(newLines)-1]
	}
	return newLines
}

func convertSchema(spec *openapi.Swagger, api *API, name string,
	schema *openapi.Schema, seen map[string]bool) (typ Type) {

	typ.Kind = KindStruct
	typ.Name = name

	requiredMap := make(map[string]bool)
	for _, r := range schema.Required {
		requiredMap[r] = true
	}

	var structProps StructProps
	for n, p := range schema.Properties {
		t := convertType(spec, api, schemaToSimpleSchema(p), seen)
		if t.Kind == KindUnknown {
			continue
		}

		var example string
		if p.Example != nil {
			example = fmt.Sprintf("%v", p.Example)
		}

		structProps.Properties = append(structProps.Properties, Property{
			Description: convertDescription(p.Description),
			Name:        n,
			Type:        t,
			Required:    requiredMap[n],
			Example:     example,
		})
	}

	sort.Sort(sortedStructProps(structProps.Properties))

	typ.StructProps = &structProps

	return
}

func convertParams(spec *openapi.Swagger, api *API, name string,
	op *openapi.Operation, seen map[string]bool) (typ Type) {

	typ.Kind = KindStruct
	typ.Name = name

	var structProps StructProps

	for _, p := range op.Parameters {
		t := convertType(spec, api, paramToSimpleSchema(p), seen)
		if t.Kind == KindUnknown {
			continue
		}

		var example string
		if p.Example != nil {
			example = fmt.Sprintf("%v", p.Example)
		}

		structProps.Properties = append(structProps.Properties, Property{
			Description: convertDescription(p.Description),
			Name:        p.Name,
			Type:        t,
			Required:    p.Required,
			Example:     example,
		})
	}

	sort.Sort(sortedStructProps(structProps.Properties))

	typ.StructProps = &structProps

	return
}

type sortedStructProps []Property

func (s sortedStructProps) Len() int      { return len(s) }
func (s sortedStructProps) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s sortedStructProps) Less(i, j int) bool {
	if s[i].Required == s[j].Required {
		return s[i].Name < s[j].Name
	}
	return s[i].Required
}

type simpleSchema struct {
	Enum       []interface{}
	Extensions map[string]interface{}
	Type       openapi.StringOrArray
	Format     string
	Ref        openapi.Ref
	Items      *openapi.SchemaOrArray
	Additional *openapi.SchemaOrBool
}

func schemaToSimpleSchema(p openapi.Schema) *simpleSchema {
	return &simpleSchema{
		Enum:       p.Enum,
		Extensions: p.Extensions,
		Type:       p.Type,
		Format:     p.Format,
		Ref:        p.Ref,
		Items:      p.Items,
		Additional: p.AdditionalProperties,
	}
}

func paramToSimpleSchema(p openapi.Parameter) *simpleSchema {
	return &simpleSchema{
		Enum:       p.Enum,
		Extensions: p.Extensions,
		Type:       openapi.StringOrArray{p.Type},
		Format:     p.Format,
	}
}

func convertType(spec *openapi.Swagger, api *API, p *simpleSchema,
	seen map[string]bool) (typ Type) {

	if len(p.Enum) > 0 {
		goname := p.Extensions["x-go-name"].(string)

		var vl []string
		for _, v := range p.Enum {
			value := v.(string)
			vl = append(vl, value)
		}

		if !seen[goname] {
			enumProps := new(EnumProps)
			enumProps.AddAll(vl)
			api.AdditionalTypes = append(api.AdditionalTypes, Type{
				Kind:      KindEnum,
				Name:      goname,
				EnumProps: enumProps,
			})
			seen[goname] = true
		} else {
			for _, e := range api.AdditionalTypes {
				if e.Kind == KindEnum && e.Name == goname {
					enumProps := e.EnumProps
					enumProps.AddAll(vl)
					sort.Strings(enumProps.Values)
					break
				}
			}
		}

		typ.Name = goname
		typ.Kind = KindEnum

		return
	}

	if len(p.Type) > 0 {
		if p.Type[0] == "string" && p.Format == "amount" {
			typ.Kind = KindDecimal
			return
		}

		if p.Type[0] == "string" && p.Format == "timestamp" {
			typ.Kind = KindTimestamp
			return
		}

		if p.Type[0] == "array" {
			var t Type
			if p.Items != nil {
				itemSchema := schemaToSimpleSchema(*p.Items.Schema)
				t = convertType(spec, api, itemSchema, seen)
			} else {
				t = Type{Kind: KindString}
			}
			typ.Kind = KindArray
			typ.ArrayProps = &ArrayProps{Type: t}
			return
		}

		if p.Type[0] == "object" && p.Additional.Allows {
			valueSchema := schemaToSimpleSchema(*p.Additional.Schema)
			t := convertType(spec, api, valueSchema, seen)
			typ.Kind = KindObject
			typ.ObjectProps = &ObjectProps{ValueType: t}
			return
		}

		if k := convertScalarType(p.Type[0]); k != KindUnknown {
			typ.Kind = k
			return
		}
	}

	name, schema := lookupRef(spec, p.Ref.String())
	if schema == nil {
		return
	}

	typ.Kind = KindStruct
	typ.Name = name

	pkgName := getPackageName(name, schema.Extensions)
	newTyp := convertSchema(spec, api, name, schema, seen)
	typ.StructProps = newTyp.StructProps

	if !seen[pkgName] {
		api.AdditionalTypes = append(api.AdditionalTypes, newTyp)
		seen[pkgName] = true
	}

	return
}

func convertScalarType(t string) Kind {
	if t == "number" {
		return KindFloat
	}
	if t == "integer" {
		return KindInteger
	}
	if t == "string" {
		return KindString
	}
	if t == "boolean" {
		return KindBoolean
	}
	return KindUnknown
}

func convertEnumName(s string) string {
	return strings.Title(strings.ToLower(s))
}

func lookupRef(spec *openapi.Swagger, ref string) (string, *openapi.Schema) {
	if !strings.HasPrefix(ref, "#/") || len(ref) < 3 {
		return "", nil
	}
	parts := strings.SplitN(ref[2:], "/", 2)
	if len(parts) != 2 {
		return "", nil
	}
	switch parts[0] {
	case "definitions":
		s, ok := spec.Definitions[parts[1]]
		if !ok {
			return "", nil
		}
		return parts[1], &s
	case "responses":
		s, ok := spec.Responses[parts[1]]
		if !ok {
			return "", nil
		}
		return parts[1], s.Schema
	}
	return "", nil
}

func getPackageName(name string, extensions map[string]interface{}) string {
	pkgi := extensions["x-go-package"]
	if pkgi == nil {
		return ""
	}
	pkg := pkgi.(string)
	ni := extensions["x-go-name"]
	if ni != nil {
		name = ni.(string)
	}
	return pkg + "." + name
}

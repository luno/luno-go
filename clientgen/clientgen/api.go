package clientgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"
)

type Kind int

const (
	KindUnknown Kind = iota
	KindFloat
	KindInteger
	KindDecimal
	KindString
	KindBoolean
	KindEnum
	KindStruct
	KindArray
	KindTimestamp
	KindObject
)

const apiBaseURL = "https://api.luno.com"

func (k Kind) IsStruct() bool {
	return k == KindStruct
}

func (k Kind) IsEnum() bool {
	return k == KindEnum
}

type API struct {
	Description     string    `json:"description"`
	Sections        []Section `json:"sections"`
	AdditionalTypes []Type    `json:"additional_types"`
}

// Return all endpoints in all sections. NOTE: endpoints with multiple tags will
// appear multiple time in this list.
func (a API) AllEndpoints() (endpoints []Endpoint) {
	for _, s := range a.Sections {
		endpoints = append(endpoints, s.Endpoints...)
	}
	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].Name < endpoints[j].Name
	})
	return
}

// DescriptionHTML returns the API description an HTML chunk.
func (a API) DescriptionHTML() template.HTML {
	return template.HTML(a.Description)
}

type Section struct {
	Name        string     `json:"title"`
	Description string     `json:"description"`
	Endpoints   []Endpoint `json:"endpoints"`
}

// ID returns the section name without spaces, which can be used as a DOM
// element ID.
func (s Section) ID() string {
	return strings.Replace(s.Name, " ", "-", -1)
}

// Title returns the section's name in title case.
func (s Section) Title() string {
	return strings.Title(s.Name)
}

// DescriptionHTML returns the section's description as an HTML chunk.
func (s Section) DescriptionHTML() template.HTML {
	return template.HTML(s.Description)
}

type Endpoint struct {
	Method       string   `json:"method"`
	Path         string   `json:"path"`
	Name         string   `json:"name"`
	Description  []string `json:"description"`
	RequiresAuth bool     `json:"requires_auth"`
	Response     Type     `json:"response"`
	Request      Type     `json:"request"`
}

// ID returns the endpoint name without spaces, which can be used as a DOM
// element ID.
func (e Endpoint) ID() string {
	return strings.ToLower(e.Name)
}

// DescriptionHTML returns the endpoint's description as an HTML chunk.
//
// Endpoint descriptions are taken from go docs, which are hard-wrapped. Since
// this is outputting HTML, the intra-paragraph line breaks will show as spaces.
func (e Endpoint) DescriptionHTML() template.HTML {
	desc := strings.Join(e.Description, "\n")
	desc = "<p>" + strings.Replace(desc, "\n\n", "</p><p>", -1) + "</p>"
	return template.HTML(desc)
}

// Definition returns the endpoint's HTTP method and full URL.
func (e Endpoint) Definition() string {
	return fmt.Sprintf("%s %s%s",
		e.Method, apiBaseURL, e.Path)
}

// ExampleRequest generates an example cURL request for this endpoint.
func (e Endpoint) ExampleRequest() string {
	var flags []string

	if e.RequiresAuth {
		flags = append(flags, fmt.Sprintf("-u api_key_id:api_key_secret"))
	}

	if e.Method != http.MethodGet {
		flags = append(flags, fmt.Sprintf("-X %s", e.Method))
	}

	path := e.Path
	for _, p := range e.Request.StructProps.Properties {
		if !p.Required {
			continue
		}
		if strings.Contains(e.Path, "{"+p.Name+"}") {
			path = strings.Replace(path, "{"+p.Name+"}", p.Example, -1)
			continue
		}
		if p.Type.Kind == KindArray {
			examples := strings.Split(p.Example, ",")
			for _, e := range examples {
				flags = append(flags, fmt.Sprintf("-F '%s=%s'", p.Name, e))
			}
			continue
		}
		flags = append(flags, fmt.Sprintf("-F '%s=%s'", p.Name, p.Example))
	}

	flags = append(flags, fmt.Sprintf("%s%s", apiBaseURL, path))

	b := bytes.NewBuffer(nil)
	b.WriteString("curl")
	b.WriteString(" ")
	b.WriteString(flags[0])
	flags = flags[1:]

	for _, f := range flags {
		b.WriteString(" ")
		b.WriteString("\\\n")
		b.WriteString("     ")
		b.WriteString(f)
	}

	return string(b.Bytes())
}

// ResponseStructure returns a JSON object with the same structure as this
// endpoint's response. The values are all zero values for the corresponding
// type.
func (e Endpoint) ResponseStructure() string {
	res := e.Response.JSONStructure()
	var i interface{}
	if err := json.Unmarshal([]byte(res), &i); err != nil {
		return ""
	}
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

type Type struct {
	Kind Kind   `json:"kind"`
	Name string `json:"name"`

	StructProps *StructProps `json:"struct_props,omitempty"`
	EnumProps   *EnumProps   `json:"enum_props,omitempty"`
	ArrayProps  *ArrayProps  `json:"array_props,omitempty"`
	ObjectProps *ObjectProps `json:"object_props,omitempty"`
}

// Doc returns the type's type as it should appear in API docs.
func (t Type) Doc() string {
	switch t.Kind {
	case KindFloat:
		return "float"
	case KindInteger:
		return "integer"
	case KindDecimal:
		return "decimal"
	case KindString:
		return "string"
	case KindBoolean:
		return "boolean"
	case KindEnum:
		return "enum"
	case KindStruct:
		return "object"
	case KindArray:
		return "array"
	case KindTimestamp:
		return "timestamp"
	}
	return ""
}

// SubDoc returns the type's subtype as it should appear in API docs.
func (t Type) SubDoc() string {
	switch t.Kind {
	case KindEnum:
		return "string"
	case KindArray:
		return t.ArrayProps.Type.Doc()
	}
	return ""
}

// JSONStructure returns the string representing a JSON object with the same
// structure as this type. Zero values are used to ensure the string can be
// unmarshalled.
func (t Type) JSONStructure() string {
	switch t.Kind {
	case KindFloat:
		return "0.0"
	case KindInteger:
		return "0"
	case KindDecimal:
		return `"0.0"`
	case KindString:
		return `""`
	case KindBoolean:
		return "false"
	case KindEnum:
		return `""`
	case KindStruct:
		if t.StructProps == nil {
			return ""
		}
		b := bytes.NewBuffer(nil)
		b.WriteString("{")
		for i, p := range t.StructProps.Properties {
			if i > 0 {
				b.WriteString(",")
			}
			s := p.Type.JSONStructure()
			b.WriteString(fmt.Sprintf(`"%s":%s`, p.Name, s))
		}
		b.WriteString("}")
		return string(b.Bytes())
	case KindArray:
		return "[" + t.ArrayProps.Type.JSONStructure() + "]"
	case KindTimestamp:
		return "0"
	}
	return ""
}

type StructProps struct {
	Properties []Property `json:"properties"`
}

type EnumProps struct {
	Values []string `json:"values"`
}

type ObjectProps struct {
	// Only string keys are supported by go-swagger
	// https://goswagger.io/faq/faq_spec.html#maps-as-swagger-parameters
	ValueType Type
}

func (ep *EnumProps) AddAll(vl []string) {
	m := make(map[string]bool)
	for _, v := range ep.Values {
		m[v] = true
	}
	for _, v := range vl {
		if !m[v] {
			ep.Values = append(ep.Values, v)
		}
	}
}

type ArrayProps struct {
	Type Type `json:"type"`
}

type Property struct {
	Name        string   `json:"name"`
	Description []string `json:"description,omitempty"`
	Type        Type     `json:"type"`
	Example     string   `json:"example,omitempty"`

	Required bool `json:"required,omitempty"`
}

// DescriptionHTML returns the property's description as an HTML chunk.
//
// Property descriptions in the API docs are shown in a table, so shouldn't be
// wrapped with <p> tags. Instead, replace intra-paragraph line breaks with
// spaces and inter-paragraph line breaks with <br> tags.
func (p Property) DescriptionHTML() template.HTML {
	desc := strings.Join(p.Description, "\n")
	desc = strings.Replace(desc, "\n\n", "<br><br>", -1)
	desc = strings.Replace(desc, "\n", " ", -1)
	return template.HTML(desc)
}

package clientgen

import (
	"encoding/json"
	openapi "github.com/go-openapi/spec"
	"net/http"
)

func LoadSpec(fn string) (*openapi.Swagger, error) {
	f, err := http.Get(fn)
	if err != nil {
		return nil, err
	}
	defer f.Body.Close()

	var spec openapi.Swagger
	if err := json.NewDecoder(f.Body).Decode(&spec); err != nil {
		return nil, err
	}

	return &spec, nil
}

package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Opts struct {
	Version      string
	Title        string
	Description  string
	MethodOffset int
}

func New(opts Opts) OpenAPI {
	return OpenAPI{
		OpenAPI:    opts.Version,
		Paths:      map[string]map[string]any{},
		Components: map[string]map[string]any{"schemas": {}},
		offset:     opts.MethodOffset,
		Info: Info{
			Title:       opts.Title,
			Description: opts.Description,
		},
	}
}

func (o OpenAPI) Convert(info []byte) ([]byte, error) {
	var micro Micro
	if err := json.Unmarshal(info, &micro); err != nil {
		return nil, err
	}
	o.Info.Version = micro.Version

	for _, v := range micro.Endpoints {
		rawSchema, err := parseSchema(v.Metadata.ResponseSchema)
		if err != nil {
			return nil, err
		}

		inlineSchema, defs, err := convertJSONSchemaToOpenAPI(rawSchema)
		if err != nil {
			return nil, err
		}

		parsedParams, err := parseParams(v.Metadata.Params)
		if err != nil {
			return nil, err
		}

		method, path, params := subjToPath(v.Subject, parsedParams, o.offset)

		for name, def := range defs {
			o.Components["schemas"][name] = def
		}

		if _, exists := o.Paths[path]; !exists {
			o.Paths[path] = map[string]any{}
		}

		o.Paths[path] = map[string]any{
			method: map[string]any{
				"tags":        []string{v.QueueGroup},
				"summary":     v.Metadata.Description,
				"description": v.Metadata.Description,
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Success",
						"content": map[string]any{
							v.Metadata.Format: map[string]any{
								"schema": inlineSchema,
							},
						},
					},
				},
			},
		}
		parameters := []map[string]any{}
		if len(params) > 0 {
			for _, p := range parsedParams {
				parameters = append(parameters, map[string]any{
					"name":     p.Name,
					"required": p.Required,
					"in":       p.In,
					"schema":   p.Schema,
				})
			}
		}
		o.Paths[path]["parameters"] = parameters

	}

	return yaml.Marshal(o)
}

func fixRefs(v any) any {
	switch x := v.(type) {
	case map[string]any:
		if ref, ok := x["$ref"].(string); ok {
			if strings.HasPrefix(ref, "#/$defs/") {
				name := strings.TrimPrefix(ref, "#/$defs/")
				x["$ref"] = "#/components/schemas/" + name
			}
		}
		for k, child := range x {
			x[k] = fixRefs(child)
		}
		return x
	case []any:
		for i, item := range x {
			x[i] = fixRefs(item)
		}
		return x
	default:
		return v
	}
}

func convertJSONSchemaToOpenAPI(schema map[string]any) (map[string]any, map[string]any, error) {
	components := map[string]map[string]any{
		"schemas": map[string]any{},
	}

	if defs, ok := schema["$defs"].(map[string]any); ok {
		for name, def := range defs {
			components["schemas"][name] = def.(map[string]any)
		}
		delete(schema, "$defs")
	}

	schema = fixRefs(schema).(map[string]any)
	c := fixRefs(components["schemas"]).(map[string]any)

	delete(schema, "$schema")

	return schema, c, nil

}

func subjToPath(subject string, pathParams []Param, methodOffset int) (string, string, []string) {
	split := strings.Split(subject, ".")
	method := strings.ToLower(split[methodOffset])
	var params []string

	parsed := split[4:]
	service := split[2]

	out := fmt.Sprintf("/%s", service)
	counter := 0

	for _, p := range parsed {
		if p == "*" && len(pathParams) != 0 {
			p = fmt.Sprintf("{%s}", pathParams[counter].Name)
			params = append(params, p)
			counter++
		}
		out += "/" + p
	}

	return method, out, params

}

func parseSchema(schema string) (map[string]any, error) {
	var out map[string]any
	if schema == "" {
		return nil, nil
	}
	return out, json.Unmarshal([]byte(schema), &out)
}

func parseParams(params string) ([]Param, error) {
	var out []Param
	if params == "" {
		return nil, nil
	}
	return out, json.Unmarshal([]byte(params), &out)
}

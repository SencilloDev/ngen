package openapi

type Metadata struct {
	Description    string `json:"description,omitempty"`
	Format         string `json:"format,omitempty"`
	RequestSchema  string `json:"request_schema,omitempty"`
	ResponseSchema string `json:"response_schema,omitempty"`
	Params         string `json:"params,omitempty"`
}
type Endpoints struct {
	Name       string   `json:"name,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	QueueGroup string   `json:"queue_group,omitempty"`
	Metadata   Metadata `json:"metadata,omitempty"`
}
type Micro struct {
	Name        string      `json:"name,omitempty"`
	ID          string      `json:"id,omitempty"`
	Version     string      `json:"version,omitempty"`
	Metadata    Metadata    `json:"metadata,omitempty"`
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Endpoints   []Endpoints `json:"endpoints,omitempty"`
}

type MicroInfo struct {
	Info Micro `json:"info"`
}

type Param struct {
	Name     string      `json:"name,omitempty"`
	Required bool        `json:"required"`
	In       string      `json:"in"`
	Schema   ParamSchema `json:"schema"`
}

type ParamSchema struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type OpenAPI struct {
	OpenAPI    string                    `json:"openapi"`
	Info       Info                      `json:"info"`
	Paths      map[string]map[string]any `json:"paths"`
	Components map[string]map[string]any `json:"components"`
	offset     int
}

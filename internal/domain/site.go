package domain

// Site representa um site monitorado pela aplicação.
type Site struct {
	Name           string            `json:"name"`
	URL            string            `json:"url"`
	Method         string            `json:"method,omitempty"`
	TimeoutSeconds int               `json:"timeout_seconds,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
}

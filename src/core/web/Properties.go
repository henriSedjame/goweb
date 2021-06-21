package web


type CorsProperties struct {
	AllowedOrigins string `json:"allowed_origins"`
	AllowedHeaders string `json:"allowed_headers"`
	AllowedMethods string `json:"allowed_methods"`
}


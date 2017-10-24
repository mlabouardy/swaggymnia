package models

type ConfigSwagger struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Host        string `json:"host"`
	BasePath    string `json:"bastPath"`
	Schemes     string `json:"schemes"`
	Description string `json:"description"`
}

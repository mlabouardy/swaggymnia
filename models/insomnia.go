package models

import "time"

type Insomnia struct {
	Type         string     `json:"_type"`
	ExportFormat int        `json:"__export_format"`
	ExportDate   time.Time  `json:"__export_date"`
	ExportSource string     `json:"__export_source"`
	Resources    []Resource `json:"resources"`
}

package dto

type StatusOption struct {
	Value     int    `json:"value"`
	Label     string `json:"label"`
	Slug      string `json:"slug"`
	ClassName string `json:"class_name"`
}

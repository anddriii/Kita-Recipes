package dto

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Slug string `json:"slug"`
}

type CategoryResponseDetail struct {
	ID     int64                     `json:"id"`
	Name   string                    `json:"name"`
	Icon   string                    `json:"icon"`
	Slug   string                    `json:"slug"`
	Recipe []*RecipeResponseCategory `json:"recipes"`
}

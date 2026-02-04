package dto

type TemplateResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FollowTemplateRequest struct {
	TemplateID uint `json:"template_id"`
}

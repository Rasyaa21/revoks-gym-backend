package dto

type SettingsResponse struct {
	PushEnabled  bool   `json:"push_enabled"`
	EmailEnabled bool   `json:"email_enabled"`
	Language     string `json:"language"`
}

type UpdateSettingsRequest struct {
	PushEnabled  *bool  `json:"push_enabled"`
	EmailEnabled *bool  `json:"email_enabled"`
	Language     string `json:"language"`
}

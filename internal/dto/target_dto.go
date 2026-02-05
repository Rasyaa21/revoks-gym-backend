package dto

type TargetResponse struct {
	ID        uint   `json:"id"`
	Period    string `json:"period"`
	Title     string `json:"title"`
	GoalValue int    `json:"goal_value"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type CreateTargetRequest struct {
	Period    string `json:"period"` // weekly|monthly
	Title     string `json:"title"`
	GoalValue int    `json:"goal_value"`
	StartDate string `json:"start_date"` // RFC3339
	EndDate   string `json:"end_date"`   // RFC3339
}

type AddTargetProgressRequest struct {
	Value      int    `json:"value"`
	RecordedAt string `json:"recorded_at"` // RFC3339; optional
}

type TargetProgressResponse struct {
	ID         uint   `json:"id"`
	Value      int    `json:"value"`
	RecordedAt string `json:"recorded_at"`
}

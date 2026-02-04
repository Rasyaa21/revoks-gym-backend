package dto

type TargetResponse struct {
	ID        uint   `json:"id"`
	Period    string `json:"period"`
	Title     string `json:"title"`
	GoalValue int    `json:"goal_value"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
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

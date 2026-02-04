package dto

type CreateWorkoutProgressRequest struct {
	Title      string `json:"title"`
	Notes      string `json:"notes"`
	RecordedAt string `json:"recorded_at"` // RFC3339; optional
}

type WorkoutProgressResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Notes      string `json:"notes"`
	RecordedAt string `json:"recorded_at"`
}

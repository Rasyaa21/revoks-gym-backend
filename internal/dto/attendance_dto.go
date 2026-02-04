package dto

type AttendanceLogResponse struct {
	ID         uint   `json:"id"`
	Direction  string `json:"direction"`
	Source     string `json:"source"`
	OccurredAt string `json:"occurred_at"`
}

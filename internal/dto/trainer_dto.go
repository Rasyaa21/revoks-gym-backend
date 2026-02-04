package dto

type TrainerResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Specialty string `json:"specialty"`
	PhotoURL  string `json:"photo_url"`
}

type TrainerScheduleResponse struct {
	ID        uint   `json:"id"`
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
}

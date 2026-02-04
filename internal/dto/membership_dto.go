package dto

type MembershipStatusResponse struct {
	Status   string `json:"status"` // active|expired
	Plan     string `json:"plan,omitempty"`
	StartsAt string `json:"starts_at,omitempty"`
	EndsAt   string `json:"ends_at,omitempty"`
}

type MembershipHistoryItem struct {
	Status   string `json:"status"`
	Plan     string `json:"plan,omitempty"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type MembershipResponse struct {
	Current *MembershipStatusResponse `json:"current"`
	History []MembershipHistoryItem   `json:"history"`
}

type RenewMembershipRequest struct {
	Months int    `json:"months"`
	Plan   string `json:"plan"`
}

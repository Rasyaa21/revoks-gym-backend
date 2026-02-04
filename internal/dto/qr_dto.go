package dto

type MyQRCodeResponse struct {
	Token            string `json:"token"`
	ExpiresAt        string `json:"expires_at"`
	MembershipStatus string `json:"membership_status"` // active|expired
}

type ScanQRRequest struct {
	Token string `json:"token"`
}

type ScanQRResponse struct {
	Accepted   bool   `json:"accepted"`
	Reason     string `json:"reason,omitempty"`
	Direction  string `json:"direction,omitempty"` // in|out
	OccurredAt string `json:"occurred_at,omitempty"`
}

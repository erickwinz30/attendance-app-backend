package types

type UserAuthInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	IsHRD bool   `json:"is_hrd"`
}

type AuthCheckResponse struct {
	Authenticated bool          `json:"authenticated"`
	User          *UserAuthInfo `json:"user,omitempty"`
}

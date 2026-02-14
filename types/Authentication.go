package types

type UserAuthInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AuthCheckResponse struct {
	Authenticated bool          `json:"authenticated"`
	User          *UserAuthInfo `json:"user,omitempty"`
	IsAttended    bool          `json:"is_attended"`
}

package middleware

type sessionAuth struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

package middleware

const (
	ErrInvalidAuthorization = "Authorization header is missing or not in 'Bearer {token}' format"
	ErrInvalidToken         = "Invalid or expired token"
	ErrInvalidClaims        = "Invalid token claims"
	ErrInvalidIssueAt       = "Invalid issue at"
	ErrTokenExpired         = "Token is expired"
)

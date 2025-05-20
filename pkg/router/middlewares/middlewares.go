package middlewares

type contextKey string

const (
	UserEmailKey contextKey = "userEmail"
	UserAgentKey contextKey = "User-Agent"
)

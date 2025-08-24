package status

type StatusMessage = string

const (
	InternalServerError StatusMessage = "Internal Server Error"
	// msgConnectUserFail  StatusMessage = "Unable to connect user"
	// msgMalformedJSON    StatusMessage = "Invalid JSON received"
	MethodNotAllowed StatusMessage = "request method not allowed"
	// ParseFormFail       StatusMessage = "Failed to parse form"
	// UnauthorizedRequest StatusMessage = "User not authorized to make change"
	NoSessionCookie StatusMessage = "No session cookie found"
	InvalidSession  StatusMessage = "Invalid session"
)

package command

// Return codes to be used by command implementations and tests
const (
	OK                = 0
	BAD_REQUEST       = 400
	INTERNAL_ERROR    = 500
	DECRYPTION_DENIED = 403
)

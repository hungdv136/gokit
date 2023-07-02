package netkit

// Define common headers
const (
	HeaderAuthorization = "Authorization"
	HeaderRequestID     = "X-REQUEST-ID"
)

// Defines common token type
const (
	TokenTypeBearer = "Bearer"
)

// Define common verdicts for common cases
// Each service can define its own verdicts
const (
	VerdictMissingAuthentication = "missing_authentication"
	VerdictInvalidToken          = "invalid_token"
	VerdictInvalidParameters     = "invalid_parameters"
	VerdictSuccess               = "success"
	VerdictFailure               = "failure"
	VerdictNotFound              = "not_found"
	VerdictDuplicate             = "duplicate"
	VerdictLimitExceeded         = "limit_exceeded"
	VerdictExpiredData           = "expired_data"
	VerdictPermissionDenied      = "permission_denied"
)

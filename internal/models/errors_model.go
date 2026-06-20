package models

type ErrorBadRequest struct {
	Message string `json:"message" example:"Bad Request"`
	Error   string `json:"error" example:"PAYLOAD_NOT_VALID"`
}

type ErrorNotFound struct {
	Message string `json:"message" example:"Not Found"`
	Error   string `json:"error" example:"NOT_FOUND"`
}

type ErrorInternalServer struct {
	Message string `json:"message" example:"Internal Server Error"`
	Error   string `json:"error" example:"UNKNOWN_ERROR"`
}

type ErrorConflict struct {
	Message string `json:"message" example:"Conflict"`
	Error   string `json:"error" example:"EMAIL_ALREADY_EXISTS"`
}

type ErrorValidation struct {
	Message string            `json:"message" example:"Validation Failed"`
	Error   string            `json:"error" example:"VALIDATION_FAILED"`
	Fields  map[string]string `json:"fields" example:"{\"email\": \"invalid email format\"}"`
}

type ErrorUnauthorized struct {
	Message string `json:"message" example:"Unauthorized"`
	Error   string `json:"error" example:"INVALID_CREDENTIALS"`
}

type ErrorForbidden struct {
	Message string `json:"message" example:"Forbidden"`
	Error   string `json:"error" example:"ACCESS_DENIED"`
}

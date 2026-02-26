package httperr

type HTTPError struct {
	Status  int
	Code    string
	Message string
	Details any
	Cause   error
}

type APIErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

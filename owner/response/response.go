package response

type GenericResponse struct {
	ErrorCode int64  `json:"error_code"`
	Message   string `json:"message"`
}

package response

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Meta struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Message: message,
		Data:    data,
	}
}

func SuccessResponseWithMeta(message string, data interface{}, meta *Meta) Response {
	return Response{
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func ErrorResponse(Error string) Response {
	// Implementation for error response
	return Response{
		Message: "An error occurred",
		Error:   Error,
	}
}

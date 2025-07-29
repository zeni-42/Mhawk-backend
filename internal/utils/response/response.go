package response

import "net/http"

type ApiResponse struct {
	Status	int				`json:"status"`
	Success	bool			`json:"success"`
	Message	string			`json:"message"`
	Data	interface{}		`json:"data,omitempty"`
	Error	interface{}		`json:"error,omitempty"`
}

func Success(data interface{}, status int, message string) ApiResponse {
	if status == 0 {
		status = http.StatusOK
	}

	if message == "" {
		message = "Success"
	}

	res := ApiResponse {
		Status: status,
		Success: true,
		Message: message,
		Data: data,
		Error: nil,
	}

	return res
}

func Error(error interface{}, status int, message string) ApiResponse {
	if message == "" {
		message = "Somthing went wrong"
	}
	
	if status == 0 {
		status = http.StatusBadRequest
	}
	
	res := ApiResponse {
		Status: status,
		Success: false,
		Message: message,
		Data: nil,
		Error: error,
	}

	return res
}
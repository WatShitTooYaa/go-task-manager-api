package main

// type Response struct {
// 	Success bool   `json:"success"`
// 	Message string `json:"message"`
// 	Data    any    `json:"data,omitempty"`
// }

// func sendResponse(w http.ResponseWriter, message string, success bool, data any, statusCode int) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	res := Response{}
// 	res.Data = data
// 	res.Message = message
// 	// res.StatusCode = statusCode
// 	res.Success = success

// 	json.NewEncoder(w).Encode(&res)
// }

// func sendErrorResponse(w http.ResponseWriter, )

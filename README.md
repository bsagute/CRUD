package models

type Response struct {
	Message      string                 `json:"message"`
	Success      bool                   `json:"success"`
	ResponseCode int                    `json:"response_code"`
	Data         map[string]interface{} `json:"data"`
	Error        interface{}            `json:"error"`
}

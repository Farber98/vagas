package models

type MsgResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewMsgResponse(msg string) MsgResponse {
	return MsgResponse{
		Message: msg,
	}
}

func NewDataResponse(msg string, data interface{}) MsgResponse {
	return MsgResponse{
		Message: msg,
		Data:    data,
	}
}

package models

type MsgResponse struct {
	Message string `json:"message,omitempty"`
}

func NewMsgResponse(msg string) MsgResponse {
	return MsgResponse{
		Message: msg,
	}
}

package model

type ErrorMsg struct {
	StatusCode  int    `json:"status_code"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
}

func NewErrorMsg(statuscode int, status, desc string) *ErrorMsg {
	return &ErrorMsg{
		StatusCode:  statuscode,
		Status:      status,
		Description: desc,
	}
}

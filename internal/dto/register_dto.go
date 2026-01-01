package dto

type RegisterDTO struct {
	TerminalId string `json:"terminalId" binding:"required"`
	MsgType    string `json:"msgType" binding:"required"`
}

package web

type InputData struct {
	Board []int `json:"board" binding:"required,len=9"`
}

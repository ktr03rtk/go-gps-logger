package model

type BaseFilePath string

type Payload struct {
	Message   []byte
	FilePaths []BaseFilePath
}

func NewPayload() *Payload {
	return &Payload{}
}


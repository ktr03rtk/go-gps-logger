package model

type BaseFilePath string

type Payload struct {
	message   []byte
	filePaths []BaseFilePath
}

func NewPayload() *Payload {
	return &Payload{}
}


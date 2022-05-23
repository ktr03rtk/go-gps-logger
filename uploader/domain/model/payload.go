package model

const (
	PayloadSize       = 4096 * 1024
	MaxProcessFileNum = 100
)

type BaseFilePath string

type Payload struct {
	Message   []byte
	FilePaths []BaseFilePath
}

func NewPayload() *Payload {
	return &Payload{
		Message:   make([]byte, 0, PayloadSize),
		FilePaths: make([]BaseFilePath, 0, MaxProcessFileNum),
	}
}

func (p *Payload) Add(msg []byte, filePath BaseFilePath) {
	p.Message = append(p.Message, msg...)
	p.FilePaths = append(p.FilePaths, filePath)
}

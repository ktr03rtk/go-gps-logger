package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayloadAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialPayload *Payload
		inputMessage   []byte
		inputFilePath  BaseFilePath
		expectedOutput *Payload
	}{
		{
			"normal case1",
			&Payload{Message: []byte("initial\nmessage\n"), FilePaths: []BaseFilePath{"file1", "file2"}},
			[]byte("added\nmessage\n"),
			BaseFilePath("file3"),
			&Payload{Message: []byte("initial\nmessage\nadded\nmessage\n"), FilePaths: []BaseFilePath{"file1", "file2", "file3"}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			payload := tt.initialPayload
			payload.Add(tt.inputMessage, tt.inputFilePath)

			assert.Exactly(t, tt.expectedOutput, payload)
		})
	}
}

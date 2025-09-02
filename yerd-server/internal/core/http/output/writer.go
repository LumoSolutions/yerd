package output

import (
	"bufio"
	"encoding/json"
	"fmt"
)

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type StreamWriter struct {
	w *bufio.Writer
}

func NewStreamWriter(w *bufio.Writer) *StreamWriter {
	return &StreamWriter{w: w}
}

func (sw *StreamWriter) writeMessage(message, messageType string) error {
	data, err := json.Marshal(&Message{
		Type: messageType,
		Text: message,
	})

	if err != nil {
		return err
	}

	return sw.write(string(data))
}

func (sw *StreamWriter) WriteSuccess(message string, args ...interface{}) {
	sw.writeMessage(
		fmt.Sprintf(message, args...),
		"success",
	)
}

func (sw *StreamWriter) WriteInfo(message string, args ...interface{}) {
	sw.writeMessage(
		fmt.Sprintf(message, args...),
		"info",
	)
}

func (sw *StreamWriter) WriteError(message string, args ...interface{}) {
	sw.writeMessage(
		fmt.Sprintf(message, args...),
		"error",
	)
}

func (sw *StreamWriter) write(message string) error {
	if _, err := sw.w.WriteString(message + "\n"); err != nil {
		return err
	}

	return sw.w.Flush()
}

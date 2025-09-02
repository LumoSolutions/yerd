package output

import (
	"bufio"

	"github.com/gofiber/fiber/v2"
)

type StreamHandler func(writer *StreamWriter) error

// UseStream sets up streaming response and executes the handler with a StreamWriter
func UseStream(c *fiber.Ctx, handler StreamHandler) error {
	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		writer := NewStreamWriter(w)

		if err := handler(writer); err != nil {
			writer.WriteError("Error: " + err.Error())
		}
	})

	return nil
}

package client

import (
	"fmt"
	"io"
	"log"

	"github.com/lazzer64/ttsls/pkg/lsp/message"
)

type MessageWriter interface {
	Send(msg message.Message) error
	RequestFailed(id int, err error)
	ParseError(id int, err error)
	InvalidRequest(id int, err error)
	MethodNotFound(id int, err error)
	InvalidParams(id int, err error)
	InternalError(id int, err error)
	LogError(msg string)
	LogWarning(msg string)
	LogInfo(msg string)
	Log(msg string)
}

type messageWriter struct {
	io.Writer
}

func (w messageWriter) Send(msg message.Message) error {
	_, err := w.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(msg.Bytes()), msg.Bytes())))
	if err != nil {
		log.Printf("Could not write response: %s", err)
		return err
	}
	log.Printf("SEND %s\n", msg.Bytes())
	return nil
}

func (w messageWriter) error(code int, id int, err error) {
	log.Println(err)
	w.Send(
		message.NewResponseMessage(
			id,
			nil,
			map[string]interface{}{
				"code":    code,
				"message": err.Error(),
			},
		),
	)
}

func (w messageWriter) log(level message.ShowMessageLevel, msg string) error {
	return w.Send(
		message.NewShowMessageNotification(
			level,
			msg,
		),
	)
}

func (w messageWriter) RequestFailed(id int, err error)  { w.error(-32803, id, err) }
func (w messageWriter) ParseError(id int, err error)     { w.error(-32700, id, err) }
func (w messageWriter) InvalidRequest(id int, err error) { w.error(-32600, id, err) }
func (w messageWriter) MethodNotFound(id int, err error) { w.error(-32601, id, err) }
func (w messageWriter) InvalidParams(id int, err error)  { w.error(-32602, id, err) }
func (w messageWriter) InternalError(id int, err error)  { w.error(-32603, id, err) }

func (w messageWriter) LogError(msg string)   { w.log(message.MessageTypeError, msg) }
func (w messageWriter) LogWarning(msg string) { w.log(message.MessageTypeWarning, msg) }
func (w messageWriter) LogInfo(msg string)    { w.log(message.MessageTypeInfo, msg) }
func (w messageWriter) Log(msg string)        { w.log(message.MessageTypeLog, msg) }

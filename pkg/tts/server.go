package tts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

//go:generate stringer -type=TabletopSimulatorResponseMessageId
type TabletopSimulatorResponseMessageId int

const (
	UNKNOWN_RESPONSE        TabletopSimulatorResponseMessageId = -1
	NEW_OBJECT_RESPONSE     TabletopSimulatorResponseMessageId = 0
	GAME_LOADED_RESPONSE    TabletopSimulatorResponseMessageId = 1
	PRINT_RESPONSE          TabletopSimulatorResponseMessageId = 2
	ERROR_RESPONSE          TabletopSimulatorResponseMessageId = 3
	CUSTOM_RESPONSE         TabletopSimulatorResponseMessageId = 4
	RETURN_RESPONSE         TabletopSimulatorResponseMessageId = 5
	GAME_SAVED_RESPONSE     TabletopSimulatorResponseMessageId = 6
	OBJECT_CREATED_RESPONSE TabletopSimulatorResponseMessageId = 7
)

type TTSMessage struct {
	MessageID    TabletopSimulatorResponseMessageId `json:"messageID"`
	Guid         string                             `json:"guid,omitempty"`
	Name         string                             `json:"name,omitempty"`
	Error        string                             `json:"error,omitempty"`
	Message      string                             `json:"message,omitempty"`
	SavePath     string                             `json:"savePath,omitempty"`
	ScriptStates []ScriptState                      `json:"scriptStates,omitempty"`
}

type Handler func(TTSMessage)

type TTS struct {
	ctx      context.Context
	listener net.Listener
}

var global = TTS{}

func (tts TTS) Serve(handle Handler) error {
	for {
		select {
		default:
			conn, err := tts.listener.Accept()
			if err != nil {
				log.Printf("TTS  Error while accepting connection: %s\n", err)
				continue
			}

			b, err := io.ReadAll(conn)
			if err != nil {
				log.Printf("TTS  Error while reading connection: %s\n", err)
				conn.Close()
				continue
			}

			log.Printf("TTS  Received %q\n", string(b))
			msg := TTSMessage{}
			err = json.Unmarshal(b, &msg)
			if err != nil {
				log.Printf("TTS  Error while decoding message: %s\n", err)
				conn.Close()
				continue
			}

			conn.Close()
			handle(msg)
		case <-tts.ctx.Done():
			return tts.ctx.Err()
		}
	}
}

func Serve(ctx context.Context, handler Handler) error {
	if global.listener == nil {
		const clientServerAddr = "127.0.0.1:39998"

		lc := net.ListenConfig{}
		listener, err := lc.Listen(ctx, "tcp", clientServerAddr)

		if err != nil {
			listener.Close()
			return fmt.Errorf("Could not start server %w\n", err)
		}
		global.listener = listener

		log.Printf("TTS  Listening for connections at %s\n", clientServerAddr)
	}
	global.ctx = ctx
	return global.Serve(handler)
}

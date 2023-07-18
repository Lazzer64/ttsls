package lsp

import (
	"context"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/lazzer64/ttsls/pkg/lsp/client"
	"github.com/lazzer64/ttsls/pkg/lsp/handler"
	"github.com/lazzer64/ttsls/pkg/lsp/protocol"
	"github.com/lazzer64/ttsls/pkg/tts"
)

type LSP struct {
	handlers    map[string]handler.Handler
	nid         int
	mtx         sync.Mutex
	initialized bool
}

func NewLSP() *LSP {
	lsp := &LSP{
		handlers: map[string]handler.Handler{},
		nid:      1,
		mtx:      sync.Mutex{},
	}

	lsp.register("initialize", lsp.initializeHandler)
	lsp.register("initialized", lsp.initializedHandler)
	lsp.register("textDocument/didOpen", handler.TextDocumentDidOpenHandler)
	lsp.register("textDocument/didClose", handler.TextDocumentDidCloseHandler)
	lsp.register("textDocument/didChange", handler.TextDocumentDidChangeHandler)
	lsp.register("textDocument/definition", handler.TextDocumentDefinitionHandler)
	lsp.register("textDocument/hover", handler.TextDocumentHoverHandler)
	lsp.register("textDocument/completion", handler.TextDocumentCompletionHandler)
	lsp.register("textDocument/codeAction", handler.TextDocumentCodeActionHandler)
	lsp.register("textDocument/signatureHelp", handler.TextDocumentCodeActionHandler)
	lsp.register("workspace/executeCommand", handler.WorkspaceExecuteCommandHandler)

	return lsp
}

func (lsp *LSP) nextId() int {
	lsp.mtx.Lock()
	nid := lsp.nid
	lsp.nid++
	lsp.mtx.Unlock()
	return nid
}

func (lsp *LSP) register(m string, h handler.Handler) {
	if _, exists := lsp.handlers[m]; exists {
		log.Panicf("Attempted to re-register method %s.", m)
	}
	lsp.handlers[m] = h
}

func (lsp *LSP) Serve(ctx context.Context, r io.Reader, w io.Writer) {
	client := client.New(w)

	msgChan := make(chan protocol.Message)
	go func(r io.Reader) {
		for {
			msg, err := protocol.ReadMessage(r)
			if err != nil {
				log.Printf("LSP  Error reading message: %s", err)
				continue
			}
			msgChan <- msg
		}
	}(r)

	for {
		select {
		case msg := <-msgChan:
			log.Printf("RECV {%d %s}\n", msg.Id, msg.Method)
			if hand, ok := lsp.handlers[msg.Method]; ok {
				go hand(client, msg)
			} else if msg.Method != "" && !strings.HasPrefix(msg.Method, "$/") {
				log.Printf("LSP  Unhandled message: \"%s\"\n", msg.Method)
			}
		case <-ctx.Done():
			log.Panicf("LSP  Context cancelled: %s\n", ctx.Err())
			return
		}
	}
}

func (lsp *LSP) initializeHandler(c client.Client, u protocol.Message) {
	msg := u.Initialize()
	capabilities := map[string]protocol.ServerCapabilities{"capabilities": {
		PositionEncoding: protocol.PositionEncodingKindUTF16,
		TextDocumentSync: protocol.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    protocol.TextDocumentSyncKindFull,
		},
		DefinitionProvider: protocol.TextDocumentRegistrationOptions{
			DocumentSelector: []struct{ Pattern string }{{
				Pattern: "*.ttslua",
			}},
		},
		HoverProvider: protocol.TextDocumentRegistrationOptions{
			DocumentSelector: []struct{ Pattern string }{{
				Pattern: "*.ttslua",
			}},
		},
		ExecuteCommandProvider: protocol.ExecuteCommandOptions{
			Commands: []string{"tts.exec", "tts.getScripts", "tts.saveAndPlay"},
		},
		CompletionProvider: protocol.CompletionOptions{
			TriggerCharacters: []string{"."},
		},
		SignatureHelpProvider: protocol.SignatureHelpOptions{
			TriggerCharacters:   []string{},
			RetriggerCharacters: []string{},
		},
	}}
	c.Send(protocol.NewResponse(msg.Id, capabilities))
}

func (lsp *LSP) initializedHandler(client client.Client, u protocol.Message) {
	lsp.initialized = true
	log.Println("LSP  Client initialized")

	go tts.Serve(context.Background(), func(msg tts.TTSMessage) {
		lsp.ttsHandler(client, msg)
	})

	client.Log(protocol.MessageTypeInfo, "ttsls initialized")
}

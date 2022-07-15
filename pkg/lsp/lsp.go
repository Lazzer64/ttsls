package lsp

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lazzer64/ttsls/pkg/lsp/handler"
	"github.com/lazzer64/ttsls/pkg/lsp/message"
	"github.com/lazzer64/ttsls/pkg/lsp/message/types"
	"github.com/lazzer64/ttsls/pkg/lsp/client"
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

type includeManager struct {
	dir string
}

func (m includeManager) WriteInclude(name string, data []byte) error {
	path := filepath.Join(m.dir, name)
	if filepath.IsAbs(name) {
		path = name
	}
	err := os.MkdirAll(filepath.Dir(path), 0644)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (m includeManager) ReadInclude(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(m.dir, name))
}

func (lsp *LSP) Serve(ctx context.Context, r io.Reader, w io.Writer) {
	client := client.New(w)

	msgChan := make(chan message.UndefinedMessage)
	go func(r io.Reader) {
		for {
			msg, err := message.ReadMessage(r)
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
			log.Printf("RECV %s\n", msg.Bytes())
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

func (lsp *LSP) initializeHandler(client client.Client, u message.UndefinedMessage) {
	msg := u.Initialize()
	capabilities := map[string]types.ServerCapabilities{"capabilities": {
		PositionEncoding: types.PositionEncodingKindUTF16,
		TextDocumentSync: types.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    types.TextDocumentSyncKindFull,
		},
		DefinitionProvider: types.DefinitionOptions{
			DocumentSelector: []types.DocumentFilter{{
				Pattern: "*.ttslua",
			}},
		},
		HoverProvider: types.HoverOptions{
			DocumentSelector: []types.DocumentFilter{{
				Pattern: "*.ttslua",
			}},
		},
		ExecuteCommandProvider: types.ExecuteCommandOptions{
			Commands: []string{"tts.exec", "tts.getScripts", "tts.saveAndPlay"},
		},
		CompletionProvider: types.CompletionOptions{
			TriggerCharacters: []string{"."},
		},
		SignatureHelpProvider: types.SignatureHelpOptions{
			TriggerCharacters: []string{},
			RetriggerCharacters: []string{},
		},
	}}
	client.Send(message.NewResponseMessage(msg.Id, capabilities, nil))
}

func (lsp *LSP) initializedHandler(client client.Client, u message.UndefinedMessage) {
	lsp.initialized = true
	log.Println("LSP  Client initialized")

	go tts.Serve(context.Background(), func(msg tts.TTSMessage) {
		lsp.ttsHandler(client, msg)
	})

	client.Send(message.NewShowMessageNotification(
		message.MessageTypeInfo,
		"ttsls initialized",
	))
}

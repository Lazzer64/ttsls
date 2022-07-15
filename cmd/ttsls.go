package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/lazzer64/ttsls/pkg/lsp"
)

func main() {
	log.SetOutput(os.Stderr)

	home, _ := os.UserHomeDir()
	logfile := filepath.Join(home, ".ttsls", "log")
	os.MkdirAll(filepath.Dir(logfile), 0644)

	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stderr, f))
	}

	s := lsp.NewLSP()

	log.Println("CMD  Server started")
	s.Serve(context.Background(), os.Stdin, os.Stdout)
}

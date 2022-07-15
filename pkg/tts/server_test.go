package tts

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestTTS_Serve(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tts := TTS{ctx, newMockListener(ctx,
		fmt.Sprintf(`{"messageID":%d}`, GAME_LOADED_RESPONSE),
		`invalid...`,
		fmt.Sprintf(`{"messageID":%d}`, GAME_SAVED_RESPONSE),
	)}

	got := []TTSMessage{}
	tts.Serve(func(m TTSMessage) { got = append(got, m) })

	want := []TTSMessage{
		{MessageID: GAME_LOADED_RESPONSE},
		{MessageID: GAME_SAVED_RESPONSE},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func newMockListener(ctx context.Context, m ...string) net.Listener {
	conns := make(chan string, len(m))
	for _, s := range m {
		conns <- s
	}
	return mocklistener{ctx, conns}
}

type mocklistener struct {
	ctx   context.Context
	conns chan string
}

func (l mocklistener) Close() error   { return nil }
func (l mocklistener) Addr() net.Addr { return &net.UnixAddr{} }
func (l mocklistener) Accept() (net.Conn, error) {
	select {
	case <-l.ctx.Done():
		return nil, io.ErrClosedPipe
	case s := <-l.conns:
		return mockconn{bytes.NewBufferString(s)}, nil
	}
}

type mockconn struct{ *bytes.Buffer }

func (c mockconn) Close() error                       { return nil }
func (c mockconn) LocalAddr() net.Addr                { return &net.UnixAddr{} }
func (c mockconn) RemoteAddr() net.Addr               { return &net.UnixAddr{} }
func (c mockconn) SetDeadline(t time.Time) error      { return nil }
func (c mockconn) SetReadDeadline(t time.Time) error  { return nil }
func (c mockconn) SetWriteDeadline(t time.Time) error { return nil }

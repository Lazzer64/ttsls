package tts

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const ttsServerAddr string = "127.0.0.1:39999"

type TabletopSimulatorRequestMessageID int

const (
	UNKNOWN_REQUEST       TabletopSimulatorRequestMessageID = -1
	GET_SCRIPTS_REQUEST   TabletopSimulatorRequestMessageID = 0
	SAVE_AND_PLAY_REQUEST TabletopSimulatorRequestMessageID = 1
	CUSTOM_REQUEST        TabletopSimulatorRequestMessageID = 2
	EXEC_REQUEST          TabletopSimulatorRequestMessageID = 3
)

func send(s []byte) error {
	log.Printf("TTS  Connecting to %s\n", ttsServerAddr)
	conn, err := net.Dial("tcp", ttsServerAddr)
	if err != nil {
		log.Printf("TTS  Could not connect to Tabletop Simulator: %s\n", err)
		return err
	}
	defer conn.Close()
	log.Println("TTS  Connection established")

	log.Println(fmt.Sprintf("TTS  Sending: %s", s))

	conn.Write(s)
	if err != nil {
		log.Printf("TTS  Could not send message to Tabletop Simulator: %s", err)
		return err
	}
	return nil
}

func (tts TTS) Exec(s string) error {
	return send([]byte(fmt.Sprintf(`{"messageID":%d,"guid":"-1","script":%q}`, EXEC_REQUEST, s)))
}

func Exec(s string) error {
	return global.Exec(s)
}

func (tts TTS) GetScripts() error {
	return send([]byte(fmt.Sprintf(`{"messageID":%d}`, GET_SCRIPTS_REQUEST)))
}

func GetScripts() error {
	return global.GetScripts()
}

func (tts TTS) SaveAndPlay(scripts ...ScriptState) error {
	encoded, err := json.Marshal(scripts)
	if err != nil {
		return err
	}
	return send([]byte(fmt.Sprintf(`{"messageID":%d,"scriptStates":%s}`, SAVE_AND_PLAY_REQUEST, encoded)))
}

func SaveAndPlay(scripts ...ScriptState) error {
	return global.SaveAndPlay(scripts...)
}

package comm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gverger/godraw/models"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pull"
	"go.nanomsg.org/mangos/v3/protocol/push"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
)

func Listen(address string, stream chan<- models.Drawing) error {
	var sock mangos.Socket
	var err error
	var msg []byte
	if sock, err = pull.NewSocket(); err != nil {
		return fmt.Errorf("can't get new pull socket: %w", err)
	}
	if err = sock.Listen(address); err != nil {
		return fmt.Errorf("can't listen on pull socket: %w", err)
	}
	fmt.Printf("Listening to %s\n", address)
	for {
		// Could also use sock.RecvMsg to get header
		msg, err = sock.Recv()
		if err != nil {
			return fmt.Errorf("cannot receive from mangos Socket: %w", err)
		}
		if string(msg) == "STOP" {
			break
		}

		var drawing models.Drawing
		err := json.Unmarshal(msg, &drawing)
		if err != nil {
			fmt.Println("cannot unmarshal")
			continue
		}

		stream <- drawing
	}
	fmt.Println("STOPPING")
	return nil
}

type MsgSender struct {
	sock mangos.Socket
}

func NewMsgSender(address string) (MsgSender, error) {
	var err error
	sender := MsgSender{}

	if sender.sock, err = push.NewSocket(); err != nil {
		return sender, fmt.Errorf("can't get new push socket: %w", err)
	}
	if err = sender.sock.Dial(address); err != nil {
		return sender, fmt.Errorf("can't dial on push socket: %w", err)
	}

	return sender, nil
}

func (s MsgSender) Close() {
	s.sock.Close()
}

func (s MsgSender) Send(drawing models.Drawing) error {
	msg, err := json.Marshal(drawing)
	if err != nil {
		return fmt.Errorf("can't marshal drawing: %w", err)
	}

	if err := s.sock.Send(msg); err != nil {
		return fmt.Errorf("can't send message on push socket: %w", err)
	}

	time.Sleep(time.Second / 10)
	return nil
}

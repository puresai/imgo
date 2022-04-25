package ws

import (
	"context"
	"net/http"

	"github.com/13sai/imgo/internal/protocol"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // disable CheckOrigin
		},
	}
)

type wsConn struct {
	conn       *websocket.Conn
	activeTime int64
	writeQueue chan *protocol.Msg
}

func NewWsConn(conn *websocket.Conn) *wsConn {
	return &wsConn{conn: conn, writeQueue: make(chan *protocol.Msg, 100)}
}

func (w *wsConn) WriteLoop(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	w.writeQueue <- &protocol.Msg{Content: []byte("ac"), Head: protocol.Head{Type: 10}}
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-w.writeQueue:
			writer, err := w.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				w.handleError(err)
				return
			}
			err = protocol.Encode(writer, msg)
			if err != nil {
				w.handleError(err)
				return
			}
			writer.Close()
		}
	}
}

func (w *wsConn) ReceiveLoop(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgType, content, err := w.conn.NextReader()
			if err != nil {
				logrus.Errorf("ReceiveLoop err=%v", err)
				return
			}

			msg, err := protocol.Decode(content)
			if err != nil {
				w.handleError(err)
				return
			}

			w.handleMsg(msgType, msg)
		}
	}
}

func (w *wsConn) handleError(err error) {
	logrus.Errorf("Error err=%v", err)
}

func (w *wsConn) handleMsg(msgType int, msg *protocol.Msg) {
	switch msg.Type {
	default:
		logrus.Infof("msg content=%v", string(msg.Content))
	}
}

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/13sai/imgo/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	g := gin.New()
	g.Use(gin.Logger())
	g.GET("/ws", func(c *gin.Context) {
		conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logrus.Errorf("ws conn fail, err=%v", err)
			return
		}

		ctx, cancel := context.WithCancel(c)
		wsConn := ws.NewWsConn(conn)
		go wsConn.ReceiveLoop(ctx, cancel)
		go wsConn.WriteLoop(ctx, cancel)
	})

	srv := &http.Server{
		Handler: g,
		Addr:    ":1234",
	}

	go srv.ListenAndServe()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sig:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Errorf("Server Shutdown:%v", err)
		}
	}
}

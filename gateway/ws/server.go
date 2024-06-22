package ws

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

type Connections struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]bool
}

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWs(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("read error: ", err)
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))

		ws.Write([]byte("thank you for the message"))
	}
}

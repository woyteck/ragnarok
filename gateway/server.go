package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"woyteck.pl/ragnarok/gateway/api"
)

type handleFunc func(conn *websocket.Conn, req *api.TalkRequest) error

type Server struct {
	baseUri     string
	connections map[string]*websocket.Conn
	router      *mux.Router
	setup       *http.Server
	handler     handleFunc
}

func NewServer(baseUri string, handler handleFunc) *Server {
	router := mux.NewRouter()

	return &Server{
		baseUri:     baseUri,
		connections: make(map[string]*websocket.Conn),
		router:      router,
		setup: &http.Server{
			Handler:      router,
			Addr:         baseUri,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		handler: handler,
	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	fmt.Println("new connection from:", conn.RemoteAddr())

	s.connections[conn.RemoteAddr().String()] = conn

	for {
		var talkReq api.TalkRequest
		err := conn.ReadJSON(&talkReq)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("UnexpectedCloseError: %v\n", err)
			} else {
				log.Println("read err:", err)
			}
			break
		}
		if err == io.EOF {
			break
		}

		fmt.Printf("%+v\n", talkReq)
		err = s.handler(conn, &talkReq)
		if err != nil {
			log.Println("handler err:", err)
			break
		}
	}
}

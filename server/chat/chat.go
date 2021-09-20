package chat

import (
	"fmt"
	"gochat/server/authentication"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader     websocket.Upgrader
	connections  map[*websocket.Conn]string
	kickChannels map[*websocket.Conn]chan bool
	UserManager  authentication.UserManager
}

type IncomingMessagePacket struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type OutgoingMessagePacket struct {
	Type    string `json:"type"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type ErrorResponse struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

type LoginRequest struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success  bool   `json:"success"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

type OutgoingActivityPacket struct {
	Username string `json:"username"`
	Type     string `json:"type"`
}

func NewServer() Server {
	return Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
		connections:  map[*websocket.Conn]string{},
		UserManager:  authentication.NewUserManager(),
		kickChannels: map[*websocket.Conn]chan bool{},
	}
}

func (s *Server) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if _, has := s.connections[conn]; has {
			delete(s.connections, conn)
			if _, has := s.kickChannels[conn]; has {
				delete(s.kickChannels, conn)
			}
		}

		_ = conn.Close()
	}()

	for {
		alive := true

		type dataPacket struct {
			kick bool
			data map[string]interface{}
		}
		dataChan := make(chan dataPacket, 1)

		s.kickChannels[conn] = make(chan bool)
		go func() {
			kickChannel, has := s.kickChannels[conn]
			if has {
				if <-kickChannel && alive {
					dataChan <- dataPacket{true, nil}
				}
			}
		}()

		go func() {
			var p map[string]interface{}
			if err := conn.ReadJSON(&p); err != nil && alive {
				dataChan <- dataPacket{
					true, nil,
				}
			} else if alive {
				dataChan <- dataPacket{false, p}
			}
		}()

		p := <-dataChan
		if p.kick {
			alive = false
			s.BroadcastKick(conn)
			break
		}
		packet := p.data

		if !HasKeys(packet, "type") {
			continue
		}

		switch packet["type"] {
		case "message":
			if HasKeys(packet, "content") {
				if sender, ok := s.connections[conn]; ok {
					s.BroadcastMessage(sender, packet["content"].(string))
				} else {
					s.BroadcastLeave(conn)
					break
				}
			}
		case "login":
			if HasKeys(packet, "username", "password") {
				packet := LoginRequest{
					Type:     packet["type"].(string),
					Username: packet["username"].(string),
					Password: packet["password"].(string),
				}

				username, valid := s.UserManager.ValidateUser(packet.Username, packet.Password)

				response := LoginResponse{
					Type:     "login-response",
					Success:  valid,
					Username: username,
				}
				if err := conn.WriteJSON(response); err != nil {
					s.BroadcastLeave(conn)
					break
				}

				s.connections[conn] = username
				s.BroadcastJoin(conn)

				continue
			}
		}
	}
}

func (s *Server) BroadcastMessage(sender string, content string) {
	packet := OutgoingMessagePacket{
		Type:    "message",
		Sender:  sender,
		Content: content,
	}

	for conn := range s.connections {
		if err := conn.WriteJSON(packet); err != nil {
			delete(s.connections, conn)
		}
	}

	fmt.Printf("<%s> %s\n", sender, content)
}

func (s *Server) WriteError(conn *websocket.Conn, err error) {
	packet := ErrorResponse{
		Error: err.Error(),
		Type:  "error",
	}
	if err := conn.WriteJSON(packet); err != nil {
		delete(s.connections, conn)
	}
}

func (s *Server) BroadcastJoin(conn *websocket.Conn) {
	user, ok := s.connections[conn]
	if !ok {
		return
	}

	packet := OutgoingActivityPacket{
		Username: user,
		Type:     "join",
	}

	for conn := range s.connections {
		if err := conn.WriteJSON(packet); err != nil {
			delete(s.connections, conn)
		}
	}
}

func (s *Server) BroadcastLeave(conn *websocket.Conn) {
	user, ok := s.connections[conn]
	if !ok {
		return
	}

	packet := OutgoingActivityPacket{
		Username: user,
		Type:     "leave",
	}

	for conn := range s.connections {
		if err := conn.WriteJSON(packet); err != nil {
			delete(s.connections, conn)
		}
	}
}

func (s *Server) BroadcastKick(conn *websocket.Conn) {
	user, ok := s.connections[conn]
	if !ok {
		return
	}

	packet := OutgoingActivityPacket{
		Username: user,
		Type:     "kick",
	}

	for conn := range s.connections {
		if err := conn.WriteJSON(packet); err != nil {
			delete(s.connections, conn)
		}
	}
}

func (s *Server) GetConnection(username string) (*websocket.Conn, bool) {
	for k, v := range s.connections {
		if v == username {
			return k, true
		}
	}
	return nil, false
}

func (s *Server) KickUser(username string) bool {
	conn, has := s.GetConnection(username)
	if !has {
		return has
	}
	var ch chan bool
	ch, has = s.kickChannels[conn]
	if has {
		ch <- true
	}
	return has
}

func HasKeys(v map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		_, ok := v[key]
		if !ok {
			return false
		}
	}
	return true
}

func (s *Server) GetConnections() map[*websocket.Conn]string {
	return s.connections
}

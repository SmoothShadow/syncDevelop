package network

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Server WebSocket服务端
type Server struct {
	addr    string
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
	onMsg   func(msg *Message)
	onJoin  func(deviceName string)
	onLeave func(deviceName string)
	history *HistoryManager
}

func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		clients: make(map[string]*websocket.Conn),
		history: NewHistoryManager(100),
	}
}

func (s *Server) OnMessage(fn func(msg *Message)) {
	s.onMsg = fn
}

func (s *Server) OnJoin(fn func(deviceName string)) {
	s.onJoin = fn
}

func (s *Server) OnLeave(fn func(deviceName string)) {
	s.onLeave = fn
}

func (s *Server) History() *HistoryManager {
	return s.history
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWS)

	log.Printf("[Server] 监听 %s", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[Server] WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	s.mu.Lock()
	s.clients[remoteAddr] = conn
	s.mu.Unlock()

	log.Printf("[Server] 客户端连接: %s", remoteAddr)

	if s.onJoin != nil {
		s.onJoin(remoteAddr)
	}

	defer func() {
		s.mu.Lock()
		delete(s.clients, remoteAddr)
		s.mu.Unlock()

		if s.onLeave != nil {
			s.onLeave(remoteAddr)
		}
		log.Printf("[Server] 客户端断开: %s", remoteAddr)
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		msg, err := DecodeMessage(data)
		if err != nil {
			log.Printf("[Server] 消息解码失败: %v", err)
			continue
		}

		if s.onMsg != nil {
			s.onMsg(msg)
		}

		// 广播给其他客户端
		s.broadcast(data, conn)
	}
}

func (s *Server) broadcast(data []byte, exclude *websocket.Conn) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		if c != exclude {
			go func(conn *websocket.Conn) {
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("[Server] 广播失败: %v", err)
				}
			}(c)
		}
	}
}

// SendToAll 向所有客户端发送消息
func (s *Server) SendToAll(msg *Message) {
	data, err := EncodeMessage(msg)
	if err != nil {
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		go func(conn *websocket.Conn) {
			conn.WriteMessage(websocket.TextMessage, data)
		}(c)
	}
}

// GetClientCount 获取连接客户端数
func (s *Server) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

// GetLocalIP 获取本机局域网IP
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

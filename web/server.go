package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ---------- 数据结构 ----------

type ItemType string

const (
	ItemText  ItemType = "text"
	ItemFile  ItemType = "file"
	ItemImage ItemType = "image"
)

type SyncItem struct {
	ID        string    `json:"id"`
	Type      ItemType  `json:"type"`
	Content   string    `json:"content"`
	FileName  string    `json:"fileName,omitempty"`
	FileSize  int64     `json:"fileSize,omitempty"`
	From      string    `json:"from"`
	Timestamp int64     `json:"timestamp"`
}

// ---------- Server ----------

type Server struct {
	addr      string
	fileDir   string
	items     []SyncItem
	mu        sync.RWMutex
	wsClients map[string]*websocket.Conn
	wsMu      sync.RWMutex
	maxItems  int
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewServer(addr, fileDir string) *Server {
	return &Server{
		addr:      addr,
		fileDir:   fileDir,
		items:     make([]SyncItem, 0),
		wsClients: make(map[string]*websocket.Conn),
		maxItems:  200,
	}
}

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

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/items", s.handleItems)
	mux.HandleFunc("/api/send", s.handleSend)
	mux.HandleFunc("/api/upload", s.handleUpload)
	mux.HandleFunc("/api/download/", s.handleDownload)
	mux.HandleFunc("/api/delete/", s.handleDelete)
	mux.HandleFunc("/api/clear", s.handleClear)
	mux.HandleFunc("/ws", s.handleWS)

	log.Printf("[Server] 监听 %s", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

// ---------- 前端 ----------

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

// ---------- API ----------

func (s *Server) handleItems(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	items := make([]SyncItem, len(s.items))
	copy(items, s.items)
	s.mu.RUnlock()

	sort.Slice(items, func(i, j int) bool {
		return items[i].Timestamp > items[j].Timestamp
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 405)
		return
	}

	var body struct {
		Text string `json:"text"`
		From string `json:"from"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if strings.TrimSpace(body.Text) == "" {
		http.Error(w, "empty text", 400)
		return
	}

	item := SyncItem{
		ID:        fmt.Sprintf("%d", time.Now().UnixMilli()),
		Type:      ItemText,
		Content:   body.Text,
		From:      body.From,
		Timestamp: time.Now().UnixMilli(),
	}
	s.addItem(item)
	s.broadcastNewItem(item)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "id": item.ID})
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 405)
		return
	}

	r.ParseMultipartForm(100 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	saveName := fmt.Sprintf("%d_%s", time.Now().UnixMilli(), header.Filename)
	savePath := filepath.Join(s.fileDir, saveName)
	os.MkdirAll(filepath.Dir(savePath), 0755)

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	from := r.FormValue("from")
	item := SyncItem{
		ID:        fmt.Sprintf("%d", time.Now().UnixMilli()),
		Type:      ItemFile,
		Content:   savePath,
		FileName:  header.Filename,
		FileSize:  written,
		From:      from,
		Timestamp: time.Now().UnixMilli(),
	}
	s.addItem(item)
	s.broadcastNewItem(item)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "id": item.ID, "fileName": header.Filename})
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/download/")

	s.mu.RLock()
	var filePath, fileName string
	for _, item := range s.items {
		if item.ID == id && item.Type == ItemFile {
			filePath = item.Content
			fileName = item.FileName
			break
		}
	}
	s.mu.RUnlock()

	if filePath == "" {
		http.Error(w, "not found", 404)
		return
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "file missing", 404)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	http.ServeFile(w, r, filePath)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "method not allowed", 405)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/delete/")

	s.mu.Lock()
	for i, item := range s.items {
		if item.ID == id {
			if item.Type == ItemFile {
				os.Remove(item.Content)
			}
			s.items = append(s.items[:i], s.items[i+1:]...)
			break
		}
	}
	s.mu.Unlock()

	s.broadcastEvent("update")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "method not allowed", 405)
		return
	}

	s.mu.Lock()
	for _, item := range s.items {
		if item.Type == ItemFile {
			os.Remove(item.Content)
		}
	}
	s.items = s.items[:0]
	s.mu.Unlock()

	s.broadcastEvent("update")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ---------- WebSocket ----------

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	clientID := fmt.Sprintf("%d", time.Now().UnixNano())

	s.wsMu.Lock()
	s.wsClients[clientID] = conn
	s.wsMu.Unlock()

	defer func() {
		s.wsMu.Lock()
		delete(s.wsClients, clientID)
		s.wsMu.Unlock()
	}()

	// 保持连接
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (s *Server) broadcastNewItem(item SyncItem) {
	data, _ := json.Marshal(map[string]interface{}{
		"event": "new_item",
		"item":  item,
	})
	s.wsBroadcast(data)
}

func (s *Server) broadcastEvent(event string) {
	data, _ := json.Marshal(map[string]string{"event": event})
	s.wsBroadcast(data)
}

func (s *Server) wsBroadcast(data []byte) {
	s.wsMu.RLock()
	defer s.wsMu.RUnlock()
	for _, conn := range s.wsClients {
		go func(c *websocket.Conn) {
			c.WriteMessage(websocket.TextMessage, data)
		}(conn)
	}
}

// ---------- 内部 ----------

func (s *Server) addItem(item SyncItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
	if len(s.items) > s.maxItems {
		removed := s.items[0]
		if removed.Type == ItemFile {
			os.Remove(removed.Content)
		}
		s.items = s.items[1:]
	}
}

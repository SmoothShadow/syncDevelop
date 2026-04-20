package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"syncDevelop/network"
	"syncDevelop/sync"
	"time"
)

//go:embed files/*
var staticFiles embed.FS

// WebUI Web管理界面
type WebUI struct {
	engine *sync.SyncEngine
	addr   string
}

func NewWebUI(engine *sync.SyncEngine, addr string) *WebUI {
	return &WebUI{
		engine: engine,
		addr:   addr,
	}
}

func (w *WebUI) Start() error {
	mux := http.NewServeMux()

	// 静态文件
	staticFS, err := fs.Sub(staticFiles, "files")
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API: 获取历史记录
	mux.HandleFunc("/api/history", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		entries := w.engine.History().GetAll()
		json.NewEncoder(res).Encode(entries)
	})

	// API: 手动发送文本
	mux.HandleFunc("/api/send", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(res, "method not allowed", 405)
			return
		}
		var body struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(res, err.Error(), 400)
			return
		}
		w.engine.SendText(body.Text)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]string{"status": "ok"})
	})

	// API: 设备状态
	mux.HandleFunc("/api/status", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]interface{}{
			"isServer":     w.engine.IsServer(),
			"clientCount":  w.engine.GetClientCount(),
			"localIP":      network.GetLocalIP(),
			"timestamp":    time.Now().UnixMilli(),
		})
	})

	// API: 文件上传
	mux.HandleFunc("/api/upload", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(res, "method not allowed", 405)
			return
		}
		file, header, err := req.FormFile("file")
		if err != nil {
			http.Error(res, err.Error(), 400)
			return
		}
		defer file.Close()

		// 构造文件消息
		fileInfo := network.FileInfo{
			Name: header.Filename,
			Size: header.Size,
		}
		infoJSON, _ := json.Marshal(fileInfo)
		msg := network.NewMessage(network.MsgFile, "webui", string(infoJSON))
		msg.Extra = header.Filename
		w.engine.History().Add(network.ClipEntry{
			ID:        msg.ID,
			Type:      network.MsgFile,
			Content:   string(infoJSON),
			FileName:  header.Filename,
			FileSize:  header.Size,
			From:      "webui",
			Timestamp: msg.Timestamp,
		})

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]string{
			"status":   "ok",
			"fileName": header.Filename,
			"size":     fmt.Sprintf("%d", header.Size),
		})
	})

	log.Printf("[WebUI] 管理界面监听 %s", w.addr)
	return http.ListenAndServe(w.addr, mux)
}

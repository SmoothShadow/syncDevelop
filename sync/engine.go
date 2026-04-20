package sync

import (
	"encoding/base64"
	"fmt"
	"log"
	"syncDevelop/clipboard"
	"syncDevelop/network"
)

// SyncEngine 同步引擎
type SyncEngine struct {
	deviceName string
	server     *network.Server
	client     *network.Client
	watcher    *clipboard.Watcher
	history    *network.HistoryManager
	isServer   bool
	fileDir    string
}

func NewSyncEngine(deviceName string, isServer bool, addr string, fileDir string) *SyncEngine {
	e := &SyncEngine{
		deviceName: deviceName,
		isServer:   isServer,
		history:    network.NewHistoryManager(100),
		fileDir:    fileDir,
	}

	if isServer {
		e.server = network.NewServer(addr)
		e.server.OnMessage(func(msg *network.Message) {
			e.onRemoteMessage(msg)
		})
	} else {
		e.client = network.NewClient(addr)
		e.client.OnMessage(func(msg *network.Message) {
			e.onRemoteMessage(msg)
		})
	}

	e.watcher = clipboard.NewWatcher(0)
	e.watcher.OnChange(e.onLocalClipboardChange)

	return e
}

func (e *SyncEngine) Start() error {
	e.watcher.Start()

	if e.isServer {
		go func() {
			if err := e.server.Start(); err != nil {
				log.Fatalf("[SyncEngine] 服务端启动失败: %v", err)
			}
		}()
	} else {
		if err := e.client.Connect(); err != nil {
			log.Printf("[SyncEngine] 连接失败: %v，将自动重连", err)
			e.client.AutoReconnect(3)
		}
	}
	return nil
}

func (e *SyncEngine) Stop() {
	e.watcher.Stop()
	if e.client != nil {
		e.client.Close()
	}
}

// onLocalClipboardChange 本地剪贴板变化 → 推送到远端
func (e *SyncEngine) onLocalClipboardChange(content clipboard.Content) {
	switch content.Type {
	case clipboard.TypeText:
		if content.Text == "" {
			return
		}
		msg := network.NewMessage(network.MsgClipText, e.deviceName, content.Text)
		e.history.Add(network.ClipEntry{
			ID:        msg.ID,
			Type:      msg.Type,
			Content:   content.Text,
			From:      e.deviceName,
			Timestamp: msg.Timestamp,
		})
		e.send(msg)
		log.Printf("[SyncEngine] 推送文本 (%d 字符)", len(content.Text))

	case clipboard.TypeImage:
		if len(content.Image) == 0 {
			return
		}
		encoded := base64.StdEncoding.EncodeToString(content.Image)
		msg := network.NewMessage(network.MsgClipImage, e.deviceName, encoded)
		e.history.Add(network.ClipEntry{
			ID:        msg.ID,
			Type:      msg.Type,
			Content:   "[图片 " + formatSize(len(content.Image)) + "]",
			From:      e.deviceName,
			Timestamp: msg.Timestamp,
		})
		e.send(msg)
		log.Printf("[SyncEngine] 推送图片 (%s)", formatSize(len(content.Image)))
	}
}

// onRemoteMessage 远端消息 → 写入本地剪贴板
func (e *SyncEngine) onRemoteMessage(msg *network.Message) {
	if msg.From == e.deviceName {
		return
	}

	switch msg.Type {
	case network.MsgClipText:
		e.history.Add(network.ClipEntry{
			ID:        msg.ID,
			Type:      msg.Type,
			Content:   msg.Payload,
			From:      msg.From,
			Timestamp: msg.Timestamp,
		})
		if err := clipboard.WriteText(msg.Payload); err != nil {
			log.Printf("[SyncEngine] 写入剪贴板失败: %v", err)
		} else {
			log.Printf("[SyncEngine] 收到文本并写入剪贴板 (%d 字符)", len(msg.Payload))
		}

	case network.MsgClipImage:
		data, err := base64.StdEncoding.DecodeString(msg.Payload)
		if err != nil {
			log.Printf("[SyncEngine] 图片解码失败: %v", err)
			return
		}
		e.history.Add(network.ClipEntry{
			ID:        msg.ID,
			Type:      msg.Type,
			Content:   "[图片 " + formatSize(len(data)) + "]",
			From:      msg.From,
			Timestamp: msg.Timestamp,
		})
		if err := clipboard.WriteImage(data); err != nil {
			log.Printf("[SyncEngine] 写入图片剪贴板失败: %v", err)
		} else {
			log.Printf("[SyncEngine] 收到图片并写入剪贴板 (%s)", formatSize(len(data)))
		}

	case network.MsgFile:
		e.history.Add(network.ClipEntry{
			ID:        msg.ID,
			Type:      msg.Type,
			Content:   msg.Payload,
			FileName:  msg.Extra,
			From:      msg.From,
			Timestamp: msg.Timestamp,
		})
		log.Printf("[SyncEngine] 收到文件: %s", msg.Extra)
	}
}

func (e *SyncEngine) send(msg *network.Message) {
	if e.isServer && e.server != nil {
		e.server.SendToAll(msg)
	} else if e.client != nil && e.client.IsConnected() {
		e.client.Send(msg)
	}
}

func (e *SyncEngine) History() *network.HistoryManager {
	return e.history
}

// SendText 手动推送文本
func (e *SyncEngine) SendText(text string) {
	msg := network.NewMessage(network.MsgClipText, e.deviceName, text)
	e.history.Add(network.ClipEntry{
		ID:        msg.ID,
		Type:      msg.Type,
		Content:   text,
		From:      e.deviceName,
		Timestamp: msg.Timestamp,
	})
	e.send(msg)
}

func (e *SyncEngine) IsServer() bool {
	return e.isServer
}

func (e *SyncEngine) GetClientCount() int {
	if e.server != nil {
		return e.server.GetClientCount()
	}
	return 0
}

func formatSize(bytes int) string {
	const (
		KB = 1024
		MB = 1024 * KB
	)
	switch {
	case bytes >= MB:
		return fmt.Sprintf("%dMB", bytes/MB)
	case bytes >= KB:
		return fmt.Sprintf("%dKB", bytes/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

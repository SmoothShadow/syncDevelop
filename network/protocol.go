package network

import (
	"encoding/json"
	"sync"
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	MsgClipText  MessageType = "clip_text"  // 文本剪贴板
	MsgClipImage MessageType = "clip_image" // 图片剪贴板
	MsgFile      MessageType = "file"       // 文件同步
	MsgFileData  MessageType = "file_data"  // 文件二进制数据
	MsgFileAck   MessageType = "file_ack"   // 文件接收确认
	MsgHeartbeat MessageType = "heartbeat"  // 心跳
	MsgHistory   MessageType = "history"    // 历史记录请求
	MsgJoin      MessageType = "join"       // 设备加入
	MsgLeave     MessageType = "leave"      // 设备离开
)

// Message 同步消息
type Message struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	From      string      `json:"from"`       // 发送设备名
	Timestamp int64       `json:"timestamp"`   // 毫秒时间戳
	Payload   string      `json:"payload"`     // 文本内容/base64图片/文件元信息
	Extra     string      `json:"extra,omitempty"` // 扩展字段(文件名等)
}

// FileInfo 文件元信息
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Hash string `json:"hash,omitempty"`
}

// ClipEntry 剪贴板历史条目
type ClipEntry struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	Content   string      `json:"content"`   // 文本内容或缩略图base64
	FileName  string      `json:"fileName,omitempty"`
	FileSize  int64       `json:"fileSize,omitempty"`
	From      string      `json:"from"`
	Timestamp int64       `json:"timestamp"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	Name      string `json:"name"`
	OS        string `json:"os"`
	IP        string `json:"ip"`
	Connected bool   `json:"connected"`
	JoinedAt  int64  `json:"joinedAt"`
}

// HistoryManager 历史记录管理
type HistoryManager struct {
	mu      sync.RWMutex
	entries []ClipEntry
	maxLen  int
}

func NewHistoryManager(maxLen int) *HistoryManager {
	if maxLen <= 0 {
		maxLen = 100
	}
	return &HistoryManager{
		entries: make([]ClipEntry, 0, maxLen),
		maxLen:  maxLen,
	}
}

func (h *HistoryManager) Add(entry ClipEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, entry)
	if len(h.entries) > h.maxLen {
		h.entries = h.entries[len(h.entries)-h.maxLen:]
	}
}

func (h *HistoryManager) GetAll() []ClipEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result := make([]ClipEntry, len(h.entries))
	copy(result, h.entries)
	return result
}

func (h *HistoryManager) GetSince(since int64) []ClipEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var result []ClipEntry
	for _, e := range h.entries {
		if e.Timestamp > since {
			result = append(result, e)
		}
	}
	return result
}

// DecodeMessage 从JSON解码消息
func DecodeMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// EncodeMessage 将消息编码为JSON
func EncodeMessage(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}

// NewMessage 创建新消息
func NewMessage(msgType MessageType, from, payload string) *Message {
	return &Message{
		ID:        generateID(),
		Type:      msgType,
		From:      from,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
}

func generateID() string {
	return time.Now().Format("20060102150405.000")
}

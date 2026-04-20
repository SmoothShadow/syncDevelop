package clipboard

import (
	"log"
	"sync"
	"time"
)

// ContentType 剪贴板内容类型
type ContentType int

const (
	TypeText  ContentType = iota
	TypeImage
)

// Content 剪贴板内容
type Content struct {
	Type ContentType
	Text string
	Image []byte
}

// OnChangeFunc 剪贴板变化回调
type OnChangeFunc func(content Content)

// Watcher 剪贴板监听器
type Watcher struct {
	interval  time.Duration
	lastText  string
	lastImageHash string
	onChange   OnChangeFunc
	stopCh    chan struct{}
	mu        sync.Mutex
	running   bool
}

func NewWatcher(interval time.Duration) *Watcher {
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	return &Watcher{
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (w *Watcher) OnChange(fn OnChangeFunc) {
	w.onChange = fn
}

func (w *Watcher) Start() {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.mu.Unlock()

	go w.watchLoop()
	log.Println("[Clipboard] 监听已启动")
}

func (w *Watcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.running {
		return
	}
	w.running = false
	close(w.stopCh)
	log.Println("[Clipboard] 监听已停止")
}

func (w *Watcher) watchLoop() {
	// 初始化当前剪贴板内容
	w.lastText, _ = readText()

	for {
		select {
		case <-w.stopCh:
			return
		default:
		}

		// 检查文本变化
		text, err := readText()
		if err == nil && text != w.lastText && text != "" {
			w.lastText = text
			if w.onChange != nil {
				w.onChange(Content{
					Type: TypeText,
					Text: text,
				})
			}
		}

		// 检查图片变化
		img, err := readImage()
		if err == nil && len(img) > 0 {
			hash := simpleHash(img)
			if hash != w.lastImageHash {
				w.lastImageHash = hash
				if w.onChange != nil {
					w.onChange(Content{
						Type:  TypeImage,
						Image: img,
					})
				}
			}
		}

		time.Sleep(w.interval)
	}
}

// WriteText 写入文本到剪贴板
func WriteText(text string) error {
	return writeText(text)
}

// WriteImage 写入图片到剪贴板
func WriteImage(data []byte) error {
	return writeImage(data)
}

// ReadText 读取剪贴板文本
func ReadText() (string, error) {
	return readText()
}

// simpleHash 简单哈希，用于检测图片变化
func simpleHash(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	// 用前16字节+后16字节+长度做快速指纹
	n := len(data)
	sample := make([]byte, 0, 34)
	if n >= 16 {
		sample = append(sample, data[:16]...)
		sample = append(sample, data[n-16:]...)
	} else {
		sample = append(sample, data...)
	}
	// 返回长度+采样字节的字符串表示
	return string(rune(n)) + string(rune(n>>8)) + string(sample)
}

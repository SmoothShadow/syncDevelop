package network

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	serverAddr string
	conn       *websocket.Conn
	onMsg      func(msg *Message)
	onConnect  func()
	onClose    func()
	stopCh     chan struct{}
}

func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
		stopCh:     make(chan struct{}),
	}
}

func (c *Client) OnMessage(fn func(msg *Message)) {
	c.onMsg = fn
}

func (c *Client) OnConnect(fn func()) {
	c.onConnect = fn
}

func (c *Client) OnClose(fn func()) {
	c.onClose = fn
}

func (c *Client) Connect() error {
	url := "ws://" + c.serverAddr + "/ws"
	log.Printf("[Client] 连接服务端: %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	if c.onConnect != nil {
		c.onConnect()
	}

	go c.readLoop()
	return nil
}

func (c *Client) readLoop() {
	defer func() {
		if c.onClose != nil {
			c.onClose()
		}
	}()

	for {
		select {
		case <-c.stopCh:
			return
		default:
		}

		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[Client] 读取失败: %v", err)
			c.conn = nil
			return
		}

		msg, err := DecodeMessage(data)
		if err != nil {
			continue
		}

		if c.onMsg != nil {
			c.onMsg(msg)
		}
	}
}

// Send 发送消息
func (c *Client) Send(msg *Message) error {
	if c.conn == nil {
		return nil
	}
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Close 关闭连接
func (c *Client) Close() {
	close(c.stopCh)
	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
	}
}

// AutoReconnect 自动重连
func (c *Client) AutoReconnect(interval time.Duration) {
	go func() {
		for {
			select {
			case <-c.stopCh:
				return
			default:
			}

			if c.conn == nil {
				if err := c.Connect(); err != nil {
					log.Printf("[Client] 重连失败: %v, %v后重试", err, interval)
					time.Sleep(interval)
					continue
				}
				log.Printf("[Client] 重连成功")
			}
			time.Sleep(interval)
		}
	}()
}

// IsConnected 是否已连接
func (c *Client) IsConnected() bool {
	return c.conn != nil
}

package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	idleMu sync.Mutex

	Uid string

	messageMu      sync.Mutex
	readMessage    []*Message
	readMessageSeq map[string]*Message

	message chan *Message

	*websocket.Conn
	s                 *Server
	lastCommunication time.Time
	MaxIdleTime       time.Duration
	done              chan struct{}
}

func NewConnection(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade fail, err: %v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),
		message:           make(chan *Message, 1),
		lastCommunication: time.Now(),
		MaxIdleTime:       s.serverOption.MaxIdleTime,
		done:              make(chan struct{}),
	}

	go conn.keepAlive()
	return conn
}

func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()

	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 已经有消息的记录，该消息已经有ack的确认
		if len(c.readMessage) == 0 {
			// 队列中没有消息
			return
		}
		if m.AckSeq >= msg.AckSeq {
			return
		}
		c.readMessageSeq[msg.Id] = msg
		return
	}

	if msg.FrameType == FrameAck {
		return
	}

	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg

}

func (c *Conn) ReadMessage() (int, []byte, error) {
	messageType, bytes, err := c.Conn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.lastCommunication = time.Time{} // 读取到客户端的消息，开始处理消息内容，连接不空闲了
	return messageType, bytes, err
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {

	c.idleMu.Lock()
	defer c.idleMu.Unlock()

	err := c.Conn.WriteMessage(messageType, data)
	c.lastCommunication = time.Now() // 发送给客户端消息后，处理结束，连接要空闲了
	return err
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return c.Conn.Close()
}

func (c *Conn) keepAlive() {
	timer := time.NewTimer(c.MaxIdleTime)
	defer func() {
		timer.Stop()
		return
	}()

	for {
		select {
		case <-timer.C:
			c.idleMu.Lock()
			lastCommunication := c.lastCommunication
			if lastCommunication.IsZero() { // 一直在通信
				c.idleMu.Unlock()
				timer.Reset(c.MaxIdleTime)
				return
			}

			idleTime := c.MaxIdleTime - time.Since(lastCommunication)
			c.idleMu.Unlock()
			if idleTime <= 0 {
				c.s.Close(c)
				return
			}
			timer.Reset(idleTime)
		case <-c.done:
			return
		}

	}
}

package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

type Client interface {
	Close() error

	Send(v interface{}) error
	SendUid(v interface{}, uids ...string) error
	Read(v interface{}) error
}

type client struct {
	*websocket.Conn
	host string

	opt DialOption

	Discover
}

func NewClient(host string, opt ...DialOptions) *client {
	dialOption := NewDialOption(opt...)

	c := &client{
		Conn: nil,
		host: host,
		opt:  dialOption,
	}

	conn, err := c.Dial()
	if err != nil {
		panic(err)
	}

	c.Conn = conn

	return c
}

func (c *client) Dial() (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   c.host,
		Path:   c.opt.Pattern,
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

func (c *client) Send(v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = c.Conn.WriteMessage(websocket.TextMessage, bytes)
	if err == nil {
		return nil
	}

	// 发送失败，重连发送
	conn, err := c.Dial()
	if err != nil {
		panic(err)
	}
	c.Conn = conn
	return c.WriteMessage(websocket.TextMessage, bytes)
}

func (c *client) SendUid(v interface{}, uids ...string) error {
	if c.Discover != nil {
		return c.Discover.Transpond(v, uids...)
	}
	return c.Send(v)
}

func (c *client) Read(v interface{}) error {
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg, v)
}

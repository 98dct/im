package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

type Server struct {
	routes map[string]HandlerFunc
	addr   string

	serverOption serverOption

	mu         sync.RWMutex
	connToUser map[*Conn]string
	userToConn map[string]*Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOption) *Server {
	serverOption := newServerOption(opts...)
	return &Server{
		routes: make(map[string]HandlerFunc),
		addr:   addr,
		//authentication: new(authentication),
		serverOption: serverOption,
		connToUser:   make(map[*Conn]string),
		userToConn:   make(map[string]*Conn),
		upgrader:     websocket.Upgrader{},
		Logger:       logx.WithContext(context.Background()),
	}
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			s.Errorf("ws recover err: %v", err)
		}
	}()

	conn := NewConnection(s, w, r)
	if conn == nil {
		return
	}

	// 鉴权
	if !s.serverOption.Authentication.Auth(w, r) {
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不具备访问权限！")}, conn)
		s.Close(conn)
		return
	}

	// 记录链接
	s.addConn(conn, r)

	// 根据连接对象，获取请求信息
	go s.handlerConn(conn)
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {

	uid := s.GetUser(conn)
	conn.Uid = uid
	// 循环读取连接上的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("read msg err : %v", err)
			// todo 关闭连接
			s.Close(conn)
			return
		}

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			s.Errorf("json unmarshal err : %v, msg: %v", err, msg)
			// todo 关闭连接
			s.Close(conn)
			return
		}

		// 根据消息类型处理
		switch message.FrameType {
		case FramePing:
			s.Send(&Message{FrameType: FramePing}, conn)
		case FrameData:
			// 根据请求的方法，查找路由，执行具体的handler
			handlerFunc, ok := s.routes[message.Method]
			if !ok {
				s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在目标方法，%v", message.Method)}, conn)
				continue
			}
			handlerFunc(s, conn, &message)
		}

	}
}

func (s *Server) addConn(conn *Conn, req *http.Request) {
	userId := s.serverOption.Authentication.UserId(req)
	s.mu.Lock()
	defer s.mu.Unlock()

	// 验证用户是否登陆过
	if c := s.userToConn[userId]; c != nil {
		// 如果有，关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = userId
	s.userToConn[userId] = conn
}

func (s *Server) GetConn(uId string) *Conn {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userToConn[uId]
}

func (s *Server) GetConns(uIds ...string) []*Conn {
	if len(uIds) == 0 {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	conns := make([]*Conn, 0, len(uIds))
	for _, uId := range uIds {
		conns = append(conns, s.userToConn[uId])
	}

	return conns
}

func (s *Server) GetUser(conn *Conn) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connToUser[conn]
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {

	s.mu.Lock()
	defer s.mu.Unlock()

	uid, ok := s.connToUser[conn]
	if !ok || uid == "" {
		return
	}
	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, uIds ...string) error {

	if len(uIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(uIds...)...)

}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	bytes, _ := json.Marshal(msg)
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
			return err
		}
	}
	return nil

}

func (s *Server) AddRoutes(routes []Router) {
	for _, route := range routes {
		s.routes[route.Method] = route.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.serverOption.pattern, s.ServeWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务！")
}

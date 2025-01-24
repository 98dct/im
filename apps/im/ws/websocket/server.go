package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"sync"
	"time"
)

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (a AckType) ToString() string {
	switch a {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}
	return "NoAck"
}

type Server struct {
	routes map[string]HandlerFunc
	addr   string

	ListenOn     string
	serverOption serverOption

	*threading.TaskRunner
	mu         sync.RWMutex
	connToUser map[*Conn]string
	userToConn map[string]*Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOption) *Server {
	serverOption := newServerOption(opts...)
	s := &Server{
		routes: make(map[string]HandlerFunc),
		addr:   addr,
		//authentication: new(authentication),
		ListenOn:     FigureOutListenOn(addr),
		serverOption: serverOption,
		TaskRunner:   threading.NewTaskRunner(serverOption.concurrency),
		connToUser:   make(map[*Conn]string),
		userToConn:   make(map[string]*Conn),
		upgrader:     websocket.Upgrader{},
		Logger:       logx.WithContext(context.Background()),
	}

	s.serverOption.discover.Register(s.ListenOn)
	return s
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

	// 服务注册
	// 将用户id和当前im服务的地址注册到redis中 hash存储
	s.serverOption.discover.BoundUser(conn.Uid)

	go s.handlerWrite(conn)

	if s.isAck(nil) {
		go s.readAck(conn)
	}
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

		if s.isAck(&message) {
			s.Infof("conn message is ack message: %v", message)
			// 服务端开启ack机制 && 当前msg是ack message
			conn.appendMsgMq(&message)
		} else {
			conn.message <- &message
		}

	}
}

func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.serverOption.ack != NoAck
	}
	return s.serverOption.ack != NoAck && message.FrameType != FrameNoAck && message.FrameType != FrameTranspond
}

func (s *Server) readAck(conn *Conn) {
	for {
		select {
		case <-conn.done:
			return
		default:
		}

		conn.messageMu.Lock()
		// 从队列中读取新的消息
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			time.Sleep(3 * time.Second)
			continue
		}

		// 读取第一条消息
		message := conn.readMessage[0]

		// 判断ack的类型
		switch s.serverOption.ack {
		case OnlyAck:
			s.Send(&Message{
				Id:        message.Id,
				FrameType: FrameAck,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()

			conn.message <- message
		case RigorAck:

			if message.AckSeq == 0 {
				// 消息确认
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()

				// 发送
				s.Send(&Message{
					Id:        message.Id,
					FrameType: FrameAck,
					AckSeq:    message.AckSeq,
				}, conn)
				conn.messageMu.Unlock()
				continue
			}

			// 客户端返回结果，在一次确认

			// 客户端没有确认，考虑是否超过ack超时时间
			//     未超过，重新发送
			//     超过，直接结束确认
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 客户端已确认
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				continue
			}

			// 客户端没有确认消息
			val := s.serverOption.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				// 已超时
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}
			// 未超时，重新发送
			conn.messageMu.Unlock()
			s.Send(&Message{
				Id:        message.Id,
				FrameType: FrameAck,
				AckSeq:    message.AckSeq,
			}, conn)

			time.Sleep(3 * time.Second)
		}

	}
}

func (s *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			return
		case message := <-conn.message:
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
				handlerFunc(s, conn, message)

				if s.isAck(message) {
					conn.messageMu.Lock()
					delete(conn.readMessageSeq, message.Id)
					conn.messageMu.Unlock()
				}
			}
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

func (s *Server) GetConns(uIds ...string) ([]*Conn, []string) {
	if len(uIds) == 0 {
		return nil, nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	conns := make([]*Conn, 0, len(uIds))
	var noExistUids []string
	for _, uId := range uIds {
		c, ok := s.userToConn[uId]
		if ok {
			conns = append(conns, c)
			continue
		}
		noExistUids = append(noExistUids, uId)
	}

	return conns, noExistUids
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

	s.serverOption.discover.RelieveUser(uid)

	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, uIds ...string) error {

	if len(uIds) == 0 {
		return nil
	}

	conns, noExistUids := s.GetConns(uIds...)
	err := s.Send(msg, conns...)
	if err != nil {
		return err
	}

	return s.serverOption.discover.Transpond(msg, noExistUids...)

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

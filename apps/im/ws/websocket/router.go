package websocket

type Router struct {
	Method  string
	Handler HandlerFunc
}

type HandlerFunc func(srv *Server, conn *Conn, msg *Message)

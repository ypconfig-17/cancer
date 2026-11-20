package cancer

import (
	"cancer/log"
	"cancer/store"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"sync"
)

var logger = log.RegisterScope("/server")

func NewAppServer() *AppServer {
	return &AppServer{
		name: "AppServer",
	}
}

type AppServer struct {
	sync.Mutex
	name string
	port int
}

func (server *AppServer) Bind(port int) {
	server.port = port
}

func (server *AppServer) Close() {
	server.Lock()
	defer server.Unlock()
}
func (server *AppServer) Accept(pattern string) {
	server.Lock()
	defer server.Unlock()

	// 初始化存储
	store, err := store.OpenStore("db", false)
	if err != nil {
		logger.Errorf("open store fail,error:%s", err)
		return
	}
	subMgr.SetStore(store)
	pubMgr.SetStore(store)

	http.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websocket.Server{
			Handshake: func(config *websocket.Config, request *http.Request) error {
				return nil
			},
			Handler: func(conn *websocket.Conn) {
				logger.Infof("客户端连接,Addr: %s,", conn.Request().RemoteAddr)
				ch := ChannelMgr.Create(conn)
				pipeline := ch.GetPipe()
				pipeline.AddLast(StompEncoderHandler)
				pipeline.AddLast(FrameDecoderHandler)
				pipeline.AddLast(MessageEnqueueHandler)
				ch.Start()

			},
		}.ServeHTTP(w, r)
	}))

	logger.Infof("server 启动成功,port:%d", server.port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", server.port), nil) //only local
	if err != nil {
		logger.Errorf("init Serve: " + err.Error())
	}

}

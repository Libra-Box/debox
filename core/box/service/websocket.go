package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ipfs/kubo/core/box/model"
)

var WsConMap = make(map[string]*websocket.Conn)

func (s *HttpServer) WebSocketHandler(c *gin.Context) {
	upgrader := websocket.Upgrader{
		//ReadBufferSize:  1024,
		//WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 获取WebSocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Infof("error: %s", err.Error())
		return
	}
	defer ws.Close()
	var req model.Req
	// 处理WebSocket消息
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Errorf("ReadMessage:%s", err)
			return
		}
		log.Infof("ReadMessage str:%s", msg)
		err = json.Unmarshal(msg, &req)
		if err != nil {
			log.Errorf("Unmarshal:%s", err)
			return
		}
		WsConMap[req.From] = ws
		s.initWsRoute(req, string(msg))
	}
}

func (s *HttpServer) initWsRoute(req model.Req, reqStr string) {
	if req.Protocol == "/user/uploadInfo" {
		go s.UploadInfo(req)
	}
	if req.Protocol == "/user/add_qr" {
		go s.Add_qr(req)
	}
	if req.Protocol == "/user/agree" {
		go s.agree(req, string(reqStr))
	}
}

// func SendWsMsg(to string, res model.Res) {
// 	respJson, _ := json.Marshal(res)
// 	log.Infof("res:%s", respJson)
// 	err := WsConMap[to].WriteMessage(websocket.TextMessage, respJson)
// 	if err != nil {
// 		log.Error(err)
// 	}
// }

func SendWsMsg(to string, msg []byte) error {
	log.Infof("WSmsg:%s", msg)
	err := WsConMap[to].WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

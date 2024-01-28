package service

import (
	"bufio"
	"strings"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p/core/network"
)

func (s *HttpServer) relay(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	str, err := rw.ReadString('\n')
	if err != nil {
		log.Info(err)
		relay_Send(err.Error(), rw)
		return
	}
	str = strings.TrimSuffix(str, "\n")
	log.Infof("relay:%s", str)
	if err != nil {
		log.Error(err)
		relay_Send(err.Error(), rw)
		return
	}
	dataJson := jsonvalue.MustUnmarshalString(string(str))
	to, err := dataJson.GetString("to")
	if err != nil {
		log.Error(err)
		relay_Send(err.Error(), rw)
		return
	}
	err = WsConMap[to].WriteMessage(websocket.TextMessage, []byte(str))
	if err != nil {
		log.Error(err)
		relay_Send(err.Error(), rw)
		return
	}
	relay_Send(string("true"), rw)
}

func relay_Send(str string, rw *bufio.ReadWriter) {
	_, err := rw.WriteString(str + "\n")
	if err != nil {
		log.Error(err)
	}
	rw.Flush()
}

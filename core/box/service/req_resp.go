package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ipfs/kubo/core/box/model"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func SendRequest(p2pHost host.Host, peerInfo peer.AddrInfo, req model.Req, protocol protocol.ID) (resp string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// info := peer.AddrInfo{
	// 	Addrs: make([]ma.Multiaddr, 0),
	// }
	if network.Connected == p2pHost.Network().Connectedness(peerInfo.ID) {
		//info.ID = peerInfo.ID
	} else {
		log.Infof("try connect RemotePeers %v", peerInfo.String())
		err = p2pHost.Connect(ctx, peerInfo)
		if err != nil {
			log.Errorf("connect RemotePeers error %v", err.Error())
			return "", errors.New("RemotePeers not connect")
		}
	}
	log.Infof("send PeerId:%s", peerInfo.ID)
	s, err := p2pHost.NewStream(ctx, peerInfo.ID, protocol)
	if err != nil {
		log.Errorf("SendRequest %v", err)
		return "", err
	}
	defer s.Close()
	data, err := json.Marshal(req)
	if err != nil {
		log.Errorf("SendRequest %v", err)
		return "", err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	_, err = rw.WriteString(fmt.Sprintf("%s\n", data))
	if err != nil {
		log.Errorf("SendRequest %v", err)
		return "", err
	}
	rw.Flush() //
	ch := make(chan string, 1)
	go func() {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Info(err)
		}
		ch <- strings.TrimSuffix(str, "\n")
	}()
loop:
	for {
		select {
		case c := <-ch:
			resp = c
			break loop
		case <-ctx.Done():
			s.Close()
			log.Infof("request time out")
			resp = ""
			break loop
			// default:
			// 	str, err := rw.ReadString('\n')
			// 	if err != nil {
			// 		log.Info(err)
			// 	}
			// 	resp = strings.TrimSuffix(str, "\n")

		}
	}
	if resp != "" {
		if resp != "true" {
			return "", fmt.Errorf(resp)
		}
		return resp, nil
	}
	return "", fmt.Errorf("request timeout")
}

// func readResponse(rw *bufio.ReadWriter) string {
// 	for {
// 		str, err := rw.ReadString('\n')
// 		if err != nil {
// 			log.Info(err)
// 		}
// 		return strings.TrimSuffix(str, "\n")
// 	}
// }

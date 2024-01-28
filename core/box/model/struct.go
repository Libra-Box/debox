package model

import (
	"github.com/libp2p/go-libp2p/core/protocol"
)

// common
type Req struct {
	Protocol   protocol.ID `json:"protocol"`
	From       string      `json:"from"`
	To         string      `json:"to"`
	RemotePeer string      `json:"remotePeer"`
	Data       interface{} `json:"data"`
}

type Res struct {
	Protocol protocol.ID `json:"protocol"`
	Code     bool        `json:"code"`
	Msg      string      `json:"msg"`
	From     string      `json:"from"`
	Data     interface{} `json:"data"`
}

type LocalPeer struct {
	LocalPeer string
}

type Add_qr_data_res struct {
	NickName  string `json:"nickName"`
	HeadImg   string `json:"headImg"`
	PubilcKey string `json:"pubilcKey"`
}

type Agree_data_res struct {
	AesKey string `json:"aesKey"`
}

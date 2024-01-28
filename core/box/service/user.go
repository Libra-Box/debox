package service

import (
	"encoding/json"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/ipfs/kubo/core/box/model"
)

func (s *HttpServer) UploadInfo(req model.Req) {
	user := model.User{}
	res := model.Res{}
	res.Protocol = req.Protocol
	data, err := json.Marshal(req.Data)
	if err != nil {
		log.Error(err)
		res.Code = false
		res.Msg = err.Error()
		msg, _ := json.Marshal(res)
		SendWsMsg(req.To, msg)
		return
	}
	resData := jsonvalue.MustUnmarshalString(string(data))
	headImg, _ := resData.GetString("headImg")
	nickName, _ := resData.GetString("nickName")
	user.HeadImg = headImg
	user.NickName = nickName
	user.Id = req.From //user publicKey md5
	err = s.Store.InserUserInfo(user)
	if err != nil {
		log.Error(err)
		res.Code = false
		res.Msg = err.Error()
		msg, _ := json.Marshal(res)
		SendWsMsg(req.To, msg)
		return
	}
	res.Code = true
	res.Data = model.LocalPeer{
		LocalPeer: s.getLocalPeer(),
	}
	msg, _ := json.Marshal(res)
	SendWsMsg(req.To, msg)
}

func (s *HttpServer) Add_qr(req model.Req) {
	res := model.Res{}
	res.Protocol = req.Protocol
	res.From = req.To
	remotePeers, err := s.Cfg.ParseBootstrapPeers2([]string{req.RemotePeer})
	if err != nil {
		res.Code = false
		res.Msg = err.Error()
		msg, _ := json.Marshal(res)
		SendWsMsg(req.From, msg)
		return
	}
	if remotePeers[0].ID != s.P2pHost.ID() {
		_, err := SendRequest(s.P2pHost, remotePeers[0], req, "/relay")
		if err != nil {
			res.Code = false
			res.Msg = "remote error:" + err.Error()
			msg, _ := json.Marshal(res)
			SendWsMsg(req.From, msg)
			return
		}
		//send sucessful
		res.Code = true
		msg, _ := json.Marshal(res)
		SendWsMsg(req.From, msg)
	} else {
		//发送自己的信息给对方
		msg, _ := json.Marshal(req)
		err = SendWsMsg(req.To, msg)
		if err != nil {
			res.Code = true
			msg, _ := json.Marshal(res)
			SendWsMsg(req.From, msg)
		}

		// //查找对方信息  把对方的信息发给自己
		// user, err := s.Store.GetUserInfo(req.To)
		// if err != nil {
		// 	log.Error(err)
		// }
		// res.Code = true
		// res.Protocol = req.Protocol
		// res.From = req.To
		// res.To = req.From
		// res.Data = model.Add_qr_data_res{
		// 	NickName: user.NickName,
		// 	HeadImg:  user.HeadImg,
		// }
		// res.RemotePeer = s.getLocalPeer()
		// resStr, _ := json.Marshal(res)
		// SendWsMsg(req.From, resStr)
	}
}

func (s *HttpServer) agree(req model.Req, reqStr string) {
	res := model.Res{}
	res.Protocol = req.Protocol
	res.From = req.To
	remotePeers, err := s.Cfg.ParseBootstrapPeers2([]string{req.RemotePeer})
	if err != nil {
		res.Code = false
		res.Msg = err.Error()
		msg, _ := json.Marshal(res)
		SendWsMsg(req.From, msg)
		return
	}
	if remotePeers[0].ID != s.P2pHost.ID() {
		_, err := SendRequest(s.P2pHost, remotePeers[0], req, "/relay")
		if err != nil {
			res.Code = false
			res.Msg = "remote error:" + err.Error()
			msg, _ := json.Marshal(res)
			SendWsMsg(req.From, msg)
			return
		}
		//send sucessful
		res.Code = true
		msg, _ := json.Marshal(res)
		SendWsMsg(req.From, msg)
	} else {
		//send to
		msg, _ := json.Marshal(req)
		err = SendWsMsg(req.To, msg)
		if err != nil {
			res.Code = true
			msg, _ := json.Marshal(res)
			SendWsMsg(req.From, msg)
		}
	}

}

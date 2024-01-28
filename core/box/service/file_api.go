package service

//
//import (
//	"context"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/golang/protobuf/proto"
//	"github.com/ipfs/go-cid"
//	files "github.com/ipfs/go-ipfs-files"
//	"github.com/ipfs/interface-go-ipfs-core/options"
//	"github.com/ipfs/interface-go-ipfs-core/path"
//	"github.com/ipfs/kubo/core/box/model"
//	"github.com/ipfs/kubo/pkg/xfile"
//	"github.com/jinzhu/gorm"
//	"github.com/libp2p/go-libp2p/core/protocol"
//	"go4.org/lock"
//	"io"
//	"io/ioutil"
//	"os"
//	"strconv"
//	"strings"
//	"time"
//)
//

//func (s *HttpServer) UploadFile_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	//protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	//peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.CommonResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.UploadFileReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	//log.Infof("ContentType: %v", c.ContentType())
//	//log.Infof("PROTOCOL: %v", protoc)
//	//log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//
//	if req.Id == "" {
//		req.Id = genUUidString()
//	}
//	u := s.ctx.Value("user").(UserData)
//	if req.Size < req.BytesFrom+int64(len(req.Data)) {
//		log.Errorf("file size error: %+v", req)
//		resp.Code = pb.Code_RequestParamError
//
//		respondProto(c, resp)
//		return
//	}
//	user, err := s.store.GetUserById(u.Id)
//	if err != nil {
//		log.Errorf("failed to get user: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	if uint64(req.Size) > user.AllocatedSpace-user.UsedSpace {
//		log.Errorf("no enough space: %v", user.AllocatedSpace-user.UsedSpace)
//		resp.Code = pb.Code_NoEnoughSpace
//		respondProto(c, resp)
//		return
//	}
//
//	FormDevice := ""
//	parent, err := s.store.GetFileById(req.ParentId)
//	if err != nil {
//		log.Errorf("failed to get file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	if !parent.IsFolder || parent.UserId != user.Id {
//		log.Errorf("parent is not folder: %v", err)
//		resp.Code = pb.Code_RequestParamError
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	FormDevice = parent.FormDevice
//	file := model.File{
//		Id:         req.Id,
//		UserId:     user.Id,
//		ParentId:   req.ParentId,
//		Name:       req.Name,
//		Ext:        xfile.Ext(req.Name),
//		Size:       int(req.Size),
//		CreatedAt:  time.Now().Unix(),
//		UpdateAt:   time.Now().Unix(),
//		FormDevice: FormDevice,
//	}
//	if req.Md5 != "" {
//		resMd5, errs := s.store.SearchFileMd5List(req.Md5)
//		if errs != nil {
//			log.Errorf("faild to SearchFileMd5List: %v", errs)
//			resp.Code = pb.Code_Failure
//			respondProto(c, resp)
//			return
//		}
//		if len(resMd5) > 0 {
//			for _, item := range resMd5 {
//				if item.Name == req.Name && item.ParentId == req.ParentId {
//					log.Errorf("file exist")
//					resp.Code = pb.Code_FileNameExist
//					respondProto(c, resp)
//					return
//				}
//			}
//			log.Errorf("FileMd5 is exist: %v", errs)
//			file.Cid = resMd5[0].Cid
//			file.Md5 = resMd5[0].Md5
//			err = s.store.CreateItem(&file)
//			if err != nil {
//				log.Errorf("failed to crete file: %v", err)
//				resp.Code = pb.Code_IoError
//				respondProto(c, resp)
//				return
//			}
//			resp.Code = pb.Code_Md5IsExist
//			respondProto(c, resp)
//			return
//		} else {
//			if old, err := s.store.GetFileInFolder(user.Id, req.ParentId, req.Name, ""); err == nil {
//				if old.Id != req.Id {
//					log.Errorf("file exist")
//					resp.Code = pb.Code_FileNameExist
//					respondProto(c, resp)
//					return
//				}
//			}
//		}
//	}
//
//	boxName := "box000!"
//	fileName := s.TempDir + "/" + fmt.Sprintf("%s%s!%s", boxName, req.Id, req.Name)
//	if req.BytesFrom == 0 {
//		_, err := os.Stat(fileName)
//		if err == nil { // file exist
//			err = os.Remove(fileName)
//			if err != nil {
//				log.Errorf("failed to remove file: %v", err)
//				resp.Code = pb.Code_IoError
//				respondProto(c, resp)
//				return
//			}
//		}
//	} else {
//		fstat, err := os.Stat(fileName)
//		if err != nil {
//			log.Errorf("failed to get file stat: %v", err)
//			resp.Code = pb.Code_IoError
//			respondProto(c, resp)
//			return
//		}
//		if fstat.Size() != req.BytesFrom {
//			log.Errorf("file size error:%v fstat.Size:%v-req.BytesFrom%v", err, fstat.Size(), req.BytesFrom)
//			resp.Code = pb.Code_RequestParamError
//			curSize := strconv.FormatInt(fstat.Size(), 10)
//			resp.Msg = curSize
//			respondProto(c, resp)
//			return
//		}
//	}
//
//	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//	if err != nil {
//		log.Errorf("failed to open file: %v", err)
//		resp.Code = pb.Code_IoError
//		respondProto(c, resp)
//		return
//	}
//	_, err = fp.Write(req.Data)
//	defer fp.Close()
//	if err != nil {
//		log.Errorf("failed to write file: %v", err)
//		resp.Code = pb.Code_IoError
//		respondProto(c, resp)
//		return
//	}
//
//	if req.Size > req.BytesFrom+int64(len(req.Data)) {
//		respondProto(c, resp)
//		return
//	}
//
//	dbFile, err := s.store.GetFileById(file.Id)
//	if err == nil {
//		// update
//		user.UsedSpace -= uint64(dbFile.Size)
//		err := s.store.UpdateFileForUpload(&file)
//		if err != nil {
//			log.Errorf("failed to update file: %v", err)
//			resp.Code = pb.Code_IoError
//			respondProto(c, resp)
//			return
//		}
//	} else {
//		// create
//		err = s.store.CreateItem(&file)
//		if err != nil {
//			log.Errorf("failed to crete file: %v", err)
//			resp.Code = pb.Code_IoError
//			respondProto(c, resp)
//			return
//		}
//	}
//	user.UsedSpace += uint64(req.Size)
//	err = s.store.UpdateUserSpace(user)
//	if err != nil {
//		log.Errorf("failed to update user: %v", err)
//		resp.Code = pb.Code_IoError
//		respondProto(c, resp)
//		return
//	}
//	//log.Errorf("end_ss: %v", time.Now())
//
//	respondProto(c, resp)
//
//	if s.coreApi != nil {
//		go func() {
//			//log.Errorf("start_c: %v", time.Now())
//			err = s.store.UpdateFileMd5(req.Id, GetFileMd5(fileName))
//			if err != nil {
//				log.Errorf("failed to UpdateFileMd5: %v", err)
//			}
//			//log.Errorf("end_c: %v", time.Now())
//			//log.Errorf("start: %v", time.Now())
//			f, err := os.Open(fileName)
//			if err != nil {
//				log.Errorf("open file error: %v", err)
//				return
//			}
//			defer f.Close()
//			node, err := s.coreApi.Unixfs().Add(context.Background(), files.NewReaderFile(f), options.Unixfs.Pin(true))
//			if err != nil {
//				boxName = "box111"
//				newName := s.TempDir + "/" + fmt.Sprintf("%s!%s!%s", boxName, req.Id, req.Name)
//				os.Rename(fileName, newName)
//				log.Errorf("failed to add file to ipfs: %v", err)
//			} else {
//				file.Cid = node.Cid().String()
//				err = s.store.UpdateFileCid(req.Id, node.Cid().String())
//				if err != nil {
//					if strings.Contains(err.Error(), "locked") == true {
//						os.Remove(fileName)
//					} else {
//						boxName = "box111"
//						newName := s.TempDir + "/" + fmt.Sprintf("%s!%s!%s", boxName, req.Id, req.Name)
//						os.Rename(fileName, newName)
//					}
//					log.Errorf("failed to UpdateFileCid: %v", err)
//				} else {
//					os.Remove(fileName)
//				}
//			}
//			//log.Errorf("end: %v", time.Now())
//		}()
//	}
//
//}
//
//func (s *HttpServer) DownloadFile_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.DownloadFileResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.DownloadFileReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//
//	u := s.ctx.Value("user").(UserData)
//	file, err := s.store.GetFileById(req.FileId)
//	if err != nil {
//		log.Errorf("failed to get file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	if file.Cid == "" {
//		log.Errorf("file cid not exist: %v", err)
//		resp.Code = pb.Code_UserNameNotExist
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("u.Id: %v", u.Id)
//	log.Infof("UserId: %v", req.UserId)
//	log.Infof("file.UserId: %v", file.UserId)
//	if req.UserId != 0 {
//		if int32(u.Id) == req.UserId {
//			if file.UserId != u.Id {
//				log.Errorf("user: %v", err)
//				resp.Code = pb.Code_RequestParamError
//				respondProto(c, resp)
//				return
//			}
//		} else {
//			if int32(file.UserId) != req.UserId {
//				log.Errorf("user: %v", err)
//				resp.Code = pb.Code_RequestParamError
//				respondProto(c, resp)
//				return
//			}
//		}
//	} else {
//		if file.UserId != u.Id {
//			log.Errorf("user: %v", err)
//			resp.Code = pb.Code_RequestParamError
//			respondProto(c, resp)
//			return
//		}
//	}
//	//log.Errorf("start: %v", time.Now())
//	cid, err := cid.Decode(file.Cid)
//	if err != nil {
//		log.Errorf("failed to decode cid %v", err.Error())
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	cidPath := path.IpfsPath(cid)
//	fileNode, err := s.coreApi.Unixfs().Get(context.Background(), cidPath)
//	if err != nil {
//		log.Errorf("failed to get file %v", err.Error())
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	buf, err := getFileBytes(fileNode, req.BytesFrom, int(req.BytesCount))
//	if err != nil {
//		if err != io.EOF {
//			log.Errorf("failed to getFileBytes %v", err.Error())
//			resp.Code = pb.Code_RequestParamError
//			respondProto(c, resp)
//			return
//		}
//	}
//	defer fileNode.Close()
//	resp.Data = buf
//	resp.Name = file.Name
//	resp.FileSize = int64(file.Size)
//	respondProto(c, resp)
//	return
//}
//
//func (s *HttpServer) DownloadFile_http1(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.DownloadFileResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.DownloadFileReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//
//	u := s.ctx.Value("user").(UserData)
//	file, err := s.store.GetFileById(req.FileId)
//	if err != nil {
//		log.Errorf("failed to get file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	if file.Cid == "" {
//		log.Errorf("file cid not exist: %v", err)
//		resp.Code = pb.Code_UserNameNotExist
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("u.Id: %v", u.Id)
//	log.Infof("UserId: %v", req.UserId)
//	log.Infof("file.UserId: %v", file.UserId)
//	if req.UserId != 0 {
//		if int32(u.Id) == req.UserId {
//			if file.UserId != u.Id {
//				log.Errorf("user: %v", err)
//				resp.Code = pb.Code_RequestParamError
//				respondProto(c, resp)
//				return
//			}
//		} else {
//			if int32(file.UserId) != req.UserId {
//				log.Errorf("user: %v", err)
//				resp.Code = pb.Code_RequestParamError
//				respondProto(c, resp)
//				return
//			}
//		}
//	} else {
//		if file.UserId != u.Id {
//			log.Errorf("user: %v", err)
//			resp.Code = pb.Code_RequestParamError
//			respondProto(c, resp)
//			return
//		}
//	}
//	log.Errorf("start: %v", time.Now())
//	tempFile := s.TempDir + "/" + fmt.Sprintf("%s-%s", file.ParentId, file.Name)
//	fstat, err := os.Stat(tempFile)
//	if err != nil {
//		if os.IsNotExist(err) {
//			err = s.retrieveFile(file.Cid, tempFile)
//			if err != nil {
//				log.Errorf("failed to retrieveFile %v", err.Error())
//				resp.Code = pb.Code_IoError
//				respondProto(c, resp)
//				return
//			}
//		} else {
//			log.Errorf("filed to read file: %v", err)
//			resp.Code = pb.Code_IoError
//			respondProto(c, resp)
//			return
//		}
//
//	} else {
//		if fstat.Size() != int64(file.Size) {
//			err = os.Remove(tempFile)
//			if err != nil {
//				log.Errorf("filed to remove file: %v", err)
//				resp.Code = pb.Code_IoError
//				respondProto(c, resp)
//				return
//			}
//			err = s.retrieveFile(file.Cid, tempFile)
//			if err != nil {
//				log.Errorf("failed to retrieveFile %v", err.Error())
//				resp.Code = pb.Code_IoError
//				respondProto(c, resp)
//				return
//			}
//		}
//	}
//	log.Errorf("end: %v", time.Now())
//	log.Errorf("start_o: %v", time.Now())
//	fp, err := os.Open(tempFile)
//	if err != nil {
//		log.Errorf("failed to open file %v", err.Error())
//		resp.Code = pb.Code_IoError
//		respondProto(c, resp)
//		return
//	}
//	log.Errorf("end_o: %v", time.Now())
//	defer fp.Close()
//
//	buf := make([]byte, req.BytesCount)
//	log.Errorf("start_r: %v", time.Now())
//	n, err := fp.ReadAt(buf, req.BytesFrom)
//	if int64(n) < req.BytesCount {
//		buf = buf[:n]
//		if n+int(req.BytesFrom) != file.Size {
//			resp.Code = pb.Code_IoError
//			log.Errorf("filed to read file: %v", err)
//			respondProto(c, resp)
//			return
//		}
//	}
//	log.Errorf("end_r: %v", time.Now())
//	resp.Data = buf
//	resp.Name = file.Name
//	resp.FileSize = int64(file.Size)
//	if n+int(req.BytesFrom) == file.Size {
//		fp.Close()
//		os.Remove(tempFile)
//	}
//	respondProto(c, resp)
//	return
//}
//
//func (s *HttpServer) retrieveFile(fileCid string, outputPath string) error {
//	lock.Lock()
//	defer lock.Unlock()
//	cid, err := cid.Decode(fileCid)
//	if err != nil {
//		log.Errorf("failed to decode cid %v", err.Error())
//		return err
//	}
//	cidPath := path.IpfsPath(cid)
//	fileNode, err := s.coreApi.Unixfs().Get(context.Background(), cidPath)
//	if err != nil {
//		log.Errorf("failed to get file %v", err.Error())
//		return err
//	}
//	defer fileNode.Close()
//
//	err = files.WriteTo(fileNode, outputPath)
//	if err != nil {
//		log.Errorf("failed to write file %v", err.Error())
//		return err
//	}
//	return err
//}
//
//func (s *HttpServer) GetFileList_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.FileListResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.FileListReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("req: %v", req)
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//
//	u := s.ctx.Value("user").(UserData)
//	count, list, err := s.store.FileList(u.Id, int(req.DirMask), req.ParentId,
//		int(req.FileType), int(req.StarMask), req.Keyword, int(req.Order), int(req.Limit), int(req.Offset), int(req.IsEqual))
//	if err != nil {
//		log.Errorf("failed to get file list: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	resp.Total = int32(count)
//	resp.Files = make([]*pb.FileItem, len(list))
//	for i, v := range list {
//		uCount := 0
//		if v.UserList != "" {
//			parts := strings.Split(v.UserList, ",")
//			//log.Errorf("parts: %v", len(parts))
//			uCount = len(parts)
//		}
//		//log.Errorf("uCount: %v", uCount)
//		//log.Errorf("uCount: %v", int32(uCount))
//		resp.Files[i] = &pb.FileItem{
//			Id:         v.Id,
//			Name:       v.Name,
//			Size:       int64(v.Size),
//			Md5:        v.Md5,
//			IsFolder:   v.IsFolder,
//			CreatedAt:  v.CreatedAt,
//			UpdatedAt:  v.UpdateAt,
//			Star:       v.Star,
//			Share:      v.Share,
//			SubFiles:   int32(v.SubFiles),
//			Ext:        v.Ext,
//			ParentId:   v.ParentId,
//			ParentName: v.ParentName,
//			ShareCount: int32(uCount),
//			IsSystem:   int32(v.IsSystem),
//			Cid:        v.Cid,
//			FormDevice: v.FormDevice,
//			StartAt:    v.StartAt,
//			EndAt:      v.EndAt,
//		}
//		if v.IsFolder {
//			resp.Files[i].Kind = "文件夹"
//		} else {
//			resp.Files[i].Kind = model.GetFileTypeString(v.Ext)
//		}
//		//if v.Share {
//		//	resp.Files[i].ShareCount = int32(strings.Count(v.UserList, ","))
//		//}
//	}
//	respondProto(c, resp)
//	return
//}
//
//func (s *HttpServer) FileRename_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.CommonResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.FileRenameReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//
//	u := s.ctx.Value("user").(UserData)
//	file, err := s.store.GetFileById(req.FileId)
//	if err != nil {
//		log.Errorf("faild to get file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	files := model.File{
//		Id:       req.FileId,
//		UserId:   u.Id,
//		ParentId: req.ParentId,
//		Name:     req.Name,
//	}
//	if _, err = s.store.GetFileInFolder(files.UserId, files.ParentId, files.Name, ""); err != gorm.ErrRecordNotFound {
//		log.Errorf("failed to create folder: %v", err)
//		resp.Code = pb.Code_FileNameExist
//		respondProto(c, resp)
//		return
//	}
//
//	if file.UserId != u.Id {
//		resp.Code = pb.Code_RequestParamError
//
//		respondProto(c, resp)
//		return
//	}
//	file.Name = req.Name
//	err = s.store.UpdateFileName(file)
//	if err != nil {
//		log.Errorf("faild to update file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	respondProto(c, resp)
//	return
//}
//
//func (s *HttpServer) FileStar_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.CommonResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.FileStarReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//
//	u := s.ctx.Value("user").(UserData)
//	err = s.store.UpdateFileStar(u.Id, req.FileIds, true)
//	if err != nil {
//		log.Errorf("faild to update file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	respondProto(c, resp)
//	return
//}
//
//func (s *HttpServer) FileUnstar_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.CommonResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.FileUnstarReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//
//	u := s.ctx.Value("user").(UserData)
//	err = s.store.UpdateFileStar(u.Id, req.FileIds, false)
//	if err != nil {
//		log.Errorf("faild to update file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	respondProto(c, resp)
//	return
//}
//
//func (s *HttpServer) FileMove_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.CommonResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.FileMoveReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//	FormDevice := ""
//	u := s.ctx.Value("user").(UserData)
//	file, err := s.store.GetFileById(req.FileId)
//	if err != nil {
//		log.Errorf("faild to get file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	if file.UserId != u.Id {
//		log.Errorf("user_id not match")
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	if file.ParentId == req.NewParentId {
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//
//	if req.NewParentId != "" && req.NewParentId != "desktop" {
//		newParent, err := s.store.GetFileById(req.NewParentId)
//		if err != nil {
//			log.Errorf("faild to get file parents: %v", err)
//			resp.Code = pb.Code_DbError
//			respondProto(c, resp)
//			return
//		}
//		if !newParent.IsFolder || newParent.UserId != u.Id {
//			log.Errorf("parent dir error")
//			resp.Code = pb.Code_RequestParamError
//			respondProto(c, resp)
//			return
//		}
//		FormDevice = newParent.FormDevice
//	}
//	if _, err := s.store.GetFileInFolder(u.Id, req.NewParentId, file.Name, ""); err == nil {
//		log.Errorf("file exist")
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("NewParentId: %v", req.NewParentId)
//	if file.IsFolder {
//		parents, err := s.store.GetFileParents(req.NewParentId)
//		if err != nil {
//			log.Errorf("faild to get file parents: %v", err)
//			resp.Code = pb.Code_DbError
//			respondProto(c, resp)
//			return
//		}
//		for _, v := range parents {
//			if v.Id == file.Id {
//				log.Errorf("param error")
//				resp.Code = pb.Code_RequestParamError
//				respondProto(c, resp)
//				return
//			}
//		}
//	}
//	file.ParentId = req.NewParentId
//	file.FormDevice = FormDevice
//	err = s.store.UpdateFileParent(file)
//	if err != nil {
//		log.Errorf("faild to update file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	respondProto(c, resp)
//	return
//}
//
//func (s *HttpServer) FileRecord_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.CommonResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.FileRecordReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//
//	u := s.ctx.Value("user").(UserData)
//	user, err := s.store.GetUserById(u.Id)
//	if err != nil {
//		log.Errorf("faild to get user: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	FormDevice := ""
//	parent, err := s.store.GetFileById(req.ParentId)
//	if err != nil {
//		log.Errorf("file is not exist: %v", err)
//		resp.Code = pb.Code_UserNameNotExist
//		respondProto(c, resp)
//		return
//	}
//	FormDevice = parent.FormDevice
//	if uint64(req.Size) > user.AllocatedSpace-user.UsedSpace {
//		log.Errorf("no enough space: %v", user.AllocatedSpace-user.UsedSpace)
//		resp.Code = pb.Code_NoEnoughSpace
//		respondProto(c, resp)
//		return
//	}
//	files := &model.File{
//		Id:         req.FileId,
//		UserId:     u.Id,
//		ParentId:   req.ParentId,
//		Name:       req.Name,
//		Cid:        req.Cid,
//		Md5:        req.Md5,
//		Size:       int(req.Size),
//		Ext:        xfile.Ext(req.Name),
//		CreatedAt:  time.Now().Unix(),
//		UpdateAt:   time.Now().Unix(),
//		FormDevice: FormDevice,
//	}
//	err = s.store.CreateItem(files)
//	if err != nil {
//		log.Errorf("failed to crete file: %v", err)
//		resp.Code = pb.Code_IoError
//		respondProto(c, resp)
//		return
//	}
//	respondProto(c, resp)
//	return
//}
//
//func (s *HttpServer) FileCopy_http(c *gin.Context) {
//	lock.Lock()
//	defer lock.Unlock()
//	protoc := protocol.ID(c.Request.Header.Get("PROTOCOL"))
//	peerId := c.Request.Header.Get("PEER_ID")
//	resp := &pb.CommonResp{
//		Code: pb.Code_Success,
//	}
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	req := &pb.FileCopyReq{}
//	err = proto.UnmarshalMerge(body, req)
//	if err != nil {
//		log.Errorf(" request error : %v", err)
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("ContentType: %v", c.ContentType())
//	log.Infof("PROTOCOL: %v", protoc)
//	log.Infof("PEER_ID: %v", peerId)
//	//token检查
//	tokenCode := s.LoginRequired(req.Token)
//	if tokenCode != pb.Code_Success {
//		resp.Code = tokenCode
//		respondProto(c, resp)
//		return
//	}
//	FormDevice := ""
//	u := s.ctx.Value("user").(UserData)
//	user, err := s.store.GetUserById(u.Id)
//	if err != nil {
//		log.Errorf("faild to get user: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//
//	file, err := s.store.GetFileById(req.FileId)
//	if err != nil {
//		log.Errorf("faild to get file: %v", err)
//		resp.Code = pb.Code_DbError
//		respondProto(c, resp)
//		return
//	}
//	log.Infof("u.Id: %v", u.Id)
//	log.Infof("UserId: %v", req.UserId)
//	log.Infof("file.UserId: %v", file.UserId)
//	if req.UserId != 0 {
//		if int32(u.Id) == req.UserId {
//			if file.UserId != u.Id {
//				log.Errorf("user_id not match")
//				resp.Code = pb.Code_RequestParamError
//				respondProto(c, resp)
//				return
//			}
//		} else {
//			if int32(file.UserId) != req.UserId {
//				log.Errorf("user_id not match")
//				resp.Code = pb.Code_RequestParamError
//				respondProto(c, resp)
//				return
//			}
//		}
//	} else {
//		if file.UserId != u.Id {
//			log.Errorf("user_id not match")
//			resp.Code = pb.Code_RequestParamError
//			respondProto(c, resp)
//			return
//		}
//	}
//
//	if file.ParentId == req.NewParentId && file.Name == req.NewFileName {
//		log.Errorf("file name duplicated")
//		resp.Code = pb.Code_RequestParamError
//		respondProto(c, resp)
//		return
//	}
//
//	if uint64(file.Size) > user.AllocatedSpace-user.UsedSpace {
//		log.Errorf("no enough space: %v", user.AllocatedSpace-user.UsedSpace)
//		resp.Code = pb.Code_NoEnoughSpace
//		respondProto(c, resp)
//		return
//	}
//
//	if req.NewParentId != "" && req.NewParentId != "desktop" {
//		newParent, err := s.store.GetFileById(req.NewParentId)
//		if err != nil {
//			log.Errorf("faild to get file parents: %v", err)
//			resp.Code = pb.Code_DbError
//			respondProto(c, resp)
//			return
//		}
//		if !newParent.IsFolder || newParent.UserId != user.Id {
//			log.Errorf("parent dir error")
//			resp.Code = pb.Code_RequestParamError
//			respondProto(c, resp)
//			return
//		}
//		FormDevice = newParent.FormDevice
//	}
//
//	if _, err := s.store.GetFileInFolder(user.Id, req.NewParentId, req.NewFileName, ""); err == nil {
//		log.Errorf("file exist")
//		resp.Code = pb.Code_FileNameExist
//		respondProto(c, resp)
//		return
//	}
//
//	if !file.IsFolder {
//		_, err = s.copyFile(file, req.NewParentId, req.FileId, req.NewFileName, u.Id, FormDevice)
//		if err != nil {
//			resp.Code = pb.Code_DbError
//			respondProto(c, resp)
//			return
//		}
//	} else {
//		parents, err := s.store.GetFileParents(req.NewParentId)
//		if err != nil {
//			log.Errorf("faild to get file parents: %v", err)
//			resp.Code = pb.Code_DbError
//			respondProto(c, resp)
//			return
//		}
//		for _, v := range parents {
//			if v.Id == file.Id {
//				log.Errorf("param error")
//				resp.Code = pb.Code_RequestParamError
//				respondProto(c, resp)
//				return
//			}
//		}
//
//		files, err := s.store.GetFileChildrenRecursively(file.Id)
//		if err != nil {
//			log.Errorf("failed to get files: %v", err)
//			resp.Code = pb.Code_DbError
//			respondProto(c, resp)
//			return
//		}
//
//		children := make(map[string][]string, 0)
//		idFile := make(map[string]*model.File, 0)
//
//		idFile["all"] = &model.File{IsFolder: true}
//		for _, v := range files {
//			idFile[v.Id] = v
//			if nodes, ok := children[v.ParentId]; ok {
//				nodes = append(nodes, v.Id)
//				children[v.ParentId] = nodes
//			} else {
//				nodes := make([]string, 0)
//				nodes = append(nodes, v.Id)
//				children[v.ParentId] = nodes
//			}
//		}
//		err = s.copyFiles(file, req.NewParentId, req.FileId, req.NewFileName, children, idFile, u.Id, FormDevice)
//		if err != nil {
//			resp.Code = pb.Code_DbError
//			respondProto(c, resp)
//			return
//		}
//	}
//	respondProto(c, resp)
//	return
//}

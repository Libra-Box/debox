package service

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core/box/ds"
	"github.com/ipfs/kubo/repo/fsrepo"
	"github.com/jinzhu/gorm"
	"github.com/klauspost/cpuid/v2"
	"github.com/libp2p/go-libp2p/core/host"
)

var log = logging.Logger("service")

type HttpServer struct {
	HttpServer *http.Server
	P2pHost    host.Host
	CoreApi    coreiface.CoreAPI
	Store      *ds.DbStore
	TempDir    string
	Cfg        *config.Config
	RootPath   string
}

func init() {
	logging.SetLogLevel("service", "info")
}

func (s *HttpServer) Init(coreApi coreiface.CoreAPI) {
	s.CoreApi = coreApi
}

//自定义Http服务
func NewHttpServer(cfg *config.Config, p2pHost host.Host) *HttpServer {
	repoPath, err := fsrepo.BestKnownPath()
	if err != nil {
		panic(err.Error())
	}
	initConfigBox()
	cfg_libra, _ := GetConfig()
	v := url.Values{}
	v.Set(`cache`, `shared`)
	v.Set(`mode`, `rwc`)
	u := url.URL{
		Scheme:   `file`,
		Opaque:   url.PathEscape(repoPath + "/box.db"),
		RawQuery: v.Encode(),
	}
	db, err := gorm.Open("sqlite3", u.String())
	db.DB().SetMaxOpenConns(0)
	//db, err := gorm.Open("sqlite3", repoPath+"/box.db")
	if err != nil {
		panic(err.Error())
	}
	db.SingularTable(true)
	dbStore := ds.NewDbStore(db)
	engine := gin.Default()
	gin.SetMode(gin.ReleaseMode) //日志模式
	server := &HttpServer{
		HttpServer: &http.Server{Addr: cfg_libra.HttpServer, Handler: engine},
		P2pHost:    p2pHost,
		Store:      dbStore,
		TempDir:    repoPath + "/temp",
		Cfg:        cfg,
		RootPath:   repoPath,
	}
	log.Infof("peer: %v", p2pHost.ID().String())
	log.Infof("CPUInfo %v", cpuid.CPU.BrandName)
	server.setCors(engine)
	server.initRoute(engine)
	os.MkdirAll(server.TempDir, 0777)
	return server
}

func (s *HttpServer) Run() {
	go func() {
		err := s.HttpServer.ListenAndServe()
		if err != nil {
			log.Infof("%v", err.Error())
		}
	}()
	//go func() {
	//	ticker := time.NewTicker(time.Second * 30)
	//	for {
	//		select {
	//		case <-ticker.C:
	//			peerAddress, _ := s.cfg.BootstrapPeers()
	//			for _, val := range peerAddress {
	//				if val.ID == s.p2pHost.ID() {
	//					continue
	//				}
	//				//log.Infof("connectInfo:%v", s.p2pHost.Network().Connectedness(val.ID).String())
	//				if s.p2pHost.Network().Connectedness(val.ID).String() != "Connected" {
	//					err := s.p2pHost.Connect(s.ctx, val)
	//					if err != nil {
	//						log.Errorf("BootstrapPeers connect error :%v", val)
	//					} else {
	//						log.Infof("BootstrapPeers connect successful:%v", val)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}()
	//
	//go func() {
	//	ticker := time.NewTicker(time.Minute * 10)
	//	for {
	//		select {
	//		case <-ticker.C:
	//
	//		}
	//	}
	//}()
}

func (s *HttpServer) initRoute(r gin.IRouter) {
	api := r.Group("/v1")
	api.GET("/ws", s.WebSocketHandler)

	//wsHandel
	//s.P2pHost.SetStreamHandler("/user/add_qr", s.Add_qr_relay)
	s.P2pHost.SetStreamHandler("/relay", s.relay)
}

func (s *HttpServer) setCors(r gin.IRouter) {
	corsCfg := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Access-Control-Allow-Origin", "Accept",
			"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-Requested-With", "PEER_ID", "PROTOCOL"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsCfg))
}

func (s *HttpServer) Shutdown() {
	err := s.HttpServer.Shutdown(context.Background())
	if err != nil {
		log.Errorf("%v", err.Error())
	}
}

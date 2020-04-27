package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"time"

	//
	_ "net/http/pprof"

	"github.com/eviltomorrow/my-develop-kit/plog"
	"github.com/eviltomorrow/my-develop-kit/transfer-go/config"
	"github.com/eviltomorrow/my-develop-kit/transfer-go/pb"
	"github.com/eviltomorrow/my-develop-kit/transfer-go/service"

	"google.golang.org/grpc"
)

const (
	nmConfig = "config"
)

var path = flag.String(nmConfig, "config.toml", "配置文件路径")

var cfg *config.Config

func main() {
	flag.Parse()
	loadConfig()
	setupLogger()
	printInfo()
	setGlobalVars()
	initEnv()
	startupService()
}

func loadConfig() {
	cfg = config.GetGlobalConfig()
	err := cfg.Load(*path)
	if err != nil {
		log.Fatalf("配置文件解析失败，nest error: %v", err)
	}
}

func setupLogger() {
	plog.InitLogger(cfg.ToLogConfig())
}

func printInfo() {
	plog.Infof("配置文件路径：%v", filepath.Join(os.Args[0], *path))
	plog.Infof("配置项信息：%s", cfg.String())
}

func setGlobalVars() {

}

func initEnv() {

}

func startupService() {
	go func() {
		if cfg.System.PProfListenPort == 0 {
			plog.Fatalf("pprof 服务端口未设置")
		}
		plog.Infof("启动 pprof 成功，监听端口 [%d]，访问页面：%s", cfg.System.PProfListenPort, fmt.Sprintf("http://localhost:%d/debug/pprof/", cfg.System.PProfListenPort))

		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.System.PProfListenPort), nil)
		if err != nil {
			plog.Fatalf("启动 pprof 失败，nest error: %v", err)
		}
	}()

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port))
	if err != nil {
		plog.Fatalf("启动监听失败，nest error: %v", err)
	}

	server := grpc.NewServer(
	// grpc.UnaryInterceptor(middleware.UnaryServerLogInterceptor),
	)

	pb.RegisterUploadFileServer(server, &service.UploadFile{})

	plog.Infof("启动 transfer 服务")
	if err := server.Serve(listen); err != nil {
		plog.Fatalf("启动服务失败，nest error: %v", err)
	}
}

// parseDuration parses lease argument string.
func parseDuration(lease string) time.Duration {
	dur, err := time.ParseDuration(lease)
	if err != nil {
		dur, err = time.ParseDuration(lease + "s")
	}
	if err != nil || dur < 0 {
		plog.Fatalf("invalid lease duration %s", lease)
	}
	return dur
}

func hasRootPrivilege() bool {
	return os.Geteuid() == 0
}

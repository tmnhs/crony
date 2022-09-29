package main

import (
	"fmt"
	"github.com/tmnhs/crony/admin/internal/handler"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/server"
	"os"
)

const (
	ServerName = "admin"
)

func main() {
	srv, err := server.NewApiServer(ServerName)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("new api server error:%s", err.Error()))
		os.Exit(1)
	}
	//logger.Infof("hello logger")
	//注册API路由业务
	srv.RegisterRouters(handler.RegisterRouters)
	//todo 合并到srv
	err = service.NewNodeWatcher().Watch()
	logger.GetLogger().Debug(fmt.Sprintf("resolver  error:%s", err))

	//todo 邮件相关操作

	//todo 定时清理日志
	err = srv.ListenAndServe()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("startup api server error:%v", err.Error()))
		os.Exit(1)
	}
	os.Exit(0)
}

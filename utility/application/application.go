package application

import (
	"babo/utility/pkg"
	"babo/utility/zlog"
	"os"
	"runtime"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type IApplication interface {
	LoadConfig(config_path string, config interface{}) error

	InitServices() error
}

func StartApp(app IApplication, app_name, config_path string, config interface{}, cancelPrint, closeDebug bool) error {
	zlog.Init(app_name, cancelPrint, false)
	defer zlog.Sync()
	defer pkg.ProtectError()

	if err := app.LoadConfig(config_path, config); err != nil {
		zlog.Error("Load config failed", zap.String("err", err.Error()))
		return err
	}

	if err := app.InitServices(); err != nil {
		zlog.Error("Init services failed", zap.String("err", err.Error()))
		return err
	}

	maxProcs := runtime.GOMAXPROCS(0)
	zlog.Info("GOMAXPROCS", zap.Int("num", maxProcs))
	zlog.Info("Server started succesfully.")

	//监听关闭信号
	pkg.WaitForTerminate()

	return nil
}

type DefaultApplication struct {
}

func (da *DefaultApplication) LoadConfig(config_path string, config interface{}) error {
	dataBytes, err := os.ReadFile(config_path)

	if err != nil {
		zlog.Error("Read yaml failed", zap.String("path", config_path))
		return err
	}

	err = yaml.Unmarshal(dataBytes, config)
	if err != nil {
		zlog.Error("umarshal yaml failed", zap.String("err", err.Error()))
		return err
	}
	//logins.Infof("yml config: %v", config)

	return nil
}

package main

import (
	"context"
	"flag"
	"gin-scaffold/internal/component_name"
	"gin-scaffold/internal/component_name/mysql"
	"gin-scaffold/internal/pkg/setting"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port      string
	runMode   string
	config    string
	isVersion bool
)

func init() {
	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}
	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	//err = setupLogger()
	//if err != nil {
	//	log.Fatalf("init.setupLogger err: %v", err)
	//}
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
	//err = setupValidator()
	//if err != nil {
	//	log.Fatalf("init.setupValidator err: %v", err)
	//}
	//err = setupTracer()
	//if err != nil {
	//	log.Fatalf("init.setupTracer err: %v", err)
	//}
}

func main() {
	gin.SetMode(component_name.ServerSetting.RunMode)
	router := component_name.NewRouter()
	s := &http.Server{
		Addr:           ":" + component_name.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    component_name.ServerSetting.ReadTimeout,
		WriteTimeout:   component_name.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}

func setupSetting() error {
	//s, err := setting.NewSetting(strings.Split(config, ",")...)
	//if err != nil {
	//	return err
	//}
	s, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = s.ReadSection("Server", &component_name.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &component_name.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &component_name.DatabaseSetting)
	if err != nil {
		return err
	}
	//err = s.ReadSection("JWT", &component_name.JWTSetting)
	//if err != nil {
	//	return err
	//}
	//err = s.ReadSection("Email", &component_name.EmailSetting)
	//if err != nil {
	//	return err
	//}

	component_name.AppSetting.DefaultContextTimeout *= time.Second
	//component_name.JWTSetting.Expire *= time.Second
	component_name.ServerSetting.ReadTimeout *= time.Second
	component_name.ServerSetting.WriteTimeout *= time.Second
	if port != "" {
		component_name.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		component_name.ServerSetting.RunMode = runMode
	}

	return nil
}

//func setupLogger() error {
//	fileName := component_name.AppSetting.LogSavePath + "/" + component_name.AppSetting.LogFileName + component_name.AppSetting.LogFileExt
//	component_name.Logger = component_name.NewLogger(&lumberjack.Logger{
//		Filename:  fileName,
//		MaxSize:   500,
//		MaxAge:    10,
//		LocalTime: true,
//	}, "", log.LstdFlags).WithCaller(2)
//
//	return nil
//}

func setupDBEngine() error {
	var err error
	component_name.DBEngine, err = mysql.NewDBEngine(component_name.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

//func setupValidator() error {
//	component_name.Validator = validator.NewCustomValidator()
//	component_name.Validator.Engine()
//	binding.Validator = component_name.Validator
//
//	return nil
//}
//
//func setupTracer() error {
//	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
//	if err != nil {
//		return err
//	}
//	component_name.Tracer = jaegerTracer
//	return nil
//}

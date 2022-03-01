package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/greendrop/todo-grpc-go-sample/domain/model"
	"github.com/greendrop/todo-grpc-go-sample/infrastructure/persistence"
	"github.com/greendrop/todo-grpc-go-sample/interface/server"
	proto_task_v1 "github.com/greendrop/todo-grpc-go-sample/proto/task/v1"
	usecase_task_v1 "github.com/greendrop/todo-grpc-go-sample/usecase/task/v1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	prepareLog()
	log.Info("Start server")

	appConfig, err := loadAppConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	persistence.AppConfig = appConfig

	gormDB, err := openDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer closeDatabase(gormDB)
	persistence.GormDB = gormDB

	grpcServer := prepareGrpcServer()
	go func() {
		err := serveGrpcServer(grpcServer)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	grpcGatewayServer, err := prepareGrpcGatewayServer()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = serveGrpcGatewayServer(grpcGatewayServer)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func prepareLog() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
}

func loadAppConfig() (*model.AppConfig, error) {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	configPath := "./configs/app_config." + appEnv + ".toml"
	var appConfig model.AppConfig
	_, err := toml.DecodeFile(configPath, &appConfig)
	if err != nil {
		return nil, err
	}

	appConfig.AppEnv = appEnv
	appConfig.Database.Url = os.Getenv("DATABASE_URL")

	return &appConfig, nil
}

func openDatabase() (*gorm.DB, error) {
	dsn := strings.Replace(persistence.AppConfig.Database.Url, "mysql://", "", 1)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	log.Info("Opened database")

	return gormDB, nil
}

func closeDatabase(gormDB *gorm.DB) {
	sqlDB, _ := gormDB.DB()

	if sqlDB != nil {
		sqlDB.Close()
		log.Info("Closed database")
	}
}

func prepareGrpcServer() *grpc.Server {
	logger := log.WithFields(log.Fields{})
	grpc_logrus.ReplaceGrpcLogger(logger)

	grpcLogrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logger, grpcLogrusOpts...),
			grpc_validator.UnaryServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logger, grpcLogrusOpts...),
			grpc_validator.StreamServerInterceptor(),
		),
	)

	taskPersistence := persistence.NewTaskPersistence()
	taskV1GetTaskListUseCase := usecase_task_v1.NewTaskV1GetTaskListUseCase(taskPersistence)
	taskV1GetTaskUseCase := usecase_task_v1.NewTaskV1GetTaskUseCase(taskPersistence)
	taskV1CreateTaskUseCase := usecase_task_v1.NewTaskV1CreateTaskUseCase(taskPersistence)
	taskV1UpdateTaskUseCase := usecase_task_v1.NewTaskV1UpdateTaskUseCase(taskPersistence)
	taskV1DeleteTaskUseCase := usecase_task_v1.NewTaskV1DeleteTaskUseCase(taskPersistence)
	taskV1ServiceServer := server.NewTaskV1ServiceServer(
		taskV1GetTaskListUseCase,
		taskV1GetTaskUseCase,
		taskV1CreateTaskUseCase,
		taskV1UpdateTaskUseCase,
		taskV1DeleteTaskUseCase,
	)
	proto_task_v1.RegisterTaskServiceServer(grpcServer, taskV1ServiceServer)

	return grpcServer
}

func serveGrpcServer(grpcServer *grpc.Server) error {
	listener, err := net.Listen("tcp", ":"+persistence.AppConfig.Server.GrpcPort)
	if err != nil {
		log.Errorf("Failed to listen:%v", err)
		return err
	}

	log.Infof("Serving gRPC at %v", listener.Addr())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Errorf("Failed to serve: %v", err)
		return err
	}

	return nil
}

func prepareGrpcGatewayServer() (*http.Server, error) {
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:"+persistence.AppConfig.Server.GrpcPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Errorf("Failed to dial server: %v", err)
		return nil, err
	}

	grpcGatewayMux := runtime.NewServeMux()
	err = proto_task_v1.RegisterTaskServiceHandler(context.Background(), grpcGatewayMux, conn)
	if err != nil {
		log.Errorf("Failed to register gateway: %v", err)
		return nil, err
	}

	grpcGatewayServer := &http.Server{
		Addr:    ":" + persistence.AppConfig.Server.GrpcGatewayPort,
		Handler: grpcGatewayMux,
	}
	return grpcGatewayServer, nil
}

func serveGrpcGatewayServer(grpcGatewayServer *http.Server) error {
	log.Info("Serving gRPC-Gateway on http://0.0.0.0:" + persistence.AppConfig.Server.GrpcGatewayPort)
	err := grpcGatewayServer.ListenAndServe()
	if err != nil {
		log.Errorf("Failed to serve: %v", err)
		return err
	}

	return nil
}

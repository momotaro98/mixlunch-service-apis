package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/momotaro98/mixlunch-service-api/logger"
	"github.com/momotaro98/mixlunch-service-api/partyservice"
	"github.com/momotaro98/mixlunch-service-api/pb"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	usService "github.com/momotaro98/mixlunch-service-api/userscheduleservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
)

const (
	ServiceAccountKeyPath = "./serviceAccount/serviceAccountKey.json"
)

func main() {
	var (
		gRPCAddr = flag.String("grpc", ":8081", "gRPC listen address")
	)
	flag.Parse()

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// Initial loading config
	var (
		logConf = &logger.Config{
			ErrorLevel: logger.Info, // Might need to be set from command argument
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE"),
		)
		usConf = &usService.Config{
			DSN: dsn,
		}
		pConf = &partyservice.Config{
			DSN:                   dsn,
			AppCredentialFilePath: ServiceAccountKeyPath,
		}
		tConf = &tagservice.Config{
			DSN: dsn,
		}
		uConf = &userservice.Config{
			DSN: dsn,
		}
	)

	// gPRC transport
	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		grpcServer := grpc.NewServer()
		pb.RegisterMixLunchServer(grpcServer, initializeGRPCServer(logConf, usConf, pConf, uConf, tConf))
		// Launch gRPC server
		log.Println("grpc:", *gRPCAddr)
		errChan <- grpcServer.Serve(listener)
	}()

	log.Fatal(<-errChan)
}

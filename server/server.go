package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	gRPC_task "rusprofile/proto"
	"rusprofile/server/internal"
	"syscall"
)

func main() {
	logrus.Println("Reading configs")

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	// environmental variables
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	logrus.Println("Starting Service...")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port")))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	gRPC_task.RegisterRusProfileServiceServer(s, &server{})
	reflection.Register(s)

	logrus.Println("Service started.")
	go func() {
		logrus.Println("Service is waiting for requests...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	<-ch
	// Stopping the server

	fmt.Println()

	logrus.Println("Stopping the server")
	s.Stop()
	logrus.Println("End of Program")
}

type server struct {
	gRPC_task.RusProfileServiceServer
}

// initialization configs for app
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// GetInfo searches for information about the company on the RusInfo site
func (s *server) GetInfo(ctx context.Context, req *gRPC_task.Request) (*gRPC_task.Response, error) {
	logrus.Printf("Request for INN: %v", req.INN)
	_, err := internal.INNValidation(req.INN)
	if err != nil {
		logrus.Errorf("Invalid INN: %s", err)
		return &gRPC_task.Response{CompanyInfo: nil}, err
	}
	info, err := internal.RusProfileParse(req.INN)
	if err != nil {
		logrus.Errorf("Info not found: %s", err)
		return &gRPC_task.Response{CompanyInfo: nil}, err
	}

	logrus.Printf("Info found for INN: %v", req.INN)
	return &gRPC_task.Response{CompanyInfo: &info}, nil
}

package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"rusprofile/internal/client/rusprofile"
	gRPC_task "rusprofile/proto"
	"syscall"
)

func main() {
	logrus.Println("Reading configs")
	config, err := parseConfig("./configs/config.yaml")
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	logrus.Println("Starting Service...")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go runRest(config.Server)
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	gRPC_task.RegisterRusProfileServiceServer(s, &server{rusProfileClient: rusprofile.NewClient(config.Client.Rusprofile.Host, config.Client.Rusprofile.Timeout)})
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

func runRest(config serverConfig) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gRPC_task.RegisterRusProfileServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%v:%v", config.Host, config.Port), opts)
	if err != nil {
		panic(err)
	}
	logrus.Printf("Server listening at %v", config.PortRest)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", config.PortRest), mux); err != nil {
		panic(err)
	}
}

type server struct {
	gRPC_task.RusProfileServiceServer
	rusProfileClient rusprofile.Client
}

// GetInfo searches for information about the company on the RusInfo site
func (s *server) GetInfo(ctx context.Context, req *gRPC_task.Request) (*gRPC_task.Response, error) {
	logrus.Printf("Request for INN: %v", req.INN)

	info, err := s.rusProfileClient.GetProfile(ctx, req.INN)
	if err != nil {
		logrus.Errorf("%s", err)
		return nil, err
	}

	logrus.Printf("Info found for INN: %v", req.INN)
	return &gRPC_task.Response{CompanyInfo: &info}, nil
}

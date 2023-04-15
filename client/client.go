package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	gRPC_task "rusprofile/proto"
	"strconv"
	"sync"
)

func main() {

	cc, client := startingServer()
	defer cc.Close()

	args, err := flagsAndCommands()

	if err != nil {
		return
	}

	wg := &sync.WaitGroup{}
	getInfoRequest(client, args, wg)
	wg.Wait()
}

// initialization of config
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// starting the server
func startingServer() (*grpc.ClientConn, gRPC_task.RusProfileServiceClient) {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	cc, err := grpc.Dial(fmt.Sprintf("localhost:%s", viper.GetString("port")), opts)
	if err != nil {
		logrus.Fatalf("could not connect: %v", err)
	}

	return cc, gRPC_task.NewRusProfileServiceClient(cc)
}

// parsing flags, commands and parameters
func flagsAndCommands() ([]string, error) {
	help := flag.Bool("help", false, "help message")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageMessage)
	}

	flag.Parse()

	if len(os.Args) == 1 {
		write(helpMessage)
		return os.Args, errors.New("parameter not found")
	}

	if *help {
		write(helpMessage)
		return os.Args, errors.New("help message only")
	}

	return os.Args, nil
}

// making request
func getInfoRequest(client gRPC_task.RusProfileServiceClient, args []string, wg *sync.WaitGroup) {
	var argI = 1
	// URLs from file (command "file")
	if args[argI] == "file" {
		file, err := os.Open(args[argI+1])
		if err != nil {
			log.Fatal("Can't open file")
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			INN, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				log.Fatal("INN isn't digit")
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := client.GetInfo(context.Background(), &gRPC_task.Request{INN: INN})
				if err != nil {
					fmt.Printf("Unexpected error: %v\n", err)
					return
				}
				write(fmt.Sprintf("Company name: %v\n INN: %v\n CPP: %v\n Director %v\n", res.CompanyInfo.Name, res.CompanyInfo.INN, res.CompanyInfo.KPP, res.CompanyInfo.Director))
			}()
		}
		return
	} else { // URLs from parameters
		for i := argI; i < len(os.Args); i++ {
			INN, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				log.Fatal("INN isn't digit")
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := client.GetInfo(context.Background(), &gRPC_task.Request{INN: INN})
				if err != nil {
					fmt.Printf("Unexpected error: %v\n", err)
					return
				}
				write(fmt.Sprintf("Company name: %v\n INN: %v\n CPP: %v\n Director %v\n", res.CompanyInfo.Name, res.CompanyInfo.INN, res.CompanyInfo.KPP, res.CompanyInfo.Director))
			}()
		}
	}
}

// writing message to console or other destination
func write(message string) {
	dst := os.Stdout
	fmt.Fprint(dst, message)
}

const usageMessage = `Please try to input INN as a parameter:
go run client/client.go [1234567891]
or use flag --help`

const helpMessage = `
RusProfile gRPC Parser:
Description:
RusProfile gRPC Parser is a CLI tool to get company info by its URL.
You can input several INNs divided by backspaces or load them from file.
  Usage: 
client/client.go [flags] INNs 
client/client.go [flags] file name.ext
  Usage Examples:
go run client/client.go 1234567891
go run client/client.go file companies.txt
  Commands:
	
file name.ext
  Flags:
    --help     Show this help message`

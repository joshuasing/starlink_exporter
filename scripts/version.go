package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/joshuasing/starlink_exporter/internal/exporter"
	"github.com/joshuasing/starlink_exporter/internal/spacex/api/device"
)

var dishAddress = flag.String("dish", exporter.DefaultDishAddress, "Dish address")

func main() {
	flag.Parse()
	if *dishAddress == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	conn, err := grpc.NewClient(*dishAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := device.NewDeviceClient(conn)
	res, err := client.Handle(ctx, &device.Request{
		Request: new(device.Request_GetStatus),
	})
	if err != nil {
		log.Fatal(err)
	}

	deviceStatus := res.GetDishGetStatus()
	deviceInfo := deviceStatus.GetDeviceInfo()

	fmt.Printf("Starlink Dishy %s\n", deviceInfo.GetHardwareVersion())
	fmt.Printf("Software version %s\n", deviceInfo.GetSoftwareVersion())
	fmt.Printf("API version %d\n", res.GetApiVersion())
}

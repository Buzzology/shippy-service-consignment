package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// This is the generated protobuf code
	pb "github.com/Buzzology/shippy-service-consignment/proto/consignment"
	vesselPb "github.com/Buzzology/shippy-service-vessel/proto/vessel"
	micro "github.com/micro/go-micro/v2"
)

const defaultHost = "mongodb://localhost:27017"

func main() {

	// Creates the new service (can include options here)
	service := micro.NewService(

		// This name must match the package name given in protobuf definition
		micro.Name("shippy.service.consignment"),
	)

	service.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	// Create a mongo db connection
	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}

	defer client.Disconnect(context.Background())

	// Prepare repository and handlers
	consignmentCollection := client.Database("shippy").Collection("consignments")
	repository := &MongoRepository{consignmentCollection}
	vesselClient := vesselPb.NewVesselService("shippy.service.vessel", service.Client())
	h := &handler{repository, vesselClient}

	// Register handlers
	pb.RegisterShippingServiceHandler(service.Server(), h)

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}

	// Run the server
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
package main

import (
	"context"
	pb "github.com/buzzology/go-microservices-tutorial/shippy-service-consignment/proto/consignment"
	vesselPb "github.com/buzzology/go-microservices-tutorial/shippy-service-vessel/proto/vessel"
	"github.com/pkg/errors"
)

type handler struct {
	repository
	vesselClient vesselPb.VesselService
}

// We created only one method on our service which is create. Is handled by gRPC server
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Call a client instance of our vessel service with our consignment details
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselPb.Specification{
		Capacity: int32(len(req.Containers)),
		MaxWeight: req.Weight,
	})

	if vesselResponse == nil {
		return errors.New("Error fetching vessel, returned nil")
	}

	if err != nil {
		return err
	}

	// Set the vessel for the consignment
	req.VesselId = vesselResponse.Vessel.Id

	// Save
	if err = s.repository.Create(ctx, MarshalConsignment(req)); err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req
	return nil
}


func(s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {

	consignments, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}

	res.Consignments = UnmarshConsignmentCollection(consignments)
	return nil
}
package main

import (
	"context"
	pb "github.com/Buzzology/shippy-service-consignment/proto/consignment"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"fmt"
)

type Consignment struct {
	ID          string     `json:"id"`
	Weight      int32      `json:"weight"`
	Description string     `json:"description"`
	Containers  Containers `json:"containers"`
	VesselID    string     `json:"vessel_id"`
}

type Container struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	UserID     string `json:"user_id"`
}

type Containers []*Container

func MarshalContainerCollection(containers []*pb.Container) []*Container {
	collection := make([]*Container, 0)
	for _, container := range containers {
		collection = append(collection, MarshalContainer(container))
	}
	return collection
}

func UnmarshalContainerCollection(containers []*Container) []*pb.Container {
	collection := make([]*pb.Container, 0)
	for _, container := range containers {
		collection = append(collection, UnmarshalContainer(container))
	}

	return collection
}

func MarshalContainer(container *pb.Container) *Container {
	return &Container{
		ID:         container.Id,
		CustomerID: container.CustomerId,
		UserID:     container.UserId,
	}
}

func MarshalConsignment(consignment *pb.Consignment) *Consignment {
	containers := MarshalContainerCollection(consignment.Containers)

	return &Consignment{
		ID:          consignment.Id,
		Weight:      consignment.Weight,
		Description: consignment.Description,
		Containers:  containers,
		VesselID:    consignment.VesselId,
	}
}

func UnmarshalConsignment(consignment *Consignment) *pb.Consignment {

	containers := UnmarshalContainerCollection(consignment.Containers)

	return &pb.Consignment{
		Id:          consignment.ID,
		Description: consignment.Description,
		Weight:      consignment.Weight,
		Containers:  containers,
		VesselId:    consignment.VesselID,
	}
}

func UnmarshConsignmentCollection(consignments []*Consignment) []*pb.Consignment {
	collection := make([]*pb.Consignment, 0)
	for _, consignment := range consignments {
		collection = append(collection, UnmarshalConsignment(consignment))
	}

	return collection
}

func UnmarshalContainer(container *Container) *pb.Container {
	return &pb.Container{
		Id:         container.ID,
		CustomerId: container.CustomerID,
		UserId:     container.UserID,
	}
}

type repository interface {
	Create(ctx context.Context, consignment *Consignment) error
	GetAll(ctx context.Context) ([]*Consignment, error)
}

// MongoRepository implementation
type MongoRepository struct {
	collection *mongo.Collection
}

func (repository *MongoRepository) Create(ctx context.Context, consignment *Consignment) error {
	_, err := repository.collection.InsertOne(ctx, consignment)
	return err
}

func (repository *MongoRepository) GetAll(ctx context.Context) ([]*Consignment, error) {
	cursor, err := repository.collection.Find(ctx, bson.D{{}})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var consignments []*Consignment
	for cursor.Next(ctx) {
		var consignment *Consignment
		if err := cursor.Decode(&consignment); err != nil {
			return nil, err
		}
		consignments = append(consignments, consignment)
	}

	return consignments, err
}

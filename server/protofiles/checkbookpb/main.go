package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	checkbookpb "projet_grpc/server/protofiles/checkbook"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	checkbookpb.UnimplementedCheckbookServiceServer
}

func (s server) CreateCheckbook(ctx context.Context, req *checkbookpb.CheckbookRequest) (*checkbookpb.CheckbookResponse, error) {
	log.Println("Number Of Pages:", req.NbPage)
	log.Println("Account ID:", req.AccountId)
	creationDate := timestamppb.Now()
	log.Println("Creation Date:", creationDate)
	var checkbook_id int32 = int32(rand.Intn(1000))
	log.Println("Checkbook Id:", checkbook_id)

	return &checkbookpb.CheckbookResponse{NbPage: req.NbPage, AccountId: req.AccountId, CreationDate: creationDate, Id: checkbook_id}, nil
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		panic(err)
	}

	fmt.Println("starting server")
	s := grpc.NewServer()
	checkbookpb.RegisterCheckbookServiceServer(s, server{})

	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}

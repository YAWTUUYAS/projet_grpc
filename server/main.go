package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	checkbookpb "projet_grpc/protofiles/checkbookpb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	checkbookpb.UnimplementedCheckbookServiceServer
}

func getNbPages(p checkbookpb.Pages) int32 {
	switch p {
	case checkbookpb.Pages_TWENTY_FIVE:
		return 25
	case checkbookpb.Pages_FIFTY:
		return 50
	default:
		return 0
	}
}

func (s server) CreateCheckbook(ctx context.Context, req *checkbookpb.CheckbookRequest) (*checkbookpb.CheckbookResponse, error) {

	pagesEnum := req.GetNbPage()
	nbPages := getNbPages(pagesEnum)

	log.Println("Pages enum:", pagesEnum)
	log.Println("Number Of Pages:", nbPages)
	log.Println("Account ID:", req.GetAccountId())

	creationDate := timestamppb.Now()
	log.Println("Creation Date:", creationDate)

	checkbookID := int32(rand.Intn(1000))
	log.Println("Checkbook Id:", checkbookID)

	return &checkbookpb.CheckbookResponse{NbPage: pagesEnum, AccountId: req.GetAccountId(), CreationDate: creationDate, Id: checkbookID}, nil
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

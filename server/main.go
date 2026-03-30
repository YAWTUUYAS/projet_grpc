package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	checkbookpb "projet_grpc/protofiles/checkbookpb/v1"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	checkbookpb.UnimplementedCheckbookServiceServer
	db map[int32]*checkbookpb.CreateCheckbookResponse
}

func getNbPages(p checkbookpb.Pages) int32 {
	switch p {
	case checkbookpb.Pages_PAGES_TWENTY_FIVE:
		return 25
	case checkbookpb.Pages_PAGES_FIFTY:
		return 50
	default:
		return 0
	}
}

func (s *server) CreateCheckbook(ctx context.Context, req *checkbookpb.CreateCheckbookRequest) (*checkbookpb.CreateCheckbookResponse, error) {

	pagesEnum := req.GetNbPage()
	nbPages := getNbPages(pagesEnum)

	log.Println("Pages enum:", pagesEnum)
	log.Println("Number Of Pages:", nbPages)
	log.Println("Account ID:", req.GetAccountId())

	creationDate := timestamppb.Now()
	log.Println("Creation Date:", creationDate)

	checkbookID := int32(rand.Intn(10))
	log.Println("Checkbook Id:", checkbookID)

	checkbook := &checkbookpb.CreateCheckbookResponse{NbPage: pagesEnum, AccountId: req.GetAccountId(), CreationDate: creationDate, Id: checkbookID}

	s.db[checkbookID] = checkbook

	log.Println("Saved in db: ", checkbookID)

	return checkbook, nil
}

func (s *server) GetCheckbooks(ctx context.Context, req *checkbookpb.GetCheckbooksRequest) (*checkbookpb.GetCheckbooksResponse, error) {

	accountId := req.GetAccountId()
	var results []*checkbookpb.CreateCheckbookResponse

	for _, checkbook := range s.db {
		if checkbook.GetAccountId() == accountId {
			results = append(results, checkbook)
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no checkbooks found for account: %s", accountId)
	}

	log.Println("Found", len(results), "checkbooks for account", accountId)

	return &checkbookpb.GetCheckbooksResponse{
		Checkbooks: results,
	}, nil
}

func (s *server) UpdateCheckbook(ctx context.Context, req *checkbookpb.UpdateCheckbookRequest) (*checkbookpb.UpdateCheckbookResponse, error) {
	accountId := req.GetAccountId()
	nbPages := req.GetNbPage()
	checkbookId := req.GetId()

	for _, checkbook := range s.db {
		if checkbook.GetId() == checkbookId {
			if checkbook.GetAccountId() != accountId {
				s.db[checkbookId].AccountId = accountId

			}
			if checkbook.GetNbPage() != nbPages {
				s.db[checkbookId].NbPage = nbPages

			}
			updatedcheckbook := &checkbookpb.UpdateCheckbookResponse{NbPage: nbPages, AccountId: accountId, CreationDate: checkbook.CreationDate, Id: checkbookId}
			log.Println("Updated checkbook for id", checkbookId)
			return updatedcheckbook, nil
		}
	}

	return nil, fmt.Errorf("no checkbooks found for id: %d or no need to update", checkbookId)

}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		panic(err)
	}

	fmt.Println("starting server")
	s := grpc.NewServer()

	srv := &server{
		db: make(map[int32]*checkbookpb.CreateCheckbookResponse),
	}

	checkbookpb.RegisterCheckbookServiceServer(s, srv)

	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}

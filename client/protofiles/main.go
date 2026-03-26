package main

import (
	"context"
	"log"

	checkbookpb "projet_grpc/server/protofiles/checkbook"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	defer cc.Close()
	c := checkbookpb.NewCheckbookServiceClient(cc)

	createCheckbook(25, "account_001", "25/03/2026", c)
	createCheckbook(50, "account_002", "26/03/2026", c)
}

func createCheckbook(nbPage int32, accountId, creationDate string, c checkbookpb.CheckbookServiceClient) {
	log.Println("creating checkbook")

	res, err := c.CreateCheckbook(context.Background(), &checkbookpb.CheckbookRequest{NbPage: nbPage, AccountId: accountId, CreationDate: creationDate})

	if err != nil {
		log.Println("error: ", err)
		panic(err)
	}

	log.Println(res.AccountId, res.CreationDate, res.Id, res.NbPage)

}

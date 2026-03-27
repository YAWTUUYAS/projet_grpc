package main

import (
	"context"
	"log"

	checkbookpb "projet_grpc/protofiles/checkbookpb"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	defer cc.Close()
	c := checkbookpb.NewCheckbookServiceClient(cc)

	createCheckbook(checkbookpb.Pages_TWENTY_FIVE, "account_001", c)
	createCheckbook(checkbookpb.Pages_FIFTY, "account_002", c)
}

func createCheckbook(nbPage checkbookpb.Pages, accountId string, c checkbookpb.CheckbookServiceClient) {
	log.Println("creating checkbook")

	res, err := c.CreateCheckbook(context.Background(), &checkbookpb.CheckbookRequest{
		NbPage:    nbPage,
		AccountId: accountId,
	})

	if err != nil {
		log.Println("error: ", err)
		panic(err)
	}

	log.Println(res.AccountId, res.CreationDate, res.Id, res.NbPage)
}

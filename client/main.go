package main

import (
	"context"
	"log"

	checkbookpb "projet_grpc/protofiles/checkbookpb/v1"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	defer cc.Close()

	c := checkbookpb.NewCheckbookServiceClient(cc)

	createCheckbook(checkbookpb.Pages_PAGES_TWENTY_FIVE, "account_001", c)
	createCheckbook(checkbookpb.Pages_PAGES_FIFTY, "account_001", c)

	getCheckbooks("account_001", c)

	updateCheckbook(9, "account_002", checkbookpb.Pages_PAGES_FIFTY, c)
	getCheckbooks("account_002", c)
}

func createCheckbook(nbPage checkbookpb.Pages, accountId string, c checkbookpb.CheckbookServiceClient) {
	log.Println("creating checkbook")

	res, err := c.CreateCheckbook(context.Background(), &checkbookpb.CreateCheckbookRequest{
		NbPage:    nbPage,
		AccountId: accountId,
	})

	if err != nil {
		log.Println("error: ", err)
		panic(err)
	}

	log.Println(res.AccountId, res.CreationDate, res.Id, res.NbPage)
}

func getCheckbooks(accountId string, c checkbookpb.CheckbookServiceClient) {
	log.Println("fetching checkbooks for:", accountId)

	res, err := c.GetCheckbooks(context.Background(), &checkbookpb.GetCheckbooksRequest{
		AccountId: accountId,
	})

	if err != nil {
		log.Println("error:", err)
		return
	}

	for _, cb := range res.Checkbooks {
		log.Println("ID:", cb.Id,
			"Pages:", cb.NbPage,
			"Date:", cb.CreationDate)
	}
}

func updateCheckbook(checkbookId int32, accountId string, nbPage checkbookpb.Pages, c checkbookpb.CheckbookServiceClient) {
	log.Println("updating checkbook for id: ", checkbookId)

	res, err := c.UpdateCheckbook(context.Background(), &checkbookpb.UpdateCheckbookRequest{
		Id: checkbookId, AccountId: accountId, NbPage: nbPage,
	})

	if err != nil {
		log.Println("error:", err)
		return
	}

	log.Println(res.AccountId, res.Id, res.NbPage)
}

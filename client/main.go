package main

import (
	"context"
	"fmt"
	"log"

	checkbookpb "projet_grpc/protofiles/checkbookpb/v1"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	// var wg sync.WaitGroup
	eg, _ := errgroup.WithContext(context.Background())

	if err != nil {
		panic(err)
	}

	defer cc.Close()

	c := checkbookpb.NewCheckbookServiceClient(cc)
	/**
	wg.Add(3)
	go createCheckbook(checkbookpb.Pages_PAGES_TWENTY_FIVE, "account_001", c)
	go createCheckbook(checkbookpb.Pages_PAGES_FIFTY, "account_001", c)
	go createCheckbook(checkbookpb.Pages_PAGES_FIFTY, "account_002", c)
	wg.Wait()

	wg.Add(2)
	go getCheckbooks("account_001", c)
	go getCheckbooks("account_002", c)
	wg.Done()

	updateCheckbook(9, "account_002", checkbookpb.Pages_PAGES_FIFTY, c)
	**/

	eg, ctx := errgroup.WithContext(context.Background())
	eg.SetLimit(3)
	eg.Go(func() error {
		return createCheckbook(checkbookpb.Pages_PAGES_TWENTY_FIVE, "account_001", c)
	})
	eg.Go(func() error {
		return createCheckbook(checkbookpb.Pages_PAGES_FIFTY, "account_001", c)
	})
	eg.Go(func() error {
		return createCheckbook(checkbookpb.Pages_PAGES_FIFTY, "account_002", c)
	})
	if err := eg.Wait(); err != nil {
		log.Println("Create error:", err)
		return
	}

	eg2, ctx := errgroup.WithContext(ctx)
	eg2.SetLimit(2)
	eg2.Go(func() error {
		return getCheckbooks("account_001", c)
	})
	eg2.Go(func() error {
		return getCheckbooks("account_002", c)
	})
	if err := eg2.Wait(); err != nil {
		log.Println("Get error:", err)
		return
	}
	fmt.Printf("all go routines completed")

}

func createCheckbook(nbPage checkbookpb.Pages, accountId string, c checkbookpb.CheckbookServiceClient) error {
	log.Println("creating checkbook")

	res, err := c.CreateCheckbook(context.Background(), &checkbookpb.CreateCheckbookRequest{
		NbPage:    nbPage,
		AccountId: accountId,
	})

	if err != nil {
		log.Println("error: ", err)
		return err
	}

	log.Println(res.AccountId, res.CreationDate, res.Id, res.NbPage)
	return nil
}

func getCheckbooks(accountId string, c checkbookpb.CheckbookServiceClient) error {
	log.Println("fetching checkbooks for:", accountId)

	res, err := c.GetCheckbooks(context.Background(), &checkbookpb.GetCheckbooksRequest{
		AccountId: accountId,
	})

	if err != nil {
		log.Println("error:", err)
		return err
	}

	for _, cb := range res.Checkbooks {
		log.Println("ID:", cb.Id,
			"Pages:", cb.NbPage,
			"Date:", cb.CreationDate)
	}
	return nil
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

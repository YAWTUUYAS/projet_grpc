package main

import (
	"context"
	checkbookpb "projet_grpc/protofiles/checkbookpb/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ctx context.Context = context.Background()

func TestGetNbPages(t *testing.T) {

	result1, err := getNbPages(checkbookpb.Pages_PAGES_TWENTY_FIVE)

	if result1 != 25 {
		t.Errorf("expected 25, got %d", result1)
	}

	require.NoError(t, err)

	result2, err := getNbPages(checkbookpb.Pages_PAGES_FIFTY)
	if result2 != 50 {
		t.Errorf("expected 50, got %d", result2)
	}
	require.NoError(t, err)

	result3, err := getNbPages(checkbookpb.Pages_PAGES_UNSPECIFIED)
	if result3 != 0 {
		t.Errorf("expected 0, got %d", result3)
	}
	require.Error(t, err)
}

func TestCreateCheckbook(t *testing.T) {
	s := &server{
		db: make(map[int32]*checkbookpb.CreateCheckbookResponse),
	}
	req := &checkbookpb.CreateCheckbookRequest{NbPage: checkbookpb.Pages_PAGES_FIFTY, AccountId: "account_test"}

	res, err := s.CreateCheckbook(ctx, req)

	if err != nil {
		t.Errorf("could not create a checkbook request")
	}
	for _, tt := range s.db {
		if tt.AccountId != res.AccountId {
			t.Errorf("Account Id does not correspond to db")
		}
		if tt.NbPage != res.NbPage {
			t.Errorf("Page Number does not correspond to db")
		}
	}
}

func TestGetCheckbooks(t *testing.T) {
	s := &server{
		db: make(map[int32]*checkbookpb.CreateCheckbookResponse),
	}

	checkbook1 := &checkbookpb.CreateCheckbookResponse{NbPage: checkbookpb.Pages_PAGES_TWENTY_FIVE, AccountId: "Account_test1", CreationDate: timestamppb.Now(), Id: 001}
	s.db[001] = checkbook1
	checkbook2 := &checkbookpb.CreateCheckbookResponse{NbPage: checkbookpb.Pages_PAGES_FIFTY, AccountId: "Account_test2", CreationDate: timestamppb.Now(), Id: 002}
	s.db[002] = checkbook2
	checkbook3 := &checkbookpb.CreateCheckbookResponse{NbPage: checkbookpb.Pages_PAGES_FIFTY, AccountId: "Account_test1", CreationDate: timestamppb.Now(), Id: 003}
	s.db[003] = checkbook3

	accountIdTest := "Account_test1"

	res, err := s.GetCheckbooks(ctx, &checkbookpb.GetCheckbooksRequest{AccountId: accountIdTest})

	if err != nil {
		t.Errorf("could not fetch a checkbook request with Account Id: %s", accountIdTest)
	}

	if len(res.Checkbooks) != 2 {
		t.Errorf("did not fetch the right checkbook request(s) with Account Id: %s", accountIdTest)
	}
}

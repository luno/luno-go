package main

import (
	"context"
	"flag"
	"log"
	"time"

	luno "github.com/luno/luno-go"
)

var (
	apiKeyID     = flag.String("api_key_id", "", "Luno API key ID")
	apiKeySecret = flag.String("api_key_secret", "", "Luno API key secret")
	debug        = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()

	cl := luno.NewClient()
	cl.SetDebug(*debug)
	cl.SetAuth(*apiKeyID, *apiKeySecret)

	ctx := context.Background()

	{
		req := luno.GetTickersRequest{}
		res, err := cl.GetTickers(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	{
		req := luno.GetTickerRequest{Pair: "XBTZAR"}
		res, err := cl.GetTicker(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	{
		req := luno.GetOrderBookRequest{Pair: "XBTZAR"}
		res, err := cl.GetOrderBook(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	{
		req := luno.ListTradesRequest{
			Pair:  "XBTZAR",
			Since: time.Now().Add(-24*time.Hour).UnixNano() / 1e6,
		}
		res, err := cl.ListTrades(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	var accountID string

	{
		req := luno.GetBalancesRequest{}
		res, err := cl.GetBalances(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		if res != nil && len(res.Balance) > 0 {
			accountID = res.Balance[0].AccountId
		}
		time.Sleep(500 * time.Millisecond)
	}

	if accountID != "" {
		req := luno.ListTransactionsRequest{
			Id:     accountID,
			MinRow: 1,
			MaxRow: 1000,
		}
		res, err := cl.ListTransactions(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	if accountID != "" {
		req := luno.ListPendingTransactionsRequest{
			Id: accountID,
		}
		res, err := cl.ListPendingTransactions(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	{
		req := luno.ListOrdersRequest{}
		res, err := cl.ListOrders(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	{
		req := luno.ListUserTradesRequest{}
		res, err := cl.ListUserTrades(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	{
		req := luno.GetFeeInfoRequest{}
		res, err := cl.GetFeeInfo(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	{
		req := luno.GetFundingAddressRequest{
			Asset: "XBT",
		}
		res, err := cl.GetFundingAddress(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}

	var withdrawalID string

	{
		req := luno.ListWithdrawalsRequest{}
		res, err := cl.ListWithdrawals(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		if res != nil && len(res.Withdrawals) > 0 {
			withdrawalID = res.Withdrawals[0].Id
		}
		time.Sleep(500 * time.Millisecond)
	}

	if withdrawalID != "" {
		req := luno.GetWithdrawalRequest{
			Id: withdrawalID,
		}
		res, err := cl.GetWithdrawal(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}
}

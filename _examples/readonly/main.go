package main

import (
	"context"
	"flag"
	"log"
	"strconv"
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
			Since: luno.Time(time.Now().Add(-24 * time.Hour)),
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
		aid, _ := strconv.ParseInt(accountID, 10, 64)
		req := luno.ListTransactionsRequest{
			Id:     aid,
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
		aid, _ := strconv.ParseInt(accountID, 10, 64)
		req := luno.ListPendingTransactionsRequest{
			Id: aid,
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
		wid, _ := strconv.ParseInt(withdrawalID, 10, 64)
		req := luno.GetWithdrawalRequest{
			Id: wid,
		}
		res, err := cl.GetWithdrawal(ctx, &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", res)
		time.Sleep(500 * time.Millisecond)
	}
}

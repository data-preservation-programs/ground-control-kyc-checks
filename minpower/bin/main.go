package main

import (
	"context"
	"log"
	"net/http"

	"github.com/filecoin-project/go-address"
	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

func main() {
	headers := http.Header{}
	api_addr := "api.chain.love"

	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(),
		"wss://"+api_addr+"/rpc/v0", "Filecoin",
		[]interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	addr, err := address.NewFromString("f02620")
	if err != nil {
		log.Fatalf("address error: %s", err)
	}

	// Now you can call any API you're interested in.
	power, err := api.StateMinerPower(context.Background(), addr, types.EmptyTSK)
	if err != nil {
		log.Fatalf("calling state miner power: %s", err)
	}
	log.Println(power)
}

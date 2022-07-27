package minpower

import (
	"context"
	"log"
	"net/http"

	"github.com/filecoin-project/go-address"
	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

// ReverseRunes returns its argument string reversed rune-wise left to right.
func ReverseRunes(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// LookupPower gets the power for the miner from the Lotus API
func LookupPower(ctx context.Context, miner string) (*lotusapi.MinerPower, error) {
	headers := http.Header{}
	api_addr := "api.chain.love"

	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(ctx,
		"wss://"+api_addr+"/rpc/v0", "Filecoin",
		[]interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	addr, err := address.NewFromString(miner)
	if err != nil {
		return nil, err
	}

	power, err := api.StateMinerPower(ctx, addr, types.EmptyTSK)
	if err != nil {
		return nil, err
	}
	log.Printf("Miner power %s: %v\n", miner, power)
	return power, nil
}

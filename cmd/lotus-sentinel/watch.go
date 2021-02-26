package main

import (
	"context"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/urfave/cli/v2"

	cli2 "github.com/filecoin-project/lotus/cli"
	api2 "github.com/filecoin-project/lotus/cmd/lotus-sentinel/api"
	"github.com/filecoin-project/lotus/node/repo"
)

var sentinelStartWatchCmd = &cli.Command{
	Name:  "watch",
	Usage: "start a watch against the chain",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name: "confidence",
		},
	},
	Action: func(cctx *cli.Context) error {
		apic, closer, err := GetSentinelNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := cli2.ReqContext(cctx)

		//confidence := abi.ChainEpoch(cctx.Int64("confidence"))

		if err := apic.SentinelWatchStart(ctx); err != nil {
			return err
		}

		return nil
	},
}

func GetSentinelNodeAPI(ctx *cli.Context) (api2.SentinelNode, jsonrpc.ClientCloser, error) {
	addr, headers, err := cli2.GetRawAPI(ctx, repo.FullNode)
	if err != nil {
		return nil, nil, err
	}

	return NewSentinelNodeRPC(ctx.Context, addr, headers)
}

func NewSentinelNodeRPC(ctx context.Context, addr string, requestHeader http.Header) (api2.SentinelNode, jsonrpc.ClientCloser, error) {
	var res api2.SentinelNodeStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.CommonStruct.Internal,
			&res.FullNodeStruct.Internal,
			&res.Internal,
		},
		requestHeader,
	)
	return &res, closer, err
}

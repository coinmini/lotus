package api

import (
	"context"

	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/events"
	"github.com/filecoin-project/lotus/node/impl"
	"github.com/filecoin-project/lotus/sentinel"
)

var log = logging.Logger("sentinel")

type SentinelNode interface {
	api.FullNode
	SentinelWatchStart(ctx context.Context) error
}

type SentinelNodeAPI struct {
	fx.In

	impl.FullNodeAPI
	Events *events.Events
}

func (m *SentinelNodeAPI) SentinelWatchStart(ctx context.Context) error {
	log.Info("starting sentinel watch")
	return m.Events.Observe(&sentinel.LoggingTipSetObserver{})
}

var _ SentinelNode = &SentinelNodeAPI{}

type SentinelNodeStruct struct {
	apistruct.FullNodeStruct

	Internal struct {
		SentinelWatchStart func(context.Context) error `perm:"read"`
	}
}

func (s *SentinelNodeStruct) SentinelWatchStart(ctx context.Context) error {
	return s.Internal.SentinelWatchStart(ctx)
}

var _ SentinelNode = &SentinelNodeStruct{}

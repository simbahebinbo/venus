package blocksub

import (
	"bytes"
	"context"

	"github.com/ipfs/go-log/v2"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/ipfs-force-community/metrics"
)

var (
	blockTopicLogger = log.Logger("net/block_validator")
	mDecodeBlkFail   = metrics.NewCounter("net/pubsub_block_decode_failure", "Number of blocks that fail to decode seen on block pubsub channel")
)

// BlockTopicValidator may be registered on go-libp2p-pubsub to validate blocksub messages.
type BlockTopicValidator struct {
	validator pubsub.ValidatorEx
	opts      []pubsub.ValidatorOpt
}

type BlockHeaderValidator interface {
	ValidateBlockMsg(context.Context, *types.BlockMsg) pubsub.ValidationResult
}

// NewBlockTopicValidator retruns a BlockTopicValidator using `bv` for message validation
func NewBlockTopicValidator(bv BlockHeaderValidator, opts ...pubsub.ValidatorOpt) *BlockTopicValidator {
	return &BlockTopicValidator{
		opts: opts,
		validator: func(ctx context.Context, p peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
			var bm types.BlockMsg
			err := bm.UnmarshalCBOR(bytes.NewReader(msg.GetData()))
			if err != nil {
				blockTopicLogger.Warnf("failed to decode blocksub payload from peer %s: %s", p.String(), err.Error())
				mDecodeBlkFail.Tick(ctx)
				return pubsub.ValidationIgnore
			}

			validateResult := bv.ValidateBlockMsg(ctx, &bm)
			if validateResult == pubsub.ValidationAccept {
				msg.ValidatorData = bm
			}
			return validateResult
		},
	}
}

func (btv *BlockTopicValidator) Topic(network string) string {
	return types.BlockTopic(network)
}

func (btv *BlockTopicValidator) Validator() pubsub.ValidatorEx {
	return btv.validator
}

func (btv *BlockTopicValidator) Opts() []pubsub.ValidatorOpt {
	return btv.opts
}

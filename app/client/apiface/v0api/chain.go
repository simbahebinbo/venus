package v0api

import (
	"context"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	acrypto "github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/dline"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/venus/venus-shared/actors/builtin/miner"
	apitypes "github.com/filecoin-project/venus/venus-shared/api/chain"
	types "github.com/filecoin-project/venus/venus-shared/chain"
)

type IChain interface {
	IAccount
	IActor
	IBeacon
	IMinerState
	IChainInfo
}

type IAccount interface {
	// Rule[perm:read]
	StateAccountKey(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error)
}

type IActor interface {
	// Rule[perm:read]
	StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error)
	// Rule[perm:read]
	ListActor(ctx context.Context) (map[address.Address]*types.Actor, error)
}

type IBeacon interface {
	// Rule[perm:read]
	BeaconGetEntry(ctx context.Context, epoch abi.ChainEpoch) (*types.BeaconEntry, error)
}

type IChainInfo interface {
	// Rule[perm:read]
	BlockTime(ctx context.Context) time.Duration
	// Rule[perm:read]
	ChainList(ctx context.Context, tsKey types.TipSetKey, count int) ([]types.TipSetKey, error)
	// Rule[perm:read]
	ChainHead(ctx context.Context) (*types.TipSet, error)
	// Rule[perm:admin]
	ChainSetHead(ctx context.Context, key types.TipSetKey) error
	// Rule[perm:read]
	ChainGetTipSet(ctx context.Context, key types.TipSetKey) (*types.TipSet, error)
	// Rule[perm:read]
	ChainGetTipSetByHeight(ctx context.Context, height abi.ChainEpoch, tsk types.TipSetKey) (*types.TipSet, error)
	// Rule[perm:read]
	ChainGetRandomnessFromBeacon(ctx context.Context, key types.TipSetKey, personalization acrypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) (abi.Randomness, error)
	// Rule[perm:read]
	ChainGetRandomnessFromTickets(ctx context.Context, tsk types.TipSetKey, personalization acrypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) (abi.Randomness, error)
	// Rule[perm:read]
	ChainGetBlock(ctx context.Context, id cid.Cid) (*types.BlockHeader, error)
	// Rule[perm:read]
	ChainGetMessage(ctx context.Context, msgID cid.Cid) (*types.Message, error)
	// Rule[perm:read]
	ChainGetBlockMessages(ctx context.Context, bid cid.Cid) (*apitypes.BlockMessages, error)
	// Rule[perm:read]
	ChainGetMessagesInTipset(ctx context.Context, key types.TipSetKey) ([]apitypes.Message, error)
	// Rule[perm:read]
	ChainGetReceipts(ctx context.Context, id cid.Cid) ([]types.MessageReceipt, error)
	// Rule[perm:read]
	ChainGetParentMessages(ctx context.Context, bcid cid.Cid) ([]apitypes.Message, error)
	// Rule[perm:read]
	ChainGetParentReceipts(ctx context.Context, bcid cid.Cid) ([]*types.MessageReceipt, error)
	//Rule[perm:read]
	StateVerifiedRegistryRootKey(ctx context.Context, tsk types.TipSetKey) (address.Address, error)
	// Rule[perm:read]
	StateVerifierStatus(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*abi.StoragePower, error)
	// Rule[perm:read]
	ChainNotify(ctx context.Context) <-chan []*apitypes.HeadChange
	// Rule[perm:read]
	GetFullBlock(ctx context.Context, id cid.Cid) (*types.FullBlock, error)
	// Rule[perm:read]
	GetActor(ctx context.Context, addr address.Address) (*types.Actor, error)
	// Rule[perm:read]
	GetParentStateRootActor(ctx context.Context, ts *types.TipSet, addr address.Address) (*types.Actor, error)
	// Rule[perm:read]
	GetEntry(ctx context.Context, height abi.ChainEpoch, round uint64) (*types.BeaconEntry, error)
	// Rule[perm:read]
	MessageWait(ctx context.Context, msgCid cid.Cid, confidence, lookback abi.ChainEpoch) (*apitypes.ChainMessage, error)
	// Rule[perm:read]
	ProtocolParameters(ctx context.Context) (*apitypes.ProtocolParams, error)
	// Rule[perm:read]
	ResolveToKeyAddr(ctx context.Context, addr address.Address, ts *types.TipSet) (address.Address, error)
	// Rule[perm:read]
	StateNetworkName(ctx context.Context) (apitypes.NetworkName, error)
	// Rule[perm:read]
	StateGetReceipt(ctx context.Context, msg cid.Cid, from types.TipSetKey) (*types.MessageReceipt, error)
	// Rule[perm:read]
	StateSearchMsg(ctx context.Context, msg cid.Cid) (*apitypes.MsgLookup, error)
	// Rule[perm:read]
	StateSearchMsgLimited(ctx context.Context, cid cid.Cid, limit abi.ChainEpoch) (*apitypes.MsgLookup, error)
	// Rule[perm:read]
	StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64) (*apitypes.MsgLookup, error)
	// Rule[perm:read]
	StateWaitMsgLimited(ctx context.Context, cid cid.Cid, confidence uint64, limit abi.ChainEpoch) (*apitypes.MsgLookup, error)
	// Rule[perm:read]
	StateNetworkVersion(ctx context.Context, tsk types.TipSetKey) (network.Version, error)
	// Rule[perm:read]
	VerifyEntry(parent, child *types.BeaconEntry, height abi.ChainEpoch) bool
	// Rule[perm:read]
	ChainExport(context.Context, abi.ChainEpoch, bool, types.TipSetKey) (<-chan []byte, error)
	// Rule[perm:read]
	ChainGetPath(ctx context.Context, from types.TipSetKey, to types.TipSetKey) ([]*apitypes.HeadChange, error)
}

type IMinerState interface {
	// Rule[perm:read]
	StateMinerSectorAllocated(ctx context.Context, maddr address.Address, s abi.SectorNumber, tsk types.TipSetKey) (bool, error)
	// Rule[perm:read]
	StateSectorPreCommitInfo(ctx context.Context, maddr address.Address, n abi.SectorNumber, tsk types.TipSetKey) (miner.SectorPreCommitOnChainInfo, error)
	// Rule[perm:read]
	StateSectorGetInfo(ctx context.Context, maddr address.Address, n abi.SectorNumber, tsk types.TipSetKey) (*miner.SectorOnChainInfo, error)
	// Rule[perm:read]
	StateSectorPartition(ctx context.Context, maddr address.Address, sectorNumber abi.SectorNumber, tsk types.TipSetKey) (*miner.SectorLocation, error)
	// Rule[perm:read]
	StateMinerSectorSize(ctx context.Context, maddr address.Address, tsk types.TipSetKey) (abi.SectorSize, error)
	// Rule[perm:read]
	StateMinerInfo(ctx context.Context, maddr address.Address, tsk types.TipSetKey) (miner.MinerInfo, error)
	// Rule[perm:read]
	StateMinerWorkerAddress(ctx context.Context, maddr address.Address, tsk types.TipSetKey) (address.Address, error)
	// Rule[perm:read]
	StateMinerRecoveries(ctx context.Context, maddr address.Address, tsk types.TipSetKey) (bitfield.BitField, error)
	// Rule[perm:read]
	StateMinerFaults(ctx context.Context, maddr address.Address, tsk types.TipSetKey) (bitfield.BitField, error)
	// Rule[perm:read]
	StateMinerProvingDeadline(ctx context.Context, maddr address.Address, tsk types.TipSetKey) (*dline.Info, error)
	// Rule[perm:read]
	StateMinerPartitions(ctx context.Context, maddr address.Address, dlIdx uint64, tsk types.TipSetKey) ([]apitypes.Partition, error)
	// Rule[perm:read]
	StateMinerDeadlines(ctx context.Context, maddr address.Address, tsk types.TipSetKey) ([]apitypes.Deadline, error)
	// Rule[perm:read]
	StateMinerSectors(ctx context.Context, maddr address.Address, sectorNos *bitfield.BitField, tsk types.TipSetKey) ([]*miner.SectorOnChainInfo, error)
	// Rule[perm:read]
	StateMarketStorageDeal(ctx context.Context, dealID abi.DealID, tsk types.TipSetKey) (*apitypes.MarketDeal, error)
	// Rule[perm:read]
	StateMinerPreCommitDepositForPower(ctx context.Context, maddr address.Address, pci miner.SectorPreCommitInfo, tsk types.TipSetKey) (big.Int, error)
	// Rule[perm:read]
	StateMinerInitialPledgeCollateral(ctx context.Context, maddr address.Address, pci miner.SectorPreCommitInfo, tsk types.TipSetKey) (big.Int, error)
	// Rule[perm:read]
	StateVMCirculatingSupplyInternal(ctx context.Context, tsk types.TipSetKey) (types.CirculatingSupply, error)
	// Rule[perm:read]
	StateCirculatingSupply(ctx context.Context, tsk types.TipSetKey) (abi.TokenAmount, error)
	// Rule[perm:read]
	StateMarketDeals(ctx context.Context, tsk types.TipSetKey) (map[string]apitypes.MarketDeal, error)
	// Rule[perm:read]
	StateMinerActiveSectors(ctx context.Context, maddr address.Address, tsk types.TipSetKey) ([]*miner.SectorOnChainInfo, error)
	// Rule[perm:read]
	StateLookupID(ctx context.Context, addr address.Address, tsk types.TipSetKey) (address.Address, error)
	// Rule[perm:read]
	StateListMiners(ctx context.Context, tsk types.TipSetKey) ([]address.Address, error)
	// Rule[perm:read]
	StateListActors(ctx context.Context, tsk types.TipSetKey) ([]address.Address, error)
	// Rule[perm:read]
	StateMinerPower(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*apitypes.MinerPower, error)
	// Rule[perm:read]
	StateMinerAvailableBalance(ctx context.Context, maddr address.Address, tsk types.TipSetKey) (big.Int, error)
	// Rule[perm:read]
	StateSectorExpiration(ctx context.Context, maddr address.Address, sectorNumber abi.SectorNumber, tsk types.TipSetKey) (*miner.SectorExpiration, error)
	// Rule[perm:read]
	StateMinerSectorCount(ctx context.Context, addr address.Address, tsk types.TipSetKey) (apitypes.MinerSectors, error)
	// Rule[perm:read]
	StateMarketBalance(ctx context.Context, addr address.Address, tsk types.TipSetKey) (apitypes.MarketBalance, error)
	// Rule[perm:read]
	StateDealProviderCollateralBounds(ctx context.Context, size abi.PaddedPieceSize, verified bool, tsk types.TipSetKey) (apitypes.DealCollateralBounds, error)
	// Rule[perm:read]
	StateVerifiedClientStatus(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*abi.StoragePower, error)
}

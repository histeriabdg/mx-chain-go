package processor

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-go/state"
)

// InterceptedTransactionHandler defines an intercepted data wrapper over transaction handler that has
// receiver and sender shard getters
type InterceptedTransactionHandler interface {
	SenderShardId() uint32
	ReceiverShardId() uint32
	Nonce() uint64
	SenderAddress() []byte
	Fee() *big.Int
	Transaction() data.TransactionHandler
}

// ShardedPool is a perspective of the sharded data pool
type ShardedPool interface {
	AddData(key []byte, data interface{}, sizeInBytes int, cacheID string)
}

type interceptedDataSizeHandler interface {
	SizeInBytes() int
}

type interceptedHeartbeatMessageHandler interface {
	interceptedDataSizeHandler
	Message() interface{}
}

type interceptedPeerAuthenticationMessageHandler interface {
	interceptedDataSizeHandler
	Message() interface{}
	Payload() []byte
	Pubkey() []byte
}

type interceptedValidatorInfo interface {
	Hash() []byte
	ValidatorInfo() *state.ShardValidatorInfo
}

package preprocess

import (
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-go/integrationTests"
	"github.com/multiversx/mx-chain-go/process"
	"github.com/multiversx/mx-chain-go/storage/txcache"
	"time"
)

type sovereignChainTransactions struct {
	*transactions
}

// NewSovereignChainTransactionPreprocessor creates a new transaction preprocessor object
func NewSovereignChainTransactionPreprocessor(
	transactions *transactions,
) (*sovereignChainTransactions, error) {
	if transactions == nil {
		return nil, process.ErrNilPreProcessor
	}

	sctp := &sovereignChainTransactions{
		transactions,
	}

	sctp.scheduledTXContinueFunc = sctp.shouldContinueProcessingScheduledTx

	return sctp, nil
}

// ProcessBlockTransactions processes all the transaction from the block.Body, updates the state
func (sctp *sovereignChainTransactions) ProcessBlockTransactions(
	header data.HeaderHandler,
	body *block.Body,
	haveTime func() bool,
) (block.MiniBlockSlice, error) {
	calculatedMiniBlocks, _, err := sctp.processTxsFromMe(body, haveTime, header.GetPrevRandSeed())
	return calculatedMiniBlocks, err
}

// CreateAndProcessMiniBlocks creates miniBlocks from selected transactions
func (sctp *sovereignChainTransactions) CreateAndProcessMiniBlocks(haveTime func() bool, randomness []byte) (block.MiniBlockSlice, error) {
	//TODO: Replace this value with the real one
	gasBandwidth := integrationTests.MaxGasLimitPerBlock

	startTime := time.Now()

	sortedTxs, remainingTxsForScheduled, err := sctp.computeSortedTxs(sctp.shardCoordinator.SelfId(), sctp.shardCoordinator.SelfId(), gasBandwidth, randomness)
	elapsedTime := time.Since(startTime)
	if err != nil {
		log.Debug("computeSortedTxs", "error", err.Error())
		return make(block.MiniBlockSlice, 0), nil
	}

	if len(sortedTxs) == 0 {
		log.Trace("no transaction found after computeSortedTxs",
			"time [s]", elapsedTime,
		)
		return make(block.MiniBlockSlice, 0), nil
	}

	if !haveTime() {
		log.Debug("time is up after computeSortedTxs",
			"num txs", len(sortedTxs),
			"time [s]", elapsedTime,
		)
		return make(block.MiniBlockSlice, 0), nil
	}

	log.Debug("elapsed time to computeSortedTxs",
		"num txs", len(sortedTxs),
		"time [s]", elapsedTime,
	)

	sortedTxsForScheduled := append(sortedTxs, remainingTxsForScheduled...)
	sortedTxsForScheduled, _ = sctp.prefilterTransactions(nil, sortedTxsForScheduled, 0, gasBandwidth)
	sctp.sortTransactionsBySenderAndNonce(sortedTxsForScheduled, randomness)

	scheduledMiniBlocks, err := sctp.createScheduledMiniBlocksFromMeAsProposer(
		haveTime,
		haveAdditionalTimeFalse,
		sortedTxsForScheduled,
		make(map[string]struct{}),
	)
	if err != nil {
		log.Debug("createScheduledMiniBlocksFromMeAsProposer", "error", err.Error())
		return make(block.MiniBlockSlice, 0), nil
	}

	return scheduledMiniBlocks, nil
}

// ProcessMiniBlock does nothing on sovereign chain
func (sctp *sovereignChainTransactions) ProcessMiniBlock(
	_ *block.MiniBlock,
	_ func() bool,
	_ func() bool,
	_ bool,
	_ bool,
	_ int,
	_ process.PreProcessorExecutionInfoHandler,
) ([][]byte, int, bool, error) {
	return nil, 0, false, nil
}

func (sctp *sovereignChainTransactions) shouldContinueProcessingScheduledTx(
	_ func(uint32) bool,
	wrappedTx *txcache.WrappedTransaction,
	_ map[string]struct{},
	mbInfo *createScheduledMiniBlocksInfo,
) (*transaction.Transaction, *block.MiniBlock, bool) {
	txHash := wrappedTx.TxHash
	senderShardID := wrappedTx.SenderShardID
	receiverShardID := wrappedTx.ReceiverShardID

	tx, ok := wrappedTx.Tx.(*transaction.Transaction)
	if !ok {
		log.Debug("wrong type assertion",
			"hash", txHash,
			"sender shard", senderShardID,
			"receiver shard", receiverShardID)
		return nil, nil, false
	}

	miniBlock, ok := mbInfo.mapMiniBlocks[receiverShardID]
	if !ok {
		log.Debug("scheduled mini block is not created", "shard", receiverShardID)
		return nil, nil, false
	}

	addressHasEnoughBalance := sctp.hasAddressEnoughInitialBalance(tx)
	if !addressHasEnoughBalance {
		mbInfo.schedulingInfo.numScheduledTxsWithInitialBalanceConsumed++
		return nil, nil, false
	}

	return tx, miniBlock, true
}

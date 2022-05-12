package staking

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon/stakingcommon"
	"github.com/ElrondNetwork/elrond-go/vm"
	"github.com/ElrondNetwork/elrond-go/vm/systemSmartContracts"
	"github.com/stretchr/testify/require"
)

func requireSliceContains(t *testing.T, s1, s2 [][]byte) {
	for _, elemInS2 := range s2 {
		require.Contains(t, s1, elemInS2)
	}
}

func requireSliceContainsNumOfElements(t *testing.T, s1, s2 [][]byte, numOfElements int) {
	foundCt := 0
	for _, elemInS2 := range s2 {
		if searchInSlice(s1, elemInS2) {
			foundCt++
		}
	}

	require.Equal(t, numOfElements, foundCt)
}

func requireSameSliceDifferentOrder(t *testing.T, s1, s2 [][]byte) {
	require.Equal(t, len(s1), len(s2))

	for _, elemInS1 := range s1 {
		require.Contains(t, s2, elemInS1)
	}
}

func searchInSlice(s1 [][]byte, s2 []byte) bool {
	for _, elemInS1 := range s1 {
		if bytes.Equal(elemInS1, s2) {
			return true
		}
	}

	return false
}

func searchInMap(validatorMap map[uint32][][]byte, pk []byte) bool {
	for _, validatorsInShard := range validatorMap {
		for _, val := range validatorsInShard {
			if bytes.Equal(val, pk) {
				return true
			}
		}
	}
	return false
}

func requireMapContains(t *testing.T, m map[uint32][][]byte, s [][]byte) {
	for _, elemInSlice := range s {
		require.True(t, searchInMap(m, elemInSlice))
	}
}

func requireMapDoesNotContain(t *testing.T, m map[uint32][][]byte, s [][]byte) {
	for _, elemInSlice := range s {
		require.False(t, searchInMap(m, elemInSlice))
	}
}

func remove(s [][]byte, elem []byte) [][]byte {
	ret := s
	for i, e := range s {
		if bytes.Equal(elem, e) {
			ret[i] = ret[len(s)-1]
			return ret[:len(s)-1]
		}
	}

	return ret
}

func unStake(t *testing.T, owner []byte, accountsDB state.AccountsAdapter, marshaller marshal.Marshalizer, stake *big.Int) {
	validatorSC := stakingcommon.LoadUserAccount(accountsDB, vm.ValidatorSCAddress)
	ownerStoredData, err := validatorSC.DataTrieTracker().RetrieveValue(owner)
	require.Nil(t, err)

	validatorData := &systemSmartContracts.ValidatorDataV2{}
	err = marshaller.Unmarshal(validatorData, ownerStoredData)
	require.Nil(t, err)

	validatorData.TotalStakeValue.Sub(validatorData.TotalStakeValue, stake)
	marshaledData, _ := marshaller.Marshal(validatorData)
	err = validatorSC.DataTrieTracker().SaveKeyValue(owner, marshaledData)
	require.Nil(t, err)

	err = accountsDB.SaveAccount(validatorSC)
	require.Nil(t, err)
	_, err = accountsDB.Commit()
	require.Nil(t, err)
}

func TestStakingV4(t *testing.T) {
	numOfMetaNodes := uint32(400)
	numOfShards := uint32(3)
	numOfEligibleNodesPerShard := uint32(400)
	numOfWaitingNodesPerShard := uint32(400)
	numOfNodesToShufflePerShard := uint32(80)
	shardConsensusGroupSize := 266
	metaConsensusGroupSize := 266
	numOfNodesInStakingQueue := uint32(60)

	totalEligible := int(numOfEligibleNodesPerShard*numOfShards) + int(numOfMetaNodes) // 1600
	totalWaiting := int(numOfWaitingNodesPerShard*numOfShards) + int(numOfMetaNodes)   // 1600

	node := NewTestMetaProcessor(
		numOfMetaNodes,
		numOfShards,
		numOfEligibleNodesPerShard,
		numOfWaitingNodesPerShard,
		numOfNodesToShufflePerShard,
		shardConsensusGroupSize,
		metaConsensusGroupSize,
		numOfNodesInStakingQueue,
	)
	node.EpochStartTrigger.SetRoundsPerEpoch(4)

	// 1. Check initial config is correct
	initialNodes := node.NodesConfig
	require.Len(t, getAllPubKeys(initialNodes.eligible), totalEligible)
	require.Len(t, getAllPubKeys(initialNodes.waiting), totalWaiting)
	require.Len(t, initialNodes.queue, int(numOfNodesInStakingQueue))
	require.Empty(t, initialNodes.shuffledOut)
	require.Empty(t, initialNodes.auction)

	// 2. Check config after staking v4 initialization
	node.Process(t, 5)
	nodesConfigStakingV4Init := node.NodesConfig
	require.Len(t, getAllPubKeys(nodesConfigStakingV4Init.eligible), totalEligible)
	require.Len(t, getAllPubKeys(nodesConfigStakingV4Init.waiting), totalWaiting)
	require.Empty(t, nodesConfigStakingV4Init.queue)
	require.Empty(t, nodesConfigStakingV4Init.shuffledOut)
	requireSameSliceDifferentOrder(t, initialNodes.queue, nodesConfigStakingV4Init.auction)

	// 3. Check config after first staking v4 epoch, WITHOUT distribution from auction -> waiting
	node.Process(t, 6)
	nodesConfigStakingV4 := node.NodesConfig
	require.Len(t, getAllPubKeys(nodesConfigStakingV4.eligible), totalEligible) // 1600

	numOfShuffledOut := int((numOfShards + 1) * numOfNodesToShufflePerShard) // 320
	require.Len(t, getAllPubKeys(nodesConfigStakingV4.shuffledOut), numOfShuffledOut)

	newWaiting := totalWaiting - numOfShuffledOut // 1280 (1600 - 320)
	require.Len(t, getAllPubKeys(nodesConfigStakingV4.waiting), newWaiting)

	// 380 (320 from shuffled out + 60 from initial staking queue -> auction from stakingV4 init)
	auctionListSize := numOfShuffledOut + len(nodesConfigStakingV4Init.auction)
	require.Len(t, nodesConfigStakingV4.auction, auctionListSize)
	requireSliceContains(t, nodesConfigStakingV4.auction, nodesConfigStakingV4Init.auction)

	require.Empty(t, nodesConfigStakingV4.queue)
	require.Empty(t, nodesConfigStakingV4.leaving)

	// 320 nodes which are now in eligible are from previous waiting list
	requireSliceContainsNumOfElements(t, getAllPubKeys(nodesConfigStakingV4.eligible), getAllPubKeys(nodesConfigStakingV4Init.waiting), numOfShuffledOut)

	// All shuffled out are from previous staking v4 init eligible
	requireMapContains(t, nodesConfigStakingV4Init.eligible, getAllPubKeys(nodesConfigStakingV4.shuffledOut))

	// All shuffled out are in auction
	requireSliceContains(t, nodesConfigStakingV4.auction, getAllPubKeys(nodesConfigStakingV4.shuffledOut))

	// No auction node from previous epoch has been moved to waiting
	requireMapDoesNotContain(t, nodesConfigStakingV4.waiting, nodesConfigStakingV4Init.auction)

	epochs := 0
	prevConfig := nodesConfigStakingV4
	numOfSelectedNodesFromAuction := numOfShuffledOut                     // 320, since we will always fill shuffled out nodes with this config
	numOfUnselectedNodesFromAuction := auctionListSize - numOfShuffledOut // 60 = 380 - 320
	for epochs < 10 {
		node.Process(t, 5)
		newNodeConfig := node.NodesConfig

		require.Len(t, getAllPubKeys(newNodeConfig.eligible), totalEligible)       // 1600
		require.Len(t, getAllPubKeys(newNodeConfig.waiting), newWaiting)           // 1280
		require.Len(t, getAllPubKeys(newNodeConfig.shuffledOut), numOfShuffledOut) // 320
		require.Len(t, newNodeConfig.auction, auctionListSize)                     // 380
		require.Empty(t, newNodeConfig.queue)
		require.Empty(t, newNodeConfig.leaving)

		// 320 nodes which are now in eligible are from previous waiting list
		requireSliceContainsNumOfElements(t, getAllPubKeys(newNodeConfig.eligible), getAllPubKeys(prevConfig.waiting), numOfShuffledOut)

		// New auction list also contains unselected nodes from previous auction list
		requireSliceContainsNumOfElements(t, newNodeConfig.auction, prevConfig.auction, numOfUnselectedNodesFromAuction)

		// All shuffled out are from previous eligible config
		requireMapContains(t, prevConfig.eligible, getAllPubKeys(newNodeConfig.shuffledOut))

		// All shuffled out are now in auction
		requireSliceContains(t, newNodeConfig.auction, getAllPubKeys(newNodeConfig.shuffledOut))

		// 320 nodes which have been selected from previous auction list are now in waiting
		requireSliceContainsNumOfElements(t, getAllPubKeys(newNodeConfig.waiting), prevConfig.auction, numOfSelectedNodesFromAuction)

		prevConfig = newNodeConfig
		epochs++
	}
}

func TestStakingV4MetaProcessor_ProcessMultipleNodesWithSameSetupExpectSameRootHash(t *testing.T) {
	numOfMetaNodes := uint32(6)
	numOfShards := uint32(3)
	numOfEligibleNodesPerShard := uint32(6)
	numOfWaitingNodesPerShard := uint32(6)
	numOfNodesToShufflePerShard := uint32(2)
	shardConsensusGroupSize := 2
	metaConsensusGroupSize := 2
	numOfNodesInStakingQueue := uint32(2)

	nodes := make([]*TestMetaProcessor, 0, numOfMetaNodes)
	for i := uint32(0); i < numOfMetaNodes; i++ {
		nodes = append(nodes, NewTestMetaProcessor(
			numOfMetaNodes,
			numOfShards,
			numOfEligibleNodesPerShard,
			numOfWaitingNodesPerShard,
			numOfNodesToShufflePerShard,
			shardConsensusGroupSize,
			metaConsensusGroupSize,
			numOfNodesInStakingQueue,
		))
		nodes[i].EpochStartTrigger.SetRoundsPerEpoch(4)
	}

	numOfEpochs := uint32(15)
	rootHashes := make(map[uint32][][]byte)
	for currEpoch := uint32(1); currEpoch <= numOfEpochs; currEpoch++ {
		for _, node := range nodes {
			rootHash, _ := node.ValidatorStatistics.RootHash()
			rootHashes[currEpoch] = append(rootHashes[currEpoch], rootHash)

			node.Process(t, 5)
			require.Equal(t, currEpoch, node.EpochStartTrigger.Epoch())
		}
	}

	for _, rootHashesInEpoch := range rootHashes {
		firstNodeRootHashInEpoch := rootHashesInEpoch[0]
		for _, rootHash := range rootHashesInEpoch {
			require.Equal(t, firstNodeRootHashInEpoch, rootHash)
		}
	}
}

func TestStakingV4_UnStakeNodesWithNotEnoughFunds(t *testing.T) {
	pubKeys := generateAddresses(0, 20)

	// Owner1 has 8 nodes, but enough stake for just 7 nodes. At the end of the epoch(staking v4 init),
	// the last node from staking queue should be unStaked
	owner1 := "owner1"
	owner1Stats := &OwnerStats{
		EligibleBlsKeys: map[uint32][][]byte{
			core.MetachainShardId: pubKeys[:3],
		},
		WaitingBlsKeys: map[uint32][][]byte{
			0: pubKeys[3:6],
		},
		StakingQueueKeys: pubKeys[6:8],
		TotalStake:       big.NewInt(7 * nodePrice),
	}

	// Owner2 has 6 nodes, but enough stake for just 5 nodes. At the end of the epoch(staking v4 init),
	// one node from waiting list should be unStaked
	owner2 := "owner2"
	owner2Stats := &OwnerStats{
		EligibleBlsKeys: map[uint32][][]byte{
			0: pubKeys[8:11],
		},
		WaitingBlsKeys: map[uint32][][]byte{
			core.MetachainShardId: pubKeys[11:14],
		},
		TotalStake: big.NewInt(5 * nodePrice),
	}

	// Owner3 has 2 nodes in staking queue with with topUp = nodePrice
	owner3 := "owner3"
	owner3Stats := &OwnerStats{
		StakingQueueKeys: pubKeys[14:16],
		TotalStake:       big.NewInt(3 * nodePrice),
	}

	// Owner4 has 1 node in staking queue with topUp = nodePrice
	owner4 := "owner4"
	owner4Stats := &OwnerStats{
		StakingQueueKeys: pubKeys[16:17],
		TotalStake:       big.NewInt(2 * nodePrice),
	}

	cfg := &InitialNodesConfig{
		MetaConsensusGroupSize:        2,
		ShardConsensusGroupSize:       2,
		MinNumberOfEligibleShardNodes: 3,
		MinNumberOfEligibleMetaNodes:  3,
		NumOfShards:                   1,
		Owners: map[string]*OwnerStats{
			owner1: owner1Stats,
			owner2: owner2Stats,
			owner3: owner3Stats,
			owner4: owner4Stats,
		},
		MaxNodesChangeConfig: []config.MaxNodesChangeConfig{
			{
				EpochEnable:            0,
				MaxNumNodes:            12,
				NodesToShufflePerShard: 1,
			},
			{
				EpochEnable:            stakingV4DistributeAuctionToWaitingEpoch,
				MaxNumNodes:            10,
				NodesToShufflePerShard: 1,
			},
		},
	}
	node := NewTestMetaProcessorWithCustomNodes(cfg)
	node.EpochStartTrigger.SetRoundsPerEpoch(4)

	// 1. Check initial config is correct
	currNodesConfig := node.NodesConfig
	require.Len(t, getAllPubKeys(currNodesConfig.eligible), 6)
	require.Len(t, getAllPubKeys(currNodesConfig.waiting), 6)
	require.Len(t, currNodesConfig.eligible[core.MetachainShardId], 3)
	require.Len(t, currNodesConfig.waiting[core.MetachainShardId], 3)
	require.Len(t, currNodesConfig.eligible[0], 3)
	require.Len(t, currNodesConfig.waiting[0], 3)

	requireSliceContainsNumOfElements(t, currNodesConfig.eligible[core.MetachainShardId], owner1Stats.EligibleBlsKeys[core.MetachainShardId], 3)
	requireSliceContainsNumOfElements(t, currNodesConfig.waiting[core.MetachainShardId], owner2Stats.WaitingBlsKeys[core.MetachainShardId], 3)
	requireSliceContainsNumOfElements(t, currNodesConfig.eligible[0], owner2Stats.EligibleBlsKeys[0], 3)
	requireSliceContainsNumOfElements(t, currNodesConfig.waiting[0], owner1Stats.WaitingBlsKeys[0], 3)

	owner1StakingQueue := owner1Stats.StakingQueueKeys
	owner3StakingQueue := owner3Stats.StakingQueueKeys
	owner4StakingQueue := owner4Stats.StakingQueueKeys
	queue := make([][]byte, 0)
	queue = append(queue, owner1StakingQueue...)
	queue = append(queue, owner3StakingQueue...)
	queue = append(queue, owner4StakingQueue...)
	require.Len(t, currNodesConfig.queue, 5)
	requireSameSliceDifferentOrder(t, currNodesConfig.queue, queue)

	require.Empty(t, currNodesConfig.shuffledOut)
	require.Empty(t, currNodesConfig.auction)

	// 2. Check config after staking v4 initialization
	node.Process(t, 5)
	currNodesConfig = node.NodesConfig
	require.Len(t, getAllPubKeys(currNodesConfig.eligible), 6)
	require.Len(t, getAllPubKeys(currNodesConfig.waiting), 5)
	require.Len(t, currNodesConfig.eligible[core.MetachainShardId], 3)
	require.Len(t, currNodesConfig.waiting[core.MetachainShardId], 2)
	require.Len(t, currNodesConfig.eligible[0], 3)
	require.Len(t, currNodesConfig.waiting[0], 3)

	// Owner1 will have the second node from queue removed, before adding all the nodes to auction list
	queue = remove(queue, owner1StakingQueue[1])
	require.Empty(t, currNodesConfig.queue)
	require.Len(t, currNodesConfig.auction, 4)
	requireSameSliceDifferentOrder(t, currNodesConfig.auction, queue)

	// Owner2 will have one of the nodes in waiting list removed
	require.Len(t, getAllPubKeys(currNodesConfig.leaving), 1)
	requireSliceContainsNumOfElements(t, getAllPubKeys(currNodesConfig.leaving), getAllPubKeys(owner2Stats.WaitingBlsKeys), 1)

	// Owner1 will unStake some EGLD => at the end of next epoch, he should have the other node from queue(now auction list) removed
	unStake(t, []byte(owner1), node.AccountsAdapter, node.Marshaller, big.NewInt(0.1*nodePrice))

	// 3. Check config in epoch = staking v4
	node.Process(t, 5)
	currNodesConfig = node.NodesConfig
	require.Len(t, getAllPubKeys(currNodesConfig.eligible), 6)
	require.Len(t, getAllPubKeys(currNodesConfig.waiting), 3)
	require.Len(t, getAllPubKeys(currNodesConfig.shuffledOut), 2)

	require.Len(t, currNodesConfig.eligible[core.MetachainShardId], 3)
	require.Len(t, currNodesConfig.waiting[core.MetachainShardId], 1)
	require.Len(t, currNodesConfig.shuffledOut[core.MetachainShardId], 1)
	require.Len(t, currNodesConfig.eligible[0], 3)
	require.Len(t, currNodesConfig.waiting[0], 2)
	require.Len(t, currNodesConfig.shuffledOut[0], 1)

	// Owner1 will have the last node from auction list removed
	queue = remove(queue, owner1StakingQueue[0])
	require.Len(t, currNodesConfig.auction, 3)
	requireSameSliceDifferentOrder(t, currNodesConfig.auction, queue)
	require.Len(t, getAllPubKeys(currNodesConfig.leaving), 1)
	require.Equal(t, getAllPubKeys(currNodesConfig.leaving)[0], owner1StakingQueue[0])

	// Owner3 will unStake EGLD => he will have negative top-up at the selection time => one of his nodes will be unStaked.
	// His other node should not have been selected => remains in auction.
	// Meanwhile, owner4 had never unStaked EGLD => his node from auction list node will be distributed to waiting
	unStake(t, []byte(owner3), node.AccountsAdapter, node.Marshaller, big.NewInt(2*nodePrice))

	// 4. Check config in epoch = staking v4 distribute auction to waiting
	node.Process(t, 5)
	currNodesConfig = node.NodesConfig
	requireSliceContainsNumOfElements(t, getAllPubKeys(currNodesConfig.leaving), owner3StakingQueue, 1)
	requireSliceContainsNumOfElements(t, currNodesConfig.auction, owner3StakingQueue, 1)
	requireSliceContainsNumOfElements(t, getAllPubKeys(currNodesConfig.waiting), owner4StakingQueue, 1)
}

func TestStakingV4_StakeNewNodes(t *testing.T) {
	pubKeys := generateAddresses(0, 20)

	// Owner1 has 6 nodes, zero top up
	owner1 := "owner1"
	owner1Stats := &OwnerStats{
		EligibleBlsKeys: map[uint32][][]byte{
			core.MetachainShardId: pubKeys[:2],
		},
		WaitingBlsKeys: map[uint32][][]byte{
			0: pubKeys[2:4],
		},
		StakingQueueKeys: pubKeys[4:6],
		TotalStake:       big.NewInt(6 * nodePrice),
	}

	// Owner2 has 4 nodes, zero top up
	owner2 := "owner2"
	owner2Stats := &OwnerStats{
		EligibleBlsKeys: map[uint32][][]byte{
			0: pubKeys[6:8],
		},
		WaitingBlsKeys: map[uint32][][]byte{
			core.MetachainShardId: pubKeys[8:10],
		},
		TotalStake: big.NewInt(4 * nodePrice),
	}
	// Owner3 has 1 node in staking queue with topUp = nodePrice
	owner3 := "owner3"
	owner3Stats := &OwnerStats{
		StakingQueueKeys: pubKeys[10:11],
		TotalStake:       big.NewInt(2 * nodePrice),
	}

	cfg := &InitialNodesConfig{
		MetaConsensusGroupSize:        1,
		ShardConsensusGroupSize:       1,
		MinNumberOfEligibleShardNodes: 1,
		MinNumberOfEligibleMetaNodes:  1,
		NumOfShards:                   1,
		Owners: map[string]*OwnerStats{
			owner1: owner1Stats,
			owner2: owner2Stats,
			owner3: owner3Stats,
		},
		MaxNodesChangeConfig: []config.MaxNodesChangeConfig{
			{
				EpochEnable:            0,
				MaxNumNodes:            8,
				NodesToShufflePerShard: 1,
			},
		},
	}
	node := NewTestMetaProcessorWithCustomNodes(cfg)
	node.EpochStartTrigger.SetRoundsPerEpoch(4)

	// 1. Check initial config is correct
	currNodesConfig := node.NodesConfig
	require.Len(t, getAllPubKeys(currNodesConfig.eligible), 4)
	require.Len(t, getAllPubKeys(currNodesConfig.waiting), 4)
	require.Len(t, currNodesConfig.eligible[core.MetachainShardId], 2)
	require.Len(t, currNodesConfig.waiting[core.MetachainShardId], 2)
	require.Len(t, currNodesConfig.eligible[0], 2)
	require.Len(t, currNodesConfig.waiting[0], 2)

	owner1StakingQueue := owner1Stats.StakingQueueKeys
	owner3StakingQueue := owner3Stats.StakingQueueKeys
	queue := make([][]byte, 0)
	queue = append(queue, owner1StakingQueue...)
	queue = append(queue, owner3StakingQueue...)
	require.Len(t, currNodesConfig.queue, 3)
	requireSameSliceDifferentOrder(t, currNodesConfig.queue, queue)

	require.Empty(t, currNodesConfig.shuffledOut)
	require.Empty(t, currNodesConfig.auction)

	// NewOwner1 stakes 1 node with top up = 2*node price; should be sent to auction list
	newOwner1 := "newOwner1"
	newNodes1 := map[string]*NodesRegisterData{
		newOwner1: {
			BLSKeys:    [][]byte{generateAddress(444)},
			TotalStake: big.NewInt(3 * nodePrice),
		},
	}
	// 2. Check config after staking v4 init when a new node is staked
	node.Process(t, 5)
	node.ProcessStake(t, newNodes1)
	currNodesConfig = node.NodesConfig
	queue = append(queue, newNodes1[newOwner1].BLSKeys...)
	require.Empty(t, currNodesConfig.queue)
	require.Empty(t, currNodesConfig.leaving)
	require.Len(t, currNodesConfig.auction, 4)
	requireSameSliceDifferentOrder(t, currNodesConfig.auction, queue)

	// NewOwner2 stakes 2 node with top up = 2*node price; should be sent to auction list
	newOwner2 := "newOwner2"
	newNodes2 := map[string]*NodesRegisterData{
		newOwner2: {
			BLSKeys:    [][]byte{generateAddress(555), generateAddress(666)},
			TotalStake: big.NewInt(4 * nodePrice),
		},
	}
	// 2. Check in epoch = staking v4 when 2 new nodes are staked
	node.Process(t, 4)
	node.ProcessStake(t, newNodes2)
	currNodesConfig = node.NodesConfig
	queue = append(queue, newNodes2[newOwner2].BLSKeys...)
	require.Empty(t, currNodesConfig.queue)
	requireSliceContainsNumOfElements(t, currNodesConfig.auction, queue, 6)

	// 3. Epoch =  staking v4 distribute auction to waiting
	// Only the new 2 owners + owner3 had enough top up to be distributed to waiting.
	// Meanwhile; owner1 which had 0 top up, still has his bls keys in auction
	node.Process(t, 5)
	currNodesConfig = node.NodesConfig
	require.Empty(t, currNodesConfig.queue)
	requireMapContains(t, currNodesConfig.waiting, newNodes1[newOwner1].BLSKeys)
	requireMapContains(t, currNodesConfig.waiting, newNodes2[newOwner2].BLSKeys)
	requireMapContains(t, currNodesConfig.waiting, owner3StakingQueue)
	requireSliceContains(t, currNodesConfig.auction, owner1StakingQueue)
}

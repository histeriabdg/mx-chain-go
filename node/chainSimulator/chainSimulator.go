package chainSimulator

import (
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/sharding"
	"github.com/multiversx/mx-chain-core-go/data/endProcess"
	"github.com/multiversx/mx-chain-go/config"
	"github.com/multiversx/mx-chain-go/node/chainSimulator/components"
	"github.com/multiversx/mx-chain-go/node/chainSimulator/configs"
	"github.com/multiversx/mx-chain-go/node/chainSimulator/dtos"
	"github.com/multiversx/mx-chain-go/node/chainSimulator/process"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("chainSimulator")

type simulator struct {
	chanStopNodeProcess    chan endProcess.ArgEndProcess
	syncedBroadcastNetwork components.SyncedBroadcastNetworkHandler
	handlers               []ChainHandler
	initialWalletKeys      *dtos.InitialWalletKeys
	nodes                  map[uint32]process.NodeHandler
	numOfShards            uint32
}

// NewChainSimulator will create a new instance of simulator
func NewChainSimulator(
	tempDir string,
	numOfShards uint32,
	pathToInitialConfig string,
	genesisTimestamp int64,
	roundDurationInMillis uint64,
	roundsPerEpoch core.OptionalUint64,
	apiInterface components.APIConfigurator,
) (*simulator, error) {
	syncedBroadcastNetwork := components.NewSyncedBroadcastNetwork()

	instance := &simulator{
		syncedBroadcastNetwork: syncedBroadcastNetwork,
		nodes:                  make(map[uint32]process.NodeHandler),
		handlers:               make([]ChainHandler, 0, numOfShards+1),
		numOfShards:            numOfShards,
		chanStopNodeProcess:    make(chan endProcess.ArgEndProcess),
	}

	err := instance.createChainHandlers(tempDir, numOfShards, pathToInitialConfig, genesisTimestamp, roundDurationInMillis, roundsPerEpoch, apiInterface)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (s *simulator) createChainHandlers(
	tempDir string,
	numOfShards uint32,
	originalConfigPath string,
	genesisTimestamp int64,
	roundDurationInMillis uint64,
	roundsPerEpoch core.OptionalUint64,
	apiInterface components.APIConfigurator,
) error {
	outputConfigs, err := configs.CreateChainSimulatorConfigs(configs.ArgsChainSimulatorConfigs{
		NumOfShards:           numOfShards,
		OriginalConfigsPath:   originalConfigPath,
		GenesisTimeStamp:      genesisTimestamp,
		RoundDurationInMillis: roundDurationInMillis,
		TempDir:               tempDir,
	})
	if err != nil {
		return err
	}

	if roundsPerEpoch.HasValue {
		outputConfigs.Configs.GeneralConfig.EpochStartConfig.RoundsPerEpoch = int64(roundsPerEpoch.Value)
	}

	for idx := range outputConfigs.ValidatorsPrivateKeys {
		node, errCreate := s.createTestNode(outputConfigs.Configs, idx, outputConfigs.GasScheduleFilename, apiInterface)
		if errCreate != nil {
			return errCreate
		}

		chainHandler, errCreate := process.NewBlocksCreator(node)
		if errCreate != nil {
			return errCreate
		}

		shardID := node.GetShardCoordinator().SelfId()
		s.nodes[shardID] = node
		s.handlers = append(s.handlers, chainHandler)
	}

	s.initialWalletKeys = outputConfigs.InitialWallets

	log.Info("running the chain simulator with the following parameters",
		"number of shards (including meta)", numOfShards+1,
		"round per epoch", outputConfigs.Configs.GeneralConfig.EpochStartConfig.RoundsPerEpoch,
		"round duration", time.Millisecond*time.Duration(roundDurationInMillis),
		"genesis timestamp", genesisTimestamp,
		"original config path", originalConfigPath,
		"temporary path", tempDir)

	return nil
}

func (s *simulator) createTestNode(
	configs *config.Configs,
	skIndex int,
	gasScheduleFilename string,
	apiInterface components.APIConfigurator,
) (process.NodeHandler, error) {
	args := components.ArgsTestOnlyProcessingNode{
		Configs:                *configs,
		ChanStopNodeProcess:    s.chanStopNodeProcess,
		SyncedBroadcastNetwork: s.syncedBroadcastNetwork,
		NumShards:              s.numOfShards,
		GasScheduleFilename:    gasScheduleFilename,
		SkIndex:                skIndex,
		APIInterface:           apiInterface,
	}

	return components.NewTestOnlyProcessingNode(args)
}

// GenerateBlocks will generate the provided number of blocks
func (s *simulator) GenerateBlocks(numOfBlocks int) error {
	for idx := 0; idx < numOfBlocks; idx++ {
		s.incrementRoundOnAllValidators()
		err := s.allNodesCreateBlocks()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *simulator) incrementRoundOnAllValidators() {
	for _, node := range s.handlers {
		node.IncrementRound()
	}
}

func (s *simulator) allNodesCreateBlocks() error {
	for _, node := range s.handlers {
		err := node.CreateNewBlock()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetNodeHandler returns the node handler from the provided shardID
func (s *simulator) GetNodeHandler(shardID uint32) process.NodeHandler {
	return s.nodes[shardID]
}

// GetRestAPIInterfaces will return a map with the rest api interfaces for every node
func (s *simulator) GetRestAPIInterfaces() map[uint32]string {
	resMap := make(map[uint32]string)
	for shardID, node := range s.nodes {
		resMap[shardID] = node.GetFacadeHandler().RestApiInterface()
	}

	return resMap
}

// GetInitialWalletKeys will return the initial wallet keys
func (s *simulator) GetInitialWalletKeys() *dtos.InitialWalletKeys {
	return s.initialWalletKeys
}

// SetState will set the provided state for a given address
func (s *simulator) SetState(address string, state map[string]string) error {
	addressConverter := s.nodes[core.MetachainShardId].GetCoreComponents().AddressPubKeyConverter()
	addressBytes, err := addressConverter.Decode(address)
	if err != nil {
		return err
	}

	shardID := sharding.ComputeShardID(addressBytes, s.numOfShards)
	testNode, ok := s.nodes[shardID]
	if !ok {
		return fmt.Errorf("cannot find a test node for the computed shard id, computed shard id: %d", shardID)
	}

	return testNode.SetState(addressBytes, state)
}

// Close will stop and close the simulator
func (s *simulator) Close() error {
	var errorStrings []string
	for _, n := range s.nodes {
		err := n.Close()
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}

	if len(errorStrings) == 0 {
		return nil
	}

	return components.AggregateErrors(errorStrings, components.ErrClose)
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *simulator) IsInterfaceNil() bool {
	return s == nil
}

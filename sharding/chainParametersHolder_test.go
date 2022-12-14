package sharding

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/testscommon/epochNotifier"
	"github.com/stretchr/testify/require"
)

func TestNewChainParametersHolder(t *testing.T) {
	t.Parallel()

	getDummyArgs := func() ArgsChainParametersHolder {
		return ArgsChainParametersHolder{
			EpochNotifier: &epochNotifier.EpochNotifierStub{},
			ChainParameters: []config.ChainParametersByEpochConfig{
				{
					EnableEpoch:                 0,
					ShardMinNumNodes:            5,
					ShardConsensusGroupSize:     5,
					MetachainMinNumNodes:        7,
					MetachainConsensusGroupSize: 7,
					RoundDuration:               4000,
					Hysteresis:                  0.2,
					Adaptivity:                  false,
				},
			},
		}
	}

	t.Run("nil epoch notifier", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		args.EpochNotifier = nil

		paramsHolder, err := NewChainParametersHolder(args)
		require.True(t, check.IfNil(paramsHolder))
		require.Equal(t, ErrNilEpochNotifier, err)
	})

	t.Run("empty chain parameters", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		args.ChainParameters = nil

		paramsHolder, err := NewChainParametersHolder(args)
		require.True(t, check.IfNil(paramsHolder))
		require.Equal(t, ErrMissingChainParameters, err)
	})

	t.Run("invalid shard consensus size", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		args.ChainParameters[0].ShardConsensusGroupSize = 0

		paramsHolder, err := NewChainParametersHolder(args)
		require.True(t, check.IfNil(paramsHolder))
		require.ErrorIs(t, err, ErrNegativeOrZeroConsensusGroupSize)
	})

	t.Run("min nodes per shard smaller than consensus size", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		args.ChainParameters[0].ShardConsensusGroupSize = 5
		args.ChainParameters[0].ShardMinNumNodes = 4

		paramsHolder, err := NewChainParametersHolder(args)
		require.True(t, check.IfNil(paramsHolder))
		require.ErrorIs(t, err, ErrMinNodesPerShardSmallerThanConsensusSize)
	})

	t.Run("invalid metachain consensus size", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		args.ChainParameters[0].MetachainConsensusGroupSize = 0

		paramsHolder, err := NewChainParametersHolder(args)
		require.True(t, check.IfNil(paramsHolder))
		require.ErrorIs(t, err, ErrNegativeOrZeroConsensusGroupSize)
	})

	t.Run("min nodes meta smaller than consensus size", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		args.ChainParameters[0].MetachainConsensusGroupSize = 5
		args.ChainParameters[0].MetachainMinNumNodes = 4

		paramsHolder, err := NewChainParametersHolder(args)
		require.True(t, check.IfNil(paramsHolder))
		require.ErrorIs(t, err, ErrMinNodesPerShardSmallerThanConsensusSize)
	})

	t.Run("invalid future chain parameters", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		newChainParameters := args.ChainParameters[0]
		newChainParameters.ShardConsensusGroupSize = 0
		args.ChainParameters = append(args.ChainParameters, newChainParameters)

		paramsHolder, err := NewChainParametersHolder(args)
		require.True(t, check.IfNil(paramsHolder))
		require.ErrorIs(t, err, ErrNegativeOrZeroConsensusGroupSize)
		require.Contains(t, err.Error(), "index 1")
	})

	t.Run("no config for epoch 0", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		args.ChainParameters[0].EnableEpoch = 37
		paramsHolder, err := NewChainParametersHolder(args)
		require.True(t, check.IfNil(paramsHolder))
		require.ErrorIs(t, err, ErrMissingConfigurationForEpochZero)
	})

	t.Run("should work and have the data ready", func(t *testing.T) {
		t.Parallel()

		args := getDummyArgs()
		secondChainParams := args.ChainParameters[0]
		secondChainParams.EnableEpoch = 5
		thirdChainParams := args.ChainParameters[0]
		thirdChainParams.EnableEpoch = 10
		args.ChainParameters = append(args.ChainParameters, secondChainParams, thirdChainParams)

		paramsHolder, err := NewChainParametersHolder(args)
		require.NoError(t, err)
		require.False(t, check.IfNil(paramsHolder))

		currentValue := paramsHolder.chainParameters[0]
		for i := 1; i < len(paramsHolder.chainParameters); i++ {
			require.Less(t, paramsHolder.chainParameters[i].EnableEpoch, currentValue.EnableEpoch)
			currentValue = paramsHolder.chainParameters[i]
		}

		require.Equal(t, uint32(0), paramsHolder.currentChainParameters.EnableEpoch)
	})
}

func TestChainParametersHolder_ChainParametersForEpoch(t *testing.T) {
	t.Parallel()

	t.Run("single configuration, should return it every time", func(t *testing.T) {
		t.Parallel()

		params := []config.ChainParametersByEpochConfig{
			{
				EnableEpoch:                 0,
				ShardConsensusGroupSize:     5,
				ShardMinNumNodes:            7,
				MetachainConsensusGroupSize: 7,
				MetachainMinNumNodes:        7,
			},
		}

		paramsHolder, _ := NewChainParametersHolder(ArgsChainParametersHolder{
			ChainParameters: params,
			EpochNotifier:   &epochNotifier.EpochNotifierStub{},
		})

		res, _ := paramsHolder.ChainParametersForEpoch(0)
		require.Equal(t, uint32(5), res.ShardConsensusGroupSize)
		require.Equal(t, uint32(7), res.MetachainConsensusGroupSize)

		res, _ = paramsHolder.ChainParametersForEpoch(1)
		require.Equal(t, uint32(5), res.ShardConsensusGroupSize)
		require.Equal(t, uint32(7), res.MetachainConsensusGroupSize)

		res, _ = paramsHolder.ChainParametersForEpoch(3700)
		require.Equal(t, uint32(5), res.ShardConsensusGroupSize)
		require.Equal(t, uint32(7), res.MetachainConsensusGroupSize)
	})

	t.Run("multiple configurations, should return the corresponding one", func(t *testing.T) {
		t.Parallel()

		params := []config.ChainParametersByEpochConfig{
			{
				EnableEpoch:                 0,
				ShardConsensusGroupSize:     5,
				ShardMinNumNodes:            7,
				MetachainConsensusGroupSize: 7,
				MetachainMinNumNodes:        7,
			},
			{
				EnableEpoch:                 10,
				ShardConsensusGroupSize:     50,
				ShardMinNumNodes:            70,
				MetachainConsensusGroupSize: 70,
				MetachainMinNumNodes:        70,
			},
			{
				EnableEpoch:                 100,
				ShardConsensusGroupSize:     500,
				ShardMinNumNodes:            700,
				MetachainConsensusGroupSize: 700,
				MetachainMinNumNodes:        700,
			},
		}

		paramsHolder, _ := NewChainParametersHolder(ArgsChainParametersHolder{
			ChainParameters: params,
			EpochNotifier:   &epochNotifier.EpochNotifierStub{},
		})

		for i := 0; i < 200; i++ {
			res, _ := paramsHolder.ChainParametersForEpoch(uint32(i))
			if i < 10 {
				require.Equal(t, uint32(5), res.ShardConsensusGroupSize)
				require.Equal(t, uint32(7), res.MetachainConsensusGroupSize)
			} else if i < 100 {
				require.Equal(t, uint32(50), res.ShardConsensusGroupSize)
				require.Equal(t, uint32(70), res.MetachainConsensusGroupSize)
			} else {
				require.Equal(t, uint32(500), res.ShardConsensusGroupSize)
				require.Equal(t, uint32(700), res.MetachainConsensusGroupSize)
			}
		}
	})
}

func TestChainParametersHolder_CurrentChainParameters(t *testing.T) {
	t.Parallel()

	params := []config.ChainParametersByEpochConfig{
		{
			EnableEpoch:                 0,
			ShardConsensusGroupSize:     5,
			ShardMinNumNodes:            7,
			MetachainConsensusGroupSize: 7,
			MetachainMinNumNodes:        7,
		},
		{
			EnableEpoch:                 10,
			ShardConsensusGroupSize:     50,
			ShardMinNumNodes:            70,
			MetachainConsensusGroupSize: 70,
			MetachainMinNumNodes:        70,
		},
	}

	paramsHolder, _ := NewChainParametersHolder(ArgsChainParametersHolder{
		ChainParameters: params,
		EpochNotifier:   &epochNotifier.EpochNotifierStub{},
	})

	paramsHolder.EpochConfirmed(0, 0)
	require.Equal(t, uint32(5), paramsHolder.CurrentChainParameters().ShardConsensusGroupSize)

	paramsHolder.EpochConfirmed(3, 0)
	require.Equal(t, uint32(5), paramsHolder.CurrentChainParameters().ShardConsensusGroupSize)

	paramsHolder.EpochConfirmed(10, 0)
	require.Equal(t, uint32(50), paramsHolder.CurrentChainParameters().ShardConsensusGroupSize)

	paramsHolder.EpochConfirmed(999, 0)
	require.Equal(t, uint32(50), paramsHolder.CurrentChainParameters().ShardConsensusGroupSize)
}

func TestChainParametersHolder_AllChainParameters(t *testing.T) {
	t.Parallel()

	params := []config.ChainParametersByEpochConfig{
		{
			EnableEpoch:                 0,
			ShardConsensusGroupSize:     5,
			ShardMinNumNodes:            7,
			MetachainConsensusGroupSize: 7,
			MetachainMinNumNodes:        7,
		},
		{
			EnableEpoch:                 10,
			ShardConsensusGroupSize:     50,
			ShardMinNumNodes:            70,
			MetachainConsensusGroupSize: 70,
			MetachainMinNumNodes:        70,
		},
	}

	paramsHolder, _ := NewChainParametersHolder(ArgsChainParametersHolder{
		ChainParameters: params,
		EpochNotifier:   &epochNotifier.EpochNotifierStub{},
	})

	returnedAllChainsParameters := paramsHolder.AllChainParameters()
	require.Equal(t, params, returnedAllChainsParameters)
	require.NotEqual(t, fmt.Sprintf("%p", returnedAllChainsParameters), fmt.Sprintf("%p", paramsHolder.chainParameters))
}

func TestChainParametersHolder_ConcurrentOperations(t *testing.T) {
	chainParams := make([]config.ChainParametersByEpochConfig, 0)
	for i := uint32(0); i <= 100; i += 5 {
		chainParams = append(chainParams, config.ChainParametersByEpochConfig{
			RoundDuration:               4000,
			Hysteresis:                  0.2,
			EnableEpoch:                 i,
			ShardConsensusGroupSize:     i*10 + 1,
			ShardMinNumNodes:            i*10 + 1,
			MetachainConsensusGroupSize: i*10 + 1,
			MetachainMinNumNodes:        i*10 + 1,
			Adaptivity:                  false,
		})
	}

	paramsHolder, _ := NewChainParametersHolder(ArgsChainParametersHolder{
		ChainParameters: chainParams,
		EpochNotifier:   &epochNotifier.EpochNotifierStub{},
	})

	numOperations := 500
	wg := sync.WaitGroup{}
	wg.Add(numOperations)
	for i := 0; i < numOperations; i++ {
		go func(idx int) {
			switch idx {
			case 0:
				paramsHolder.EpochConfirmed(uint32(idx), 0)
			case 1:
				_ = paramsHolder.CurrentChainParameters()
			case 2:
				_, _ = paramsHolder.ChainParametersForEpoch(uint32(idx))
			case 3:
				_ = paramsHolder.AllChainParameters()
			}

			wg.Done()
		}(i % 4)
	}

	wg.Wait()
}

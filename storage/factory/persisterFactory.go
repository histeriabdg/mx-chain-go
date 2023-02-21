package factory

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-go/config"
	"github.com/multiversx/mx-chain-go/storage"
	"github.com/multiversx/mx-chain-go/storage/database"
	"github.com/multiversx/mx-chain-go/storage/storageunit"
	"github.com/pelletier/go-toml"
)

const (
	dbConfigFileName         = "config.toml"
	defaultType              = "LvlDBSerial"
	defaultBatchDelaySeconds = 2
	defaultMaxBatchSize      = 100
	defaultMaxOpenFiles      = 10
)

// PersisterFactory is the factory which will handle creating new databases
type PersisterFactory struct {
	dbType              string
	batchDelaySeconds   int
	maxBatchSize        int
	maxOpenFiles        int
	shardIDProvider     storage.ShardIDProvider
	shardIDProviderType string
	numShards           uint32
}

// NewPersisterFactory will return a new instance of a PersisterFactory
func NewPersisterFactory(config config.DBConfig, shardIDProvider storage.ShardIDProvider) (*PersisterFactory, error) {
	if check.IfNil(shardIDProvider) {
		return nil, storage.ErrNilShardIDProvider
	}

	return &PersisterFactory{
		dbType:              config.Type,
		batchDelaySeconds:   config.BatchDelaySeconds,
		maxBatchSize:        config.MaxBatchSize,
		maxOpenFiles:        config.MaxOpenFiles,
		shardIDProvider:     shardIDProvider,
		shardIDProviderType: config.ShardIDProviderType,
		numShards:           config.NumShards,
	}, nil
}

// Create will return a new instance of a DB with a given path
func (pf *PersisterFactory) Create(path string) (storage.Persister, error) {
	if len(path) == 0 {
		return nil, errors.New("invalid file path")
	}

	dbConfig, err := pf.getDBConfig(path)
	if err != nil {
		return nil, err
	}

	persister, err := pf.createDB(path, dbConfig)
	if err != nil {
		return nil, err
	}

	err = pf.createPersisterConfigFile(path, dbConfig)
	if err != nil {
		return nil, err
	}

	return persister, nil
}

func (pf *PersisterFactory) getDBConfig(path string) (*config.DBConfig, error) {
	dbConfigFromFile := &config.DBConfig{}
	err := core.LoadTomlFile(dbConfigFromFile, pf.getPersisterConfigFilePath(path))
	if err == nil {
		log.Debug("getDBConfig: loaded db config from toml config file", "path", dbConfigFromFile)
		return dbConfigFromFile, nil
	}

	empty := checkIfDirIsEmpty(path)
	if !empty {
		dbConfig := &config.DBConfig{
			Type:              defaultType,
			BatchDelaySeconds: defaultBatchDelaySeconds,
			MaxBatchSize:      defaultMaxBatchSize,
			MaxOpenFiles:      defaultMaxOpenFiles,
		}

		log.Debug("getDBConfig: loaded default db config")
		return dbConfig, nil
	}

	dbConfig := &config.DBConfig{
		Type:                pf.dbType,
		BatchDelaySeconds:   pf.batchDelaySeconds,
		MaxBatchSize:        pf.maxBatchSize,
		MaxOpenFiles:        pf.maxOpenFiles,
		ShardIDProviderType: pf.shardIDProviderType,
		NumShards:           pf.numShards,
	}

	log.Debug("getDBConfig: loaded db config from main config file")
	return dbConfig, nil
}

func checkIfDirIsEmpty(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Debug("getDBConfig: failed to check if dir is empty", "path", path, "error", err.Error())
		return true
	}

	if len(files) == 0 {
		return true
	}

	return false
}

func (pf *PersisterFactory) createDB(path string, dbConfig *config.DBConfig) (storage.Persister, error) {
	dbType := storageunit.DBType(dbConfig.Type)
	switch dbType {
	case storageunit.LvlDB:
		return database.NewLevelDB(path, dbConfig.BatchDelaySeconds, dbConfig.MaxBatchSize, dbConfig.MaxOpenFiles)
	case storageunit.LvlDBSerial:
		return database.NewSerialDB(path, dbConfig.BatchDelaySeconds, dbConfig.MaxBatchSize, dbConfig.MaxOpenFiles)
	case storageunit.ShardedLvlDBSerial:
		shardIDProvider, err := pf.createShardIDProvider()
		if err != nil {
			return nil, err
		}
		return database.NewShardedDB(storageunit.LvlDBSerial, path, dbConfig.BatchDelaySeconds, dbConfig.MaxBatchSize, dbConfig.MaxOpenFiles, shardIDProvider)
	case storageunit.MemoryDB:
		return database.NewMemDB(), nil
	default:
		return nil, storage.ErrNotSupportedDBType
	}
}

func (pf *PersisterFactory) createPersisterConfigFile(path string, dbConfig *config.DBConfig) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// in memory db, no files available
			log.Debug("createPersisterConfigFile: provided path not available, config file will not be created")
			return nil
		}

		return err
	}

	configFilePath := pf.getPersisterConfigFilePath(path)
	f, err := core.OpenFile(configFilePath)
	if err == nil {
		// config file already exists, no need to create config
		return nil
	}

	defer func() {
		_ = f.Close()
	}()

	err = SaveTomlFile(dbConfig, configFilePath)
	if err != nil {
		return err
	}

	return nil
}

// SaveTomlFile will open and save data to toml file
// TODO: move to core
func SaveTomlFile(dest interface{}, relativePath string) error {
	f, err := os.Create(relativePath)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	return toml.NewEncoder(f).Encode(dest)
}

func (pf *PersisterFactory) getPersisterConfigFilePath(path string) string {
	return filepath.Join(
		path,
		dbConfigFileName,
	)
}

func (pf *PersisterFactory) createShardIDProvider() (storage.ShardIDProvider, error) {
	switch storageunit.ShardIDProviderType(pf.shardIDProviderType) {
	case storageunit.BinarySplit:
		return database.NewShardIDProvider(pf.numShards)
	default:
		return nil, storage.ErrNotSupportedShardIDProviderType
	}
}

// CreateDisabled will return a new disabled persister
func (pf *PersisterFactory) CreateDisabled() storage.Persister {
	return &disabledPersister{}
}

// IsInterfaceNil returns true if there is no value under the interface
func (pf *PersisterFactory) IsInterfaceNil() bool {
	return pf == nil
}

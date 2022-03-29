package mock

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/ElrondNetwork/elrond-go/common"
)

// StorerMock -
type StorerMock struct {
	mut  sync.Mutex
	data map[string][]byte
}

// NewStorerMock -
func NewStorerMock() *StorerMock {
	return &StorerMock{
		data: make(map[string][]byte),
	}
}

// GetFromEpoch -
func (sm *StorerMock) GetFromEpoch(key []byte, _ uint32, priority common.StorageAccessType) ([]byte, error) {
	return sm.Get(key, priority)
}

// GetBulkFromEpoch -
func (sm *StorerMock) GetBulkFromEpoch(keys [][]byte, _ uint32, priority common.StorageAccessType) (map[string][]byte, error) {
	retValue := map[string][]byte{}
	for _, key := range keys {
		value, err := sm.Get(key, priority)
		if err != nil {
			continue
		}
		retValue[string(key)] = value
	}

	return retValue, nil
}

// Put -
func (sm *StorerMock) Put(key, data []byte, _ common.StorageAccessType) error {
	sm.mut.Lock()
	defer sm.mut.Unlock()
	sm.data[string(key)] = data

	return nil
}

// PutInEpoch -
func (sm *StorerMock) PutInEpoch(key, data []byte, _ uint32, priority common.StorageAccessType) error {
	return sm.Put(key, data, priority)
}

// Get -
func (sm *StorerMock) Get(key []byte, _ common.StorageAccessType) ([]byte, error) {
	sm.mut.Lock()
	defer sm.mut.Unlock()

	val, ok := sm.data[string(key)]
	if !ok {
		return nil, fmt.Errorf("key: %s not found", base64.StdEncoding.EncodeToString(key))
	}

	return val, nil
}

// SearchFirst -
func (sm *StorerMock) SearchFirst(_ []byte, _ common.StorageAccessType) ([]byte, error) {
	return nil, errors.New("not implemented")
}

// Close -
func (sm *StorerMock) Close() error {
	return nil
}

// Has -
func (sm *StorerMock) Has(_ []byte, _ common.StorageAccessType) error {
	return errors.New("not implemented")
}

// RemoveFromCurrentEpoch -
func (sm *StorerMock) RemoveFromCurrentEpoch(_ []byte, _ common.StorageAccessType) error {
	return errors.New("not implemented")
}

// Remove -
func (sm *StorerMock) Remove(_ []byte, _ common.StorageAccessType) error {
	return errors.New("not implemented")
}

// GetOldestEpoch -
func (sm *StorerMock) GetOldestEpoch() (uint32, error) {
	return 0, nil
}

// ClearCache -
func (sm *StorerMock) ClearCache() {
}

// DestroyUnit -
func (sm *StorerMock) DestroyUnit() error {
	return nil
}

// RangeKeys -
func (sm *StorerMock) RangeKeys(handler func(key []byte, val []byte) bool) {
	if handler == nil {
		return
	}

	sm.mut.Lock()
	defer sm.mut.Unlock()

	for k, v := range sm.data {
		shouldContinue := handler([]byte(k), v)
		if !shouldContinue {
			return
		}
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (sm *StorerMock) IsInterfaceNil() bool {
	return sm == nil
}

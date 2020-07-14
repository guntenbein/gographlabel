package block_repository

import (
	"github.com/guntenbein/gographlabel"
	"sync"
)

type BlockProvider interface {
	ReadBlocks(hierarchyId string) (blocks []gographlabel.BlockOrder, err error)
	WriteBlocks(hierarchyId string, blocks []gographlabel.BlockOrder) error
}

type InMemoryBlockRepository struct {
	storage map[string][]gographlabel.BlockOrder
	mutex   sync.RWMutex
}

func (mbr InMemoryBlockRepository) ReadBlocks(hierarchyId string) (blocks []gographlabel.BlockOrder, err error) {
	mbr.mutex.RLock()
	defer mbr.mutex.RUnlock()
	block, ok := mbr.storage[hierarchyId]
	if !ok || block == nil {
		return []gographlabel.BlockOrder{}, nil
	}
	return block, nil
}

func (mbr InMemoryBlockRepository) WriteBlocks(hierarchyId string, blocks []gographlabel.BlockOrder) error {
	mbr.mutex.Lock()
	defer mbr.mutex.Unlock()
	mbr.storage[hierarchyId] = blocks
	return nil
}

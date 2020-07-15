package hierarchylock

import "sync"

type SingleServiceLock struct {
	lockMap map[string]*sync.Mutex
	mutex   sync.Mutex
}

func NewSingleServiceLock() *SingleServiceLock {
	return &SingleServiceLock{make(map[string]*sync.Mutex), sync.Mutex{}}
}

func (ssl *SingleServiceLock) LockHierarchy(hierarchyId string) (err error) {
	lock := ssl.getLockFor(hierarchyId)
	lock.Lock()
	return nil
}

func (ssl *SingleServiceLock) UnlockHierarchy(hierarchyId string) (err error) {
	lock := ssl.getLockFor(hierarchyId)
	lock.Unlock()
	return nil
}

func (ssl *SingleServiceLock) getLockFor(hierarchyId string) *sync.Mutex {
	ssl.mutex.Lock()
	defer ssl.mutex.Unlock()
	lock, ok := ssl.lockMap[hierarchyId]
	if !ok {
		lock = &sync.Mutex{}
		ssl.lockMap[hierarchyId] = lock
	}
	return lock
}

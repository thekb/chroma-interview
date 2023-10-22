package chromainterview

import "time"

type Store interface {
	Set(key string, value int)
	Get(key string) int
	Delete(key string)
}

type AtomicStore interface {
	Begin() AtomicStore
	Commit()
	Rollback()
	Store
}

type store struct {
	data map[string]int
}

var _ Store = (*store)(nil)

func newStore() *store {
	return &store{
		data: make(map[string]int),
	}
}

func (s *store) Set(key string, value int) {
	s.data[key] = value
}

func (s *store) Get(key string) int {
	if v, ok := s.data[key]; ok {
		return v
	}
	return -1
}

func (s *store) Delete(key string) {
	delete(s.data, key)
}

type pendingOp struct {
	opType  string
	opKey   string
	opValue int
}

type atomicStore struct {
	tid      int64
	pending  []pendingOp
	s        *store
	previous *atomicStore
}

var _ AtomicStore = (*atomicStore)(nil)
var _ Store = (*atomicStore)(nil)

func newAtomicStore() *atomicStore {
	return &atomicStore{
		s: newStore(),
	}
}

func (as *atomicStore) Set(key string, value int) {
	// if transaction is active append to the log
	if as.tid > 0 {
		as.pending = append(as.pending, pendingOp{
			opType:  "set",
			opKey:   key,
			opValue: value,
		})
	} else {
		as.s.Set(key, value)
	}
}

func (as *atomicStore) Get(key string) int {
	// read the pending operations log in reverse, as new operation
	// takes precedence over old operation
	for i := len(as.pending) - 1; i >= 0; i-- {
		op := as.pending[i]
		// if there is a delete operation return -1
		if op.opType == "delete" && op.opKey == key {
			return -1
		}
		// if there is a set operation return the value
		if op.opType == "set" && op.opKey == key {
			return op.opValue
		}
	}

	// in case of miss
	// if it is a nested transaction
	// query the parent
	if as.previous != nil {
		return as.previous.Get(key)
	}
	// if not a nested transaction
	// query the store
	return as.s.Get(key)
}

func (as *atomicStore) Delete(key string) {
	// if transaction is active append to log
	if as.tid > 0 {
		as.pending = append(as.pending, pendingOp{
			opType: "delete",
			opKey:  key,
		})
	} else {
		as.s.Delete(key)
	}

}

func (as *atomicStore) Begin() AtomicStore {
	tid := time.Now().UnixNano()
	return &atomicStore{
		tid:      tid,
		s:        as.s,
		previous: as,
	}
}

func (s *atomicStore) Commit() {
	for _, po := range s.pending {
		switch po.opType {
		case "set":
			s.s.Set(po.opKey, po.opValue)
		case "delete":
			s.s.Delete(po.opKey)
		}
	}
}

func (s *atomicStore) Rollback() {
	s.pending = nil
}

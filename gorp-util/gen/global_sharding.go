package main

import (
	"sync"
	"sync/atomic"

	"gopkg.in/gorp.v2"
)

// @TABLE(global_sharding)
type GlobalSharding struct {
	ID         int64  `db:"id"`
	UpdatedSeq int64  `db:"updated_seq"` // 用于乐观锁
	Table      string `db:"table"`
	CurrentID  int64  `db:"current_id"` // 当前已被获取的最大ID
	Batch      int64  `db:"batch"`      // 一次获取多少个
	Shard      int64  `db:"shard"`      // 多少ID一个shard
}

func (GlobalSharding) TableName() string {
	return "global_sharding"
}

type GlobalIDGenerator struct {
	hold sync.Map
	db   *gorp.DbMap
}

type idLocker struct {
	sync.RWMutex
	counts int64
	col    *GlobalSharding
}

func NewGlobalIDGenerator(db *gorp.DbMap) (*GlobalIDGenerator, error) {
	res := &GlobalIDGenerator{
		db: db,
	}

	var tmp []*GlobalSharding
	if _, err := (&GlobalSharding{}).Where().FetchAll(db, &tmp); err != nil {
		return nil, err
	}
	for _, v := range tmp {
		res.hold.Store(v.Table, &idLocker{counts: v.Batch})
	}
	return res, nil
}

func (g *GlobalIDGenerator) Get(tablename string) (int64, error) {
	if v, ok := g.hold.Load(tablename); ok {
		l := v.(*idLocker)
		for {
			l.RLock()
			val := atomic.AddInt64(&l.counts, 1)
			l.RUnlock()
			if val > l.col.Batch {
				if err := g.fetchBatch(l); err != nil {
					return 0, err
				}
				// 重新获取
			} else {
				return val, nil
			}
		}

	} else {
		panic("not found table in global_sharding:" + tablename)
	}
}

func (g *GlobalIDGenerator) fetchBatch(l *idLocker) error {
	l.Lock()
	defer l.Unlock()
	if l.counts < l.col.Batch {
		// 已经被其他协程获取了
		return nil
	}
	// 极限情况下， l.counts == l.col.Batch,
	// 也就是说，当前协程获取的时候，其他协程已经消费了追加的一批，还是需要重新获取

	for {
		if err := l.col.Load(g.db, l.col.ID); err != nil {
			return err
		}
		l.col.CurrentID += l.col.Batch
		if err := l.col.Update(g.db); err != nil {
			if _, ok := err.(gorp.OptimisticLockError); ok {
				// 乐观锁失败，重新获取
				continue
			}
			// 其他错误
			return err
		}
		break
	}

	// 更新成功
	l.counts = 0
	return nil

}

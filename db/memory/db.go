package memory

import (
	dbm "github.com/cnaize/mz-center/db"
	"github.com/cnaize/mz-center/log"
	"github.com/cnaize/mz-center/model"
	"sync"
	"time"
)

type DB struct {
	mu       sync.RWMutex
	searches map[model.SearchRequest]map[model.User]model.SearchResponseList
}

func NewDB() *DB {
	db := &DB{
		searches: make(map[model.SearchRequest]map[model.User]model.SearchResponseList),
	}

	go db.gc()

	return db
}

func (db *DB) GetSearchRequestList(offset, count int) (model.SearchRequestList, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if count == 0 {
		count = dbm.MaxItemsCount
	} else if count > dbm.MaxItemsCount {
		count = dbm.MaxItemsCount
	}

	var res model.SearchRequestList
	var i int
	for r := range db.searches {
		if i >= count {
			break
		}

		res.Items = append(res.Items, r)
		i++
	}

	return res, nil
}

func (db *DB) AddSearchRequest(request model.SearchRequest) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.searches[request]; !ok {
		db.searches[request] = make(map[model.User]model.SearchResponseList)
	}

	return nil
}

func (db *DB) GetSearchResponseList(request model.SearchRequest, offset, count int) (model.SearchResponseList, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if count == 0 {
		count = dbm.MaxItemsCount
	} else if count > dbm.MaxItemsCount {
		count = dbm.MaxItemsCount
	}

	var res model.SearchResponseList
	uresp, ok := db.searches[request]
	if !ok {
		return res, dbm.ErrorNotFound
	}

	for _, r := range uresp {
		if len(res.Items) >= count {
			break
		}

		res.Items = append(res.Items, r.Items...)
	}

	return res, nil
}

func (db *DB) AddSearchResponseList(user model.User, request model.SearchRequest, response model.SearchResponseList) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if uresp, ok := db.searches[request]; !ok {
		return nil
	} else if _, ok := uresp[user]; ok {
		return nil
	}

	for _, r := range response.Items {
		r.Owner = user
	}

	db.searches[request][user] = response

	return nil
}

func (db *DB) gc() {
	gc := func() {
		db.mu.Lock()
		defer db.mu.Unlock()

		var i int
		for r := range db.searches {
			if time.Since(r.CreatedAt) > 8*time.Second {
				delete(db.searches, r)
				i++
			}
		}

		if i > 0 {
			log.Debug("MemoryDB.gc(): search requests cleaned: %d", i)
		}
	}

	for {
		time.Sleep(1 * time.Second)
		gc()
	}
}

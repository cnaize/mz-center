package memory

import (
	dbm "github.com/cnaize/mz-center/db"
	"github.com/cnaize/mz-center/log"
	"github.com/cnaize/mz-center/model"
	"sync"
	"time"
)

type DB struct {
	ms         sync.RWMutex
	searches   map[string]map[*model.User]*model.SearchResponseList
	mr         sync.Mutex
	searchReqs []*model.SearchRequest
}

func NewDB() *DB {
	db := &DB{
		searches: make(map[string]map[*model.User]*model.SearchResponseList),
	}

	go db.gc()

	return db
}

func (db *DB) GetSearchRequestList(offset, count uint) (*model.SearchRequestList, error) {
	db.ms.RLock()
	defer db.ms.RUnlock()

	res := &model.SearchRequestList{TotalReqCount: uint(len(db.searches))}
	if res.TotalReqCount < offset {
		return res, nil
	}

	if count == 0 {
		count = dbm.MaxItemsCount
	} else if count > dbm.MaxItemsCount {
		count = dbm.MaxItemsCount
	}

	var o uint
	var c uint
	for s := range db.searches {
		if o < offset {
			continue
		}
		if c >= count {
			break
		}

		res.Items = append(res.Items, &model.SearchRequest{Text: s})
		o++
		c++
	}

	return res, nil
}

func (db *DB) AddSearchRequest(request *model.SearchRequest) error {
	db.ms.Lock()
	defer db.ms.Unlock()

	if _, ok := db.searches[request.Text]; !ok {
		db.mr.Lock()
		defer db.mr.Unlock()

		request.CreatedAt = time.Now()
		db.searchReqs = append(db.searchReqs, request)
		db.searches[request.Text] = make(map[*model.User]*model.SearchResponseList)
	}

	return nil
}

func (db *DB) GetSearchResponseList(request *model.SearchRequest, offset, count uint) (*model.SearchResponseList, error) {
	db.ms.RLock()
	defer db.ms.RUnlock()

	uresp, ok := db.searches[request.Text]
	if !ok {
		return nil, dbm.ErrorNotFound
	}

	if count == 0 {
		count = dbm.MaxItemsCount
	} else if count > dbm.MaxItemsCount {
		count = dbm.MaxItemsCount
	}

	res := &model.SearchResponseList{}
	var o uint
	var t uint
	for _, r := range uresp {
		t += uint(len(r.Items))

		if o < offset {
			continue
		}
		if uint(len(res.Items)) >= count {
			continue
		}

		res.Items = append(res.Items, r.Items...)
		o++
	}

	res.Request = &model.SearchRequest{TotalRespCount: t}

	return res, nil
}

func (db *DB) AddSearchResponseList(user *model.User, request *model.SearchRequest, response *model.SearchResponseList) error {
	db.ms.Lock()
	defer db.ms.Unlock()

	if uresp, ok := db.searches[request.Text]; !ok {
		return nil
	} else if _, ok := uresp[user]; ok {
		return nil
	}

	for _, r := range response.Items {
		r.Owner = user
	}

	db.searches[request.Text][user] = response

	return nil
}

func (db *DB) gc() {
	gc := func() {
		db.ms.Lock()
		db.mr.Lock()
		defer db.mr.Unlock()
		defer db.ms.Unlock()

		now := time.Now()
		var i int
		for ; i < len(db.searchReqs); i++ {
			sr := db.searchReqs[i]

			if now.Sub(sr.CreatedAt) < 8*time.Second {
				break
			}

			delete(db.searches, sr.Text)
		}

		db.searchReqs = db.searchReqs[i:]

		if i > 0 {
			log.Debug("MemoryDB.gc(): search requests cleaned: %d", i)
		}
	}

	for {
		time.Sleep(1 * time.Second)
		gc()
	}
}

package memory

import (
	dbm "github.com/cnaize/mz-center/db"
	"github.com/cnaize/mz-common/log"
	"github.com/cnaize/mz-common/model"
	"sync"
	"time"
)

type DB struct {
	ms         sync.RWMutex
	searches   map[string]map[*model.User]*model.SearchResponseList
	mr         sync.RWMutex
	searchReqs []*model.SearchRequest
}

func NewDB() *DB {
	db := &DB{
		searches: make(map[string]map[*model.User]*model.SearchResponseList),
	}

	//go db.gc()

	return db
}

func (db *DB) GetSearchRequestList(offset, count uint) (*model.SearchRequestList, error) {
	db.ms.RLock()
	db.mr.RLock()
	defer db.ms.RUnlock()
	defer db.mr.RUnlock()

	res := model.SearchRequestList{TotalReqCount: uint(len(db.searchReqs))}
	if res.TotalReqCount < offset {
		return &res, nil
	}

	if count == 0 || count > dbm.MaxRequestItemsCount {
		count = dbm.MaxRequestItemsCount
	}

	var o uint
	var c uint
	for _, r := range db.searchReqs {
		if r.TotalRespCount >= dbm.MaxResponseItemsCount {
			continue
		}
		if o < offset {
			continue
		}
		if c >= count {
			break
		}

		res.Items = append(res.Items, r)
		o++
		c++
	}

	return &res, nil
}

func (db *DB) AddSearchRequest(request *model.SearchRequest) (*model.SearchRequest, error) {
	db.ms.Lock()
	db.mr.Lock()
	defer db.ms.Unlock()
	defer db.mr.Unlock()

	if _, ok := db.searches[request.Text]; !ok {
		db.searchReqs = append(db.searchReqs, request)
		db.searches[request.Text] = make(map[*model.User]*model.SearchResponseList)
	} else {
		request = db.findRequest(request)
	}

	request.UpdatedAt = time.Now()

	return request, nil
}

func (db *DB) GetSearchResponseList(request *model.SearchRequest, offset, count uint) (*model.SearchResponseList, error) {
	db.ms.RLock()
	db.mr.RLock()
	defer db.ms.RUnlock()
	defer db.mr.RUnlock()

	uresp, ok := db.searches[request.Text]
	if !ok {
		return nil, dbm.ErrorNotFound
	}

	if count == 0 || count > dbm.MaxResponseItemsCount {
		count = dbm.MaxResponseItemsCount
	}

	var res model.SearchResponseList
	var o uint
	for _, r := range uresp {
		if o < offset {
			continue
		}
		if uint(len(res.Items)) >= count {
			break
		}

		res.Items = append(res.Items, r.Items...)
		o++
	}

	res.Request = db.findRequest(request)

	return &res, nil
}

func (db *DB) AddSearchResponseList(user *model.User, request *model.SearchRequest, response *model.SearchResponseList) error {
	db.ms.Lock()
	db.mr.RLock()
	defer db.ms.Unlock()
	defer db.mr.RUnlock()

	request = db.findRequest(request)
	if request == nil {
		return nil
	}

	if request.TotalRespCount >= dbm.MaxResponseItemsCount {
		return nil
	}

	if uresp, ok := db.searches[request.Text]; !ok {
		return nil
	} else if _, ok := uresp[user]; ok {
		return nil
	}

	for _, r := range response.Items {
		r.Owner = user
	}

	db.searches[request.Text][user] = response
	request.TotalRespCount = request.TotalRespCount + uint(len(response.Items))

	return nil
}

func (db *DB) findRequest(inRequest *model.SearchRequest) *model.SearchRequest {
	for _, r := range db.searchReqs {
		if r.Text == inRequest.Text {
			return r
		}
	}

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

			if now.Sub(sr.UpdatedAt) < dbm.RequestsCleanPeriod {
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

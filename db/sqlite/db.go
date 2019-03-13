package sqlite

import (
	"fmt"
	"github.com/cnaize/mz-common/log"
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"sync"
	"time"
)

const (
	requestsCleanPeriod = 10 * time.Second
)

type DB struct {
	mu sync.Mutex
	db *gorm.DB
}

func New() (*DB, error) {
	conn, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("open failed: %+v", err)
	}

	if err := prepare(conn); err != nil {
		return nil, fmt.Errorf("prepare failed: %+v", err)
	}

	db := DB{
		db: conn,
	}

	go db.runGC()

	return &db, nil
}

func (db *DB) runGC() {
	gc := func() {
		db.mu.Lock()
		defer db.mu.Unlock()

		line := time.Now().Add(-requestsCleanPeriod)
		var reqIDs []string
		if err := db.db.Model(&model.SearchRequest{}).Where("updated_at < ?", line).Pluck("id", &reqIDs).Error; err != nil {
			return
		}

		if len(reqIDs) == 0 {
			return
		}

		db.db.Delete(&model.SearchResponse{}, "search_request_id IN (?)", reqIDs)
		db.db.Delete(&model.SearchRequest{}, "id IN (?)", reqIDs)

		log.Debug("DB.gc(): requests removed: %d", len(reqIDs))
	}

	for {
		gc()
		time.Sleep(1 * time.Second)
	}
}

func prepare(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.User{}).Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.Media{}).Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.MediaRequest{}).Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.MediaResponse{}).Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.SearchRequest{}).Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.SearchResponse{}).Error; err != nil {
		return err
	}

	return nil
}

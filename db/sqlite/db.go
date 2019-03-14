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
	rSearchCleanPeriod = 10 * time.Second
	rMediaCleanPeriod  = 10 * time.Second
	mediaCleanPeriod   = 10 * time.Second
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
	rSearchGC := func() {
		db.mu.Lock()
		defer db.mu.Unlock()

		line := time.Now().Add(-rSearchCleanPeriod)

		var reqIDs []string
		if err := db.db.Model(&model.SearchRequest{}).Where("updated_at < ?", line).Pluck("id", &reqIDs).Error; err != nil {
			return
		}

		if len(reqIDs) == 0 {
			return
		}

		queue := db.db.Delete(&model.SearchResponse{}, "search_request_id IN (?)", reqIDs)
		if queue.RowsAffected > 0 {
			log.Debug("DB.rSearchGC(): requests removed: %d", queue.RowsAffected)
		}

		queue = db.db.Delete(&model.SearchRequest{}, "id IN (?)", reqIDs)
		if queue.RowsAffected > 0 {
			log.Debug("DB.rSearchGC(): responses removed: %d", queue.RowsAffected)
		}
	}

	rMediaGC := func() {
		db.mu.Lock()
		defer db.mu.Unlock()

		line := time.Now().Add(-rMediaCleanPeriod)

		queue := db.db.Delete(&model.MediaRequest{}, "updated_at < ?", line)
		if queue.RowsAffected > 0 {
			log.Debug("DB.rMediaGC(): requests removed: %d", queue.RowsAffected)
		}

		queue = db.db.Delete(&model.MediaResponse{}, "updated_at < ?", line)
		if queue.RowsAffected > 0 {
			log.Debug("DB.rMediaGC(): responses removed: %d", queue.RowsAffected)
		}
	}

	mediaGC := func() {
		db.mu.Lock()
		defer db.mu.Unlock()

		line := time.Now().Add(-mediaCleanPeriod)

		var srMediaIDs []uint
		db.db.Model(&model.Media{}).Joins("INNER JOIN search_responses ON search_responses.media_id = media.id").
			Where("media.updated_at < ?", line).
			Pluck("media.id", &srMediaIDs)
		var mrMediaIDs []uint
		db.db.Model(&model.Media{}).Joins("INNER JOIN media_responses ON media_responses.media_id = media.id").
			Where("media.updated_at < ?", line).
			Pluck("media.id", &mrMediaIDs)

		queue := db.db.Delete(&model.Media{}).Not("id", append(srMediaIDs, mrMediaIDs...))
		if queue.RowsAffected > 0 {
			log.Debug("DB.mediaGC(): media removed: %d", queue.RowsAffected)
		}
	}

	rSearchTicker := time.NewTicker(2 * time.Second)
	rMediaTicker := time.NewTicker(3 * time.Second)
	mediaTicker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-rSearchTicker.C:
			rSearchGC()
		case <-rMediaTicker.C:
			rMediaGC()
		case <-mediaTicker.C:
			mediaGC()
		}
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

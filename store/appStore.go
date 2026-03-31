package store

import (
	"github.com/bartek5186/waitlistfox/internal"
	"gorm.io/gorm"
)

type Datastore interface {
	Db() *gorm.DB
	Config() *internal.Config
	Waitlist() *WaitlistStore
}

type AppStore struct {
	db       *gorm.DB
	config   *internal.Config
	waitlist *WaitlistStore
}

func NewAppStore(db *gorm.DB, cfg *internal.Config) *AppStore {
	return &AppStore{
		db:       db,
		config:   cfg,
		waitlist: NewWaitlistStore(db),
	}
}

func (s *AppStore) Db() *gorm.DB {
	return s.db
}

func (s *AppStore) Config() *internal.Config {
	return s.config
}

func (s *AppStore) Waitlist() *WaitlistStore {
	return s.waitlist
}

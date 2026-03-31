package services

import (
	"github.com/bartek5186/waitlistfox/store"
	"github.com/sirupsen/logrus"
)

type AppService struct {
	Store    store.Datastore
	Waitlist *WaitlistService
	logger   *logrus.Logger
}

func NewAppService(store store.Datastore, logger *logrus.Logger) *AppService {
	return &AppService{
		Store:    store,
		Waitlist: NewWaitlistService(store, logger),
		logger:   logger,
	}
}

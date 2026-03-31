package store

import (
	"context"
	"strings"
	"time"

	"github.com/bartek5186/waitlistfox/models"
	"gorm.io/gorm"
)

type WaitlistStore struct {
	db *gorm.DB
}

func NewWaitlistStore(db *gorm.DB) *WaitlistStore {
	return &WaitlistStore{db: db}
}

func (s *WaitlistStore) CreateSubscriber(ctx context.Context, subscriber *models.WaitlistSubscriber) error {
	return s.db.WithContext(ctx).Create(subscriber).Error
}

func (s *WaitlistStore) UnsubscribeByEmail(ctx context.Context, email string, reason string, meta models.WaitlistRequestMeta, clientTimestamp *time.Time) (*models.WaitlistUnsubscribe, error) {
	var event models.WaitlistUnsubscribe

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var subscriber models.WaitlistSubscriber
		if err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(email))).First(&subscriber).Error; err != nil {
			return err
		}

		event = models.WaitlistUnsubscribe{
			SubscriberID:    &subscriber.ID,
			Email:           subscriber.Email,
			UserType:        subscriber.UserType,
			CampaignID:      subscriber.CampaignID,
			Reason:          reason,
			UserAgent:       strings.TrimSpace(meta.UserAgent),
			Referrer:        strings.TrimSpace(meta.Referrer),
			IPAddress:       strings.TrimSpace(meta.IPAddress),
			ClientTimestamp: clientTimestamp,
		}

		if err := tx.Create(&event).Error; err != nil {
			return err
		}

		return tx.Delete(&subscriber).Error
	})
	if err != nil {
		return nil, err
	}

	return &event, nil
}

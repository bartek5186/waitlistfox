package models

import "time"

const (
	WaitlistUserTypePassenger      = "passenger"
	WaitlistUserTypeCarrier        = "carrier"
	WaitlistReasonFoundAlternative = "found_alternative"
	WaitlistReasonTooManyEmails    = "too_many_emails"
	WaitlistReasonNotInterested    = "not_interested"
	WaitlistReasonPrivacyConcerns  = "privacy_concerns"
	WaitlistReasonMistake          = "signed_up_by_mistake"
	WaitlistReasonOther            = "other"
)

type WaitlistSubscriber struct {
	ID                      uint       `gorm:"primaryKey" json:"id"`
	UserType                string     `gorm:"size:32;not null;index:idx_waitlist_user_type" json:"user_type"`
	CampaignID              string     `gorm:"size:100;not null;index:idx_waitlist_campaign_id" json:"campaign_id"`
	Email                   string     `gorm:"size:255;not null;uniqueIndex:ux_waitlist_email" json:"email"`
	Phone                   string     `gorm:"size:50" json:"phone,omitempty"`
	ConsentWaitlist         bool       `gorm:"not null;default:true" json:"consent_waitlist"`
	ConsentMarketing        bool       `gorm:"not null;default:false" json:"consent_marketing"`
	ConsentCookiesAnalytics bool       `gorm:"not null;default:false" json:"consent_cookies_analytics"`
	ConsentCookiesMarketing bool       `gorm:"not null;default:false" json:"consent_cookies_marketing"`
	RecaptchaScore          *float64   `gorm:"type:decimal(3,2)" json:"recaptcha_score,omitempty"`
	UserAgent               string     `gorm:"type:text" json:"user_agent,omitempty"`
	Referrer                string     `gorm:"size:500" json:"referrer,omitempty"`
	IPAddress               string     `gorm:"size:45" json:"ip_address,omitempty"`
	ClientTimestamp         *time.Time `json:"client_timestamp,omitempty"`
	CreatedAt               time.Time  `gorm:"not null;autoCreateTime;index:idx_waitlist_created_at" json:"created_at"`
	NotifiedAt              *time.Time `json:"notified_at,omitempty"`
}

func (WaitlistSubscriber) TableName() string {
	return "waitlist_subscribers"
}

type WaitlistUnsubscribe struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	SubscriberID    *uint      `json:"subscriber_id,omitempty"`
	Email           string     `gorm:"size:255;not null;index:idx_waitlist_unsubscribes_email" json:"email"`
	UserType        string     `gorm:"size:32" json:"user_type,omitempty"`
	CampaignID      string     `gorm:"size:100;index:idx_waitlist_unsubscribes_campaign_id" json:"campaign_id,omitempty"`
	Reason          string     `gorm:"size:64;not null" json:"reason"`
	UserAgent       string     `gorm:"type:text" json:"user_agent,omitempty"`
	Referrer        string     `gorm:"size:500" json:"referrer,omitempty"`
	IPAddress       string     `gorm:"size:45" json:"ip_address,omitempty"`
	ClientTimestamp *time.Time `json:"client_timestamp,omitempty"`
	CreatedAt       time.Time  `gorm:"not null;autoCreateTime;index:idx_waitlist_unsubscribes_created_at" json:"created_at"`
}

func (WaitlistUnsubscribe) TableName() string {
	return "waitlist_unsubscribes"
}

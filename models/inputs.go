package models

import "time"

type WaitlistConsentsInput struct {
	Waitlist         bool `json:"waitlist"`
	Marketing        bool `json:"marketing"`
	CookiesAnalytics bool `json:"cookies_analytics"`
	CookiesMarketing bool `json:"cookies_marketing"`
}

type WaitlistSubscribeInput struct {
	UserType       string                `json:"userType" validate:"required,oneof=passenger carrier"`
	CampaignID     string                `json:"campaignId,omitempty" validate:"omitempty,max=100"`
	Email          string                `json:"email" validate:"required,email,max=255"`
	Phone          string                `json:"phone" validate:"omitempty,max=50"`
	Consents       WaitlistConsentsInput `json:"consents"`
	RecaptchaToken string                `json:"recaptcha_token" validate:"omitempty,max=4096"`
	Timestamp      *time.Time            `json:"timestamp,omitempty"`
	UserAgent      string                `json:"user_agent" validate:"omitempty,max=4000"`
	Referrer       string                `json:"referrer" validate:"omitempty,httpurl,max=500"`
}

type WaitlistRequestMeta struct {
	IPAddress string
	UserAgent string
	Referrer  string
}

type WaitlistUnsubscribeInput struct {
	Email     string     `json:"email" validate:"required,email,max=255"`
	Reason    string     `json:"reason" validate:"required,oneof=found_alternative too_many_emails not_interested privacy_concerns signed_up_by_mistake other"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	UserAgent string     `json:"user_agent" validate:"omitempty,max=4000"`
	Referrer  string     `json:"referrer" validate:"omitempty,httpurl,max=500"`
}

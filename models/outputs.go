package models

type HealthResponse struct {
	Status string `json:"status"`
	App    string `json:"app"`
}

type WaitlistSubscribeResponse struct {
	App            string   `json:"app"`
	Status         string   `json:"status"`
	Message        string   `json:"message"`
	SubscriberID   uint     `json:"subscriber_id"`
	Email          string   `json:"email"`
	UserType       string   `json:"user_type"`
	CampaignID     string   `json:"campaign_id"`
	RecaptchaScore *float64 `json:"recaptcha_score,omitempty"`
}

type WaitlistUnsubscribeResponse struct {
	App     string `json:"app"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Email   string `json:"email"`
	Reason  string `json:"reason"`
}

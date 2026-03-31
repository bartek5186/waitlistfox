package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bartek5186/waitlistfox/internal/i18n"
	"github.com/bartek5186/waitlistfox/models"
	"github.com/bartek5186/waitlistfox/store"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrWaitlistConsentRequired = errors.New("waitlist consent required")
	ErrEmailAlreadySubscribed  = errors.New("email already subscribed")
	ErrEmailNotSubscribed      = errors.New("email not subscribed")
	ErrRecaptchaTokenRequired  = errors.New("recaptcha token required")
	ErrRecaptchaRejected       = errors.New("recaptcha rejected")
	ErrRecaptchaActionMismatch = errors.New("recaptcha action mismatch")
	ErrRecaptchaMisconfigured  = errors.New("recaptcha misconfigured")
)

type WaitlistService struct {
	Store         store.Datastore
	logger        *logrus.Logger
	recaptchaHTTP *http.Client
}

type recaptchaVerifyResponse struct {
	Success    bool     `json:"success"`
	Score      float64  `json:"score"`
	Action     string   `json:"action"`
	ErrorCodes []string `json:"error-codes"`
}

func NewWaitlistService(store store.Datastore, logger *logrus.Logger) *WaitlistService {
	timeout := time.Duration(store.Config().Recaptcha.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	return &WaitlistService{
		Store:  store,
		logger: logger,
		recaptchaHTTP: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *WaitlistService) Health(locale string) models.HealthResponse {
	return models.HealthResponse{
		Status: i18n.T(locale, "health.status", nil),
		App:    s.Store.Config().AppName,
	}
}

func (s *WaitlistService) Subscribe(ctx context.Context, locale string, in models.WaitlistSubscribeInput, meta models.WaitlistRequestMeta) (*models.WaitlistSubscribeResponse, error) {
	campaignID := strings.TrimSpace(in.CampaignID)
	if !in.Consents.Waitlist {
		return nil, ErrWaitlistConsentRequired
	}

	email := strings.ToLower(strings.TrimSpace(in.Email))
	phone := strings.TrimSpace(in.Phone)
	referrer := strings.TrimSpace(in.Referrer)
	if referrer == "" {
		referrer = strings.TrimSpace(meta.Referrer)
	}
	userAgent := strings.TrimSpace(in.UserAgent)
	if userAgent == "" {
		userAgent = strings.TrimSpace(meta.UserAgent)
	}

	score, err := s.verifyRecaptcha(ctx, strings.TrimSpace(in.RecaptchaToken), meta.IPAddress)
	if err != nil {
		return nil, err
	}

	subscriber := &models.WaitlistSubscriber{
		UserType:                in.UserType,
		CampaignID:              campaignID,
		Email:                   email,
		Phone:                   phone,
		ConsentWaitlist:         in.Consents.Waitlist,
		ConsentMarketing:        in.Consents.Marketing,
		ConsentCookiesAnalytics: in.Consents.CookiesAnalytics,
		ConsentCookiesMarketing: in.Consents.CookiesMarketing,
		RecaptchaScore:          score,
		UserAgent:               userAgent,
		Referrer:                referrer,
		IPAddress:               strings.TrimSpace(meta.IPAddress),
		ClientTimestamp:         in.Timestamp,
	}

	if err := s.Store.Waitlist().CreateSubscriber(ctx, subscriber); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailAlreadySubscribed
		}

		return nil, err
	}

	return &models.WaitlistSubscribeResponse{
		App:            s.Store.Config().AppName,
		Status:         "created",
		Message:        i18n.T(locale, "waitlist.subscribe.success", nil),
		SubscriberID:   subscriber.ID,
		Email:          subscriber.Email,
		UserType:       subscriber.UserType,
		CampaignID:     subscriber.CampaignID,
		RecaptchaScore: subscriber.RecaptchaScore,
	}, nil
}

func (s *WaitlistService) Unsubscribe(ctx context.Context, locale string, in models.WaitlistUnsubscribeInput, meta models.WaitlistRequestMeta) (*models.WaitlistUnsubscribeResponse, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	reason := strings.TrimSpace(in.Reason)
	referrer := strings.TrimSpace(in.Referrer)
	if referrer == "" {
		referrer = strings.TrimSpace(meta.Referrer)
	}
	userAgent := strings.TrimSpace(in.UserAgent)
	if userAgent == "" {
		userAgent = strings.TrimSpace(meta.UserAgent)
	}

	event, err := s.Store.Waitlist().UnsubscribeByEmail(ctx, email, reason, models.WaitlistRequestMeta{
		IPAddress: strings.TrimSpace(meta.IPAddress),
		UserAgent: userAgent,
		Referrer:  referrer,
	}, in.Timestamp)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmailNotSubscribed
		}

		return nil, err
	}

	return &models.WaitlistUnsubscribeResponse{
		App:     s.Store.Config().AppName,
		Status:  "deleted",
		Message: i18n.T(locale, "waitlist.unsubscribe.success", nil),
		Email:   event.Email,
		Reason:  event.Reason,
	}, nil
}

func (s *WaitlistService) verifyRecaptcha(ctx context.Context, token string, ipAddress string) (*float64, error) {
	cfg := s.Store.Config().Recaptcha
	if !cfg.Enabled {
		return nil, nil
	}

	if strings.TrimSpace(cfg.SecretKey) == "" {
		return nil, ErrRecaptchaMisconfigured
	}
	if token == "" {
		return nil, ErrRecaptchaTokenRequired
	}

	verifyURL := strings.TrimSpace(cfg.VerifyURL)
	if verifyURL == "" {
		verifyURL = "https://www.google.com/recaptcha/api/siteverify"
	}

	form := url.Values{}
	form.Set("secret", cfg.SecretKey)
	form.Set("response", token)
	if ipAddress != "" {
		form.Set("remoteip", ipAddress)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("recaptcha request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.recaptchaHTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("recaptcha call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recaptcha status: %s", resp.Status)
	}

	var verifyResp recaptchaVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("recaptcha decode: %w", err)
	}

	if !verifyResp.Success {
		s.logger.WithField("error_codes", strings.Join(verifyResp.ErrorCodes, ",")).Warn("recaptcha verify rejected")
		return nil, ErrRecaptchaRejected
	}

	expectedAction := strings.TrimSpace(cfg.ExpectedAction)
	if expectedAction != "" && verifyResp.Action != expectedAction {
		return nil, ErrRecaptchaActionMismatch
	}

	minScore := cfg.MinimumScore
	if minScore <= 0 {
		minScore = 0.5
	}
	if verifyResp.Score < minScore {
		return nil, ErrRecaptchaRejected
	}

	score := verifyResp.Score
	return &score, nil
}

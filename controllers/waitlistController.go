package controllers

import (
	"errors"
	"net/http"

	"github.com/bartek5186/waitlistfox/internal/i18n"
	"github.com/bartek5186/waitlistfox/models"
	"github.com/bartek5186/waitlistfox/services"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type WaitlistController struct {
	appService *services.AppService
	logger     *logrus.Logger
}

func NewWaitlistController(appService *services.AppService, logger *logrus.Logger) *WaitlistController {
	return &WaitlistController{
		appService: appService,
		logger:     logger,
	}
}

func (c *WaitlistController) Health(ec echo.Context) error {
	locale := localeFromContext(ec)
	return ec.JSON(http.StatusOK, c.appService.Waitlist.Health(locale))
}

func (c *WaitlistController) Subscribe(ec echo.Context) error {
	ctx := ec.Request().Context()
	locale := localeFromContext(ec)

	var in models.WaitlistSubscribeInput
	if err := ec.Bind(&in); err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": i18n.T(locale, "waitlist.error.invalid_payload", nil)})
	}
	if err := ec.Validate(&in); err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	meta := models.WaitlistRequestMeta{
		IPAddress: ec.RealIP(),
		UserAgent: ec.Request().UserAgent(),
		Referrer:  ec.Request().Referer(),
	}

	out, err := c.appService.Waitlist.Subscribe(ctx, locale, in, meta)
	if err != nil {
		status := http.StatusInternalServerError
		message := i18n.T(locale, "waitlist.error.internal", nil)

		switch {
		case errors.Is(err, services.ErrWaitlistConsentRequired):
			status = http.StatusBadRequest
			message = i18n.T(locale, "waitlist.error.consent_required", nil)
		case errors.Is(err, services.ErrEmailAlreadySubscribed):
			status = http.StatusConflict
			message = i18n.T(locale, "waitlist.error.duplicate_email", nil)
		case errors.Is(err, services.ErrRecaptchaTokenRequired):
			status = http.StatusBadRequest
			message = i18n.T(locale, "waitlist.error.recaptcha_required", nil)
		case errors.Is(err, services.ErrRecaptchaRejected), errors.Is(err, services.ErrRecaptchaActionMismatch):
			status = http.StatusForbidden
			message = i18n.T(locale, "waitlist.error.recaptcha_rejected", nil)
		case errors.Is(err, services.ErrRecaptchaMisconfigured):
			status = http.StatusInternalServerError
			message = i18n.T(locale, "waitlist.error.recaptcha_misconfigured", nil)
		}

		c.logger.WithError(err).Error("waitlist subscribe failed")
		return ec.JSON(status, map[string]string{"error": message})
	}

	return ec.JSON(http.StatusCreated, out)
}

func (c *WaitlistController) Unsubscribe(ec echo.Context) error {
	ctx := ec.Request().Context()
	locale := localeFromContext(ec)

	var in models.WaitlistUnsubscribeInput
	if err := ec.Bind(&in); err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": i18n.T(locale, "waitlist.error.invalid_payload", nil)})
	}
	if err := ec.Validate(&in); err != nil {
		return ec.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	meta := models.WaitlistRequestMeta{
		IPAddress: ec.RealIP(),
		UserAgent: ec.Request().UserAgent(),
		Referrer:  ec.Request().Referer(),
	}

	out, err := c.appService.Waitlist.Unsubscribe(ctx, locale, in, meta)
	if err != nil {
		status := http.StatusInternalServerError
		message := i18n.T(locale, "waitlist.error.internal", nil)

		switch {
		case errors.Is(err, services.ErrEmailNotSubscribed):
			status = http.StatusNotFound
			message = i18n.T(locale, "waitlist.error.not_found", nil)
		}

		c.logger.WithError(err).Error("waitlist unsubscribe failed")
		return ec.JSON(status, map[string]string{"error": message})
	}

	return ec.JSON(http.StatusOK, out)
}

func localeFromContext(ec echo.Context) string {
	if lang, ok := ec.Get("lang").(string); ok && lang != "" {
		return lang
	}

	return "pl"
}

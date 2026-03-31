![waitlistFox logo](waitingfox-logo.png)

# waitlistFox

`waitlistFox` is a simple backend API for waitlist signup built on top of Echo and MySQL.

## Structure

```text
config/
controllers/
internal/
models/
services/
static/
store/
main.go
```

## What's included

- public `POST /v1/waitlist/subscribe` endpoint
- MySQL persistence for `waitlist_subscribers`
- duplicate email protection through a unique constraint
- consent fields, campaign ID, request metadata, and notification timestamp
- optional reCAPTCHA verification configured in `config/config.json`
- `internal/config.go` with configuration loading and MySQL connection setup
- `internal/logger.go` with a JSON logger writing to both file and stdout
- `internal/validator.go`
- `internal/middleware/language.go`
- `internal/i18n/` with a simple translation loader
- `static/index.html`

## Running

```bash
cd waitlistfox
go run . -migrate=true
```

By default, the application reads configuration from `config/config.json`.

## Endpoints

- `GET /health`
- `POST /v1/waitlist/subscribe`
- `POST /v1/waitlist/unsubscribe`

Example payload:

```json
{
  "userType": "passenger",
  "email": "john@example.com",
  "phone": "+48123456789",
  "consents": {
    "waitlist": true,
    "marketing": false,
    "cookies_analytics": true,
    "cookies_marketing": false
  },
  "recaptcha_token": "client-token",
  "timestamp": "2026-03-30T12:34:56.789Z",
  "user_agent": "Mozilla/5.0...",
  "referrer": "https://google.com",
  "campaignId": "homepage-hero"
}
```

`campaignId` is optional. When reCAPTCHA is enabled in config, the backend verifies the token against the configured `minimum_score`. Scores below the threshold are rejected.

Unsubscribe payload:

```json
{
  "email": "john@example.com",
  "reason": "not_interested",
  "timestamp": "2026-03-31T09:18:00.000Z",
  "user_agent": "Mozilla/5.0...",
  "referrer": "https://waitlist.example.com/account"
}
```

## Configuration

The `recaptcha` section supports:

- `enabled`
- `secret_key`
- `verify_url`
- `expected_action`
- `minimum_score`
- `timeout_seconds`

Default score guidance implemented by config:

- `0.9 - 1.0`: accept
- `0.7 - 0.8`: accept
- `0.5 - 0.6`: accept, but monitor
- `0.3 - 0.4`: reject or require additional verification
- `0.0 - 0.2`: reject

## HTML Preview

![waitlistFox HTML view](Waitingfox.png)

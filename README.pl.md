# waitlistFox

`waitlistFox` to prosty backend API do zapisu na waitlistę oparty o Echo i MySQL.

## Struktura

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

## Co jest gotowe

- publiczny endpoint `POST /v1/waitlist/subscribe`
- persystencja MySQL dla tabeli `waitlist_subscribers`
- ochrona przed duplikacją e-maila przez unikalny constraint
- pola zgód, ID kampanii, metadane requestu i `notified_at`
- opcjonalna weryfikacja reCAPTCHA konfigurowana w `config/config.json`
- `internal/config.go` z ładowaniem konfiguracji i połączeniem MySQL
- `internal/logger.go` z loggerem JSON do pliku i stdout
- `internal/validator.go`
- `internal/middleware/language.go`
- `internal/i18n/` z prostym loaderem tłumaczeń
- `static/index.html`

## Uruchomienie

```bash
cd waitlistfox
go run . -migrate=true
```

Domyślnie aplikacja czyta konfigurację z `config/config.json`.

## Endpointy

- `GET /health`
- `POST /v1/waitlist/subscribe`
- `POST /v1/waitlist/unsubscribe`

Przykładowy payload:

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

`campaignId` jest opcjonalne. Jeśli reCAPTCHA jest włączona w configu, backend weryfikuje token względem ustawionego `minimum_score`. Wyniki poniżej progu są odrzucane.

Payload wypisu:

```json
{
  "email": "john@example.com",
  "reason": "not_interested",
  "timestamp": "2026-03-31T09:18:00.000Z",
  "user_agent": "Mozilla/5.0...",
  "referrer": "https://waitlist.example.com/account"
}
```

## Konfiguracja

Sekcja `recaptcha` obsługuje:

- `enabled`
- `secret_key`
- `verify_url`
- `expected_action`
- `minimum_score`
- `timeout_seconds`

Domyślne zalecenia dla score, które możesz odwzorować w configu:

- `0.9 - 1.0`: akceptuj
- `0.7 - 0.8`: akceptuj
- `0.5 - 0.6`: akceptuj, ale monitoruj
- `0.3 - 0.4`: odrzuć albo wymagaj dodatkowej weryfikacji
- `0.0 - 0.2`: odrzuć

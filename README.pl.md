![waitlistFox logo](waitingfox-logo.png)

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

## Docker

Repo zawiera obraz Dockera i `compose.yaml` tylko dla API. Nie uruchamia MySQL. Bazę trzeba wystawić osobno.

1. Przygotuj serwerowy config jako `config/config.json`.
2. Zbuduj i uruchom API:

```bash
docker compose up -d --build
```

3. Uruchom migracje, kiedy są potrzebne:

```bash
docker compose run --rm api ./waitlistfox -config /app/config/config.json -migrate=true
```

Przydatne szczegóły runtime:

- API nasłuchuje w kontenerze na porcie `8081`
- port hosta kontroluje `WAITLISTFOX_PORT`
- config jest montowany z `./config` do `/app/config`
- logi są zapisywane do `./log`
- możesz też włączyć migracje przy starcie przez `MIGRATE_ON_START=true`

## Docker Produkcyjny

Dla wdrożenia na serwer bez GHCR użyj `compose.prod.yaml` i `deploy.sh`.

Pliki:

- `compose.prod.yaml`: produkcyjny Docker Compose tylko dla API
- `.env.production.example`: przykładowe wartości runtime dla portu, bind IP, timezone i nazwy kontenera
- `deploy.sh`: prosty helper do builda, migracji, startu, logów i statusu

Sugerowany setup na serwerze:

1. Sklonuj repo na serwer.
2. Utwórz właściwy config aplikacji:

```bash
cp config/config.example.json config/config.json
```

3. Utwórz produkcyjne wartości env:

```bash
cp .env.production.example .env.production
```

4. Edytuj oba pliki:
- `config/config.json`: host DB, user, hasło, nazwa bazy, reCAPTCHA
- `.env.production`: bind IP hosta i publiczny port

5. Nadaj prawa wykonywania skryptowi deploy:

```bash
chmod +x deploy.sh
```

6. Pierwsze wdrożenie:

```bash
./deploy.sh
```

7. Kolejne wdrożenia:

```bash
./deploy.sh
```

Przydatne komendy:

```bash
./deploy.sh status
./deploy.sh logs
./deploy.sh migrate
./deploy.sh deploy --skip-git-pull
./deploy.sh deploy --skip-migrate
```

Co robi `deploy.sh` w trybie `deploy`:

1. `git pull --ff-only`
2. `docker compose -f compose.prod.yaml build api`
3. opcjonalnie uruchamia migracje w jednorazowym kontenerze
4. `docker compose -f compose.prod.yaml up -d --remove-orphans api`

Zalecany wzorzec produkcyjny:

- trzymaj API zbindowane do `127.0.0.1`
- wystaw je publicznie przez Nginx albo Caddy jako reverse proxy
- trzymaj `config/config.json` tylko na serwerze i poza gitem

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

# Control Panel Architecture

## Overview
The monitoring dashboard keeps provider health, balances, and pricing in sync with each vendor's API so admins can triage outages and rebalance model routing automatically. It is intentionally lightweight—Express serves a handful of read-only endpoints while background workers periodically refresh telemetry into a shared in-memory cache.

## Data flow
1. **Provider definitions** (`src/config/providers.js`) enumerate every channel along with the URLs, required headers, and optional normalization helpers.
2. **Polling workers** (`src/services/monitor.js`) iterate over providers, calling each configured `health`, `balance`, and `models` endpoint. Responses are normalized, timestamped, and stored inside a `state` object.
3. **Scheduler** (`src/index.js`) runs `monitorRefresh()` on startup and again via a cron job (default every 5 minutes). The same job is also exposed to manual triggers if future automation layers need it.
4. **API layer** exports the latest cache so the UI or CLI scripts can draw tables, gauges, or notifications without requerying external services.

## Component responsibilities
- **Provider configs**: define endpoints, headers, and how to coerce responses into `{ amount, currency }` for balances or `{ price, modelId, modelName }` for pricing.
- **Monitor service**: encapsulates the HTTP logic with retries/fallbacks, aggregates errors, and builds the cheapest-model list sorted by price.
- **Server**: exposes JSON endpoints, handles cron scheduling, and gracefully starts after the first refresh for consistent data.

## Scheduling & resilience
- Health/balance requests are fire-and-forget but never block the whole refresh loop; failures are logged and mark the provider as `offline` until the next successful poll.
- Prices are refreshed alongside health and can be used for generating dynamic lists or alerting when high-cost providers become cheapest.
- The cron expression lives in `.env` to allow faster iteration during debugging (e.g., `*/1 * * * *` for every minute).

## Next steps
- Plug provider-specific API clients (using the APIs documented at each vendor) by filling mapper logic in `src/config/providers.js`.
- Layer a frontend (React/Next/Vite) or CLI (`curl`-friendly scripts) on top of the `/api/summary` endpoint to visualize uptime and spend.
- Hook into alerting systems (Slack, email) when a balance drops below a threshold or when cheapest provider changes.

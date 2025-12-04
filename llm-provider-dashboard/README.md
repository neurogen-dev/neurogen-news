# LLM Provider Monitoring Dashboard

This project is a lightweight control panel for monitoring LLM access providers that act as proxy OpenAI-compatible APIs. It continuously polls each channel to verify availability, track balances, and refresh model pricing so the admin team can prioritize the cheapest high-quality routes.

## Key features
- Scheduled health checks per provider with open/closed state tracking
- Balance refresh for each connected account (from their dashboards)
- Price registry that sorts every vendor/model pair by price for quick decisioning
- Simple REST API for an admin UI or automation to scrape current state

## Setup
1. Install dependencies: `npm install`
2. Copy `.env.example` to `.env` and fill provider API keys as needed
3. Run `npm start` to launch the dashboard server
4. The server refreshes data every `POLL_INTERVAL_MINUTES` (default 5)

## API surface
- `GET /api/providers/status` – returns last health check statuses
- `GET /api/providers/balances` – provides per-provider balance snapshots
- `GET /api/providers/prices` – sorted price list across all models
- `GET /api/summary` – bundle with uptime, balance, and cheapest-model insights

## Extending providers
Providers are defined in `src/config/providers.js`. For each provider you can override:
- `baseUrl`: the root endpoint
- `endpoints`: relative paths for `health`, `balance`, and `models`
- `headers`: extra headers per request
- `mappers`: simple functions to normalize provider responses for balances/models

Add new providers by supplying credentials (see `.env.example`) and adjusting the mapper logic so the polling service can derive numeric balances and price tags.

## Architecture Notes
The app spins up an Express server while `node-cron` and `monitorRefresh()` keep the data in memory. Each refresh stores the timestamp, handles per-provider failures gracefully, and writes a unified price table that the UI layer can sort from cheap to expensive. See `docs/architecture.md` for the full breakdown and workflow diagrams.

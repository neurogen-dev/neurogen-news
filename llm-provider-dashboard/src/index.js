import express from 'express';
import cron from 'node-cron';
import dotenv from 'dotenv';
import { refreshAllProviders, state } from './services/monitor.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT ? Number(process.env.PORT) : 4000;
const pollIntervalMinutes = Number(process.env.POLL_INTERVAL_MINUTES ?? 5);
const cronExpression = process.env.CRON_EXPRESSION ?? (pollIntervalMinutes > 59
  ? `0 */${Math.max(1, Math.floor(pollIntervalMinutes / 60))} * * *`
  : `*/${Math.max(1, pollIntervalMinutes)} * * * *`);

const safeJson = (payload) => ({ success: true, data: payload });

app.get('/api/providers/status', (_req, res) => {
  res.json(safeJson(state.statuses));
});

app.get('/api/providers/balances', (_req, res) => {
  res.json(safeJson(state.balances));
});

app.get('/api/providers/prices', (_req, res) => {
  res.json(safeJson(state.priceTable));
});

app.get('/api/summary', (_req, res) => {
  res.json(safeJson({
    statuses: state.statuses,
    balances: state.balances,
    priceTable: state.priceTable,
    summary: state.summary
  }));
});

const startScheduler = () => {
  if (!cronExpression) {
    return Promise.resolve();
  }

  cron.schedule(cronExpression, () => {
    console.log(`Running refresh job (${cronExpression})`);
    refreshAllProviders();
  });
  return Promise.resolve();
};

const bootstrap = async () => {
  console.log('Initializing provider dashboard...');
  await refreshAllProviders();
  await startScheduler();
  app.listen(PORT, () => {
    console.log(`Dashboard listening on http://localhost:${PORT}`);
    console.log(`Polling every ${pollIntervalMinutes} minute(s) via '${cronExpression}'`);
  });
};

bootstrap().catch((error) => {
  console.error('Failed to start dashboard', error);
  process.exit(1);
});

import axios from 'axios';
import dotenv from 'dotenv';
import { providers } from '../config/providers.js';

dotenv.config();

const TIMEOUT_MS = 10000;

const buildHeaders = (provider) => {
  const authKey = provider.apiKeyEnv ? process.env[provider.apiKeyEnv] : undefined;
  const headers = typeof provider.headers === 'function'
    ? provider.headers()
    : { ...(provider.headers ?? {}) };

  if (authKey) {
    headers[provider.authHeader ?? 'Authorization'] = authKey;
  }

  return headers;
};

const callEndpoint = async (baseUrl, path, options = {}) => {
  if (!path) return { skipped: true, reason: 'no endpoint configured' };

  const url = `${baseUrl}${path}`;
  try {
    const resp = await axios.get(url, {
      timeout: TIMEOUT_MS,
      ...options
    });
    return { ok: true, data: resp.data ?? resp };
  } catch (error) {
    return {
      ok: false,
      error: error.response?.status
        ? `${error.response.status} ${error.response.statusText}`
        : `${error.code ?? 'ERR'}: ${error.message}`
    };
  }
};

const testProvider = async (provider) => {
  console.log(`\n=== Provider: ${provider.name} (${provider.id}) ===`);

  // health / balance / models (как в основном мониторе)
  const baseHeaders = buildHeaders(provider);

  const [health, balance, models] = await Promise.all([
    callEndpoint(provider.baseUrl, provider.endpoints?.health, { headers: baseHeaders }),
    callEndpoint(provider.baseUrl, provider.endpoints?.balance, { headers: baseHeaders }),
    callEndpoint(provider.baseUrl, provider.endpoints?.models, { headers: baseHeaders })
  ]);

  console.log('health:', health.ok ? 'OK' : 'FAIL', health.ok ? '' : `(${health.error})`);
  console.log('balance:', balance.ok ? 'OK' : 'FAIL', balance.ok ? '' : `(${balance.error})`);
  console.log('models:', models.ok ? 'OK' : 'FAIL', models.ok ? '' : `(${models.error})`);

  // Пытаемся вызвать стандартный New API /api/models с системным токеном панели
  const panelTokenEnv = provider.panelAccessTokenEnv;
  const panelToken = panelTokenEnv ? process.env[panelTokenEnv] : null;

  if (!panelToken) {
    console.log('new-api /api/models: SKIPPED (no panel token in env)');
    return;
  }

  const newApiModels = await callEndpoint(provider.baseUrl, '/api/models', {
    headers: {
      Authorization: `Bearer ${panelToken}`
    }
  });

  if (newApiModels.ok) {
    const payload = newApiModels.data;
    const keys = payload?.data ? Object.keys(payload.data) : [];
    console.log(`new-api /api/models: OK (channels: ${keys.length})`);
  } else if (newApiModels.skipped) {
    console.log(`new-api /api/models: SKIPPED (${newApiModels.reason})`);
  } else {
    console.log(`new-api /api/models: FAIL (${newApiModels.error})`);
  }
};

const main = async () => {
  console.log('Running provider test script with env keys...\n');
  for (const provider of providers) {
    // eslint-disable-next-line no-await-in-loop
    await testProvider(provider);
  }
  console.log('\nDone.');
};

main().catch((err) => {
  console.error('Test script failed:', err);
  process.exit(1);
});



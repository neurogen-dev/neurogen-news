import axios from 'axios';
import dayjs from 'dayjs';
import { providers } from '../config/providers.js';

const TIMEOUT_MS = 10000;
const DEFAULT_CURRENCY = 'USD';

export const state = {
  statuses: {},
  balances: {},
  priceTable: [],
  summary: {
    lastRefreshed: null,
    cheapestModel: null,
    totalProviders: providers.length
  }
};

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

const callEndpoint = async (provider, key) => {
  const endpoint = provider.endpoints?.[key];
  if (!endpoint) {
    return { error: 'endpoint not configured' };
  }

  try {
    const url = `${provider.baseUrl}${endpoint}`;
    const response = await axios.get(url, {
      headers: buildHeaders(provider),
      timeout: TIMEOUT_MS
    });
    return response.data ?? response;
  } catch (error) {
    return { error: `${error.code ?? 'error'}: ${error.message}` };
  }
};

// Вспомогательный вызов стандартного New API фронтового эндпоинта /api/models.
// Документация: https://docs.newapi.pro/api/get-available-models-list/
// Мы используем его только при наличии panelAccessTokenEnv у провайдера.
const callNewApiModels = async (provider) => {
  if (!provider.panelAccessTokenEnv) {
    return null;
  }

  const token = process.env[provider.panelAccessTokenEnv];
  if (!token) {
    return null;
  }

  try {
    const url = `${provider.baseUrl}/api/models`;
    const response = await axios.get(url, {
      headers: {
        Authorization: `Bearer ${token}`
        // При необходимости сюда можно добавить заголовок New-Api-User,
        // если ты захочешь жёстко указать user_id.
      },
      timeout: TIMEOUT_MS
    });
    return response.data ?? response;
  } catch (error) {
    return { error: `${error.code ?? 'error'}: ${error.message}` };
  }
};

const normaliseBalance = (provider, payload) => {
  const mapper = provider.mappers?.balance;
  const raw = mapper ? mapper(payload) : payload;
  const rawAmount = raw?.amount ?? raw?.total ?? raw?.balance ?? null;
  const numericAmount = rawAmount == null ? null : Number(rawAmount);

  return {
    amount: Number.isNaN(numericAmount) ? null : numericAmount,
    currency: raw?.currency ?? DEFAULT_CURRENCY
  };
};

const normaliseModels = (provider, payload) => {
  const mapper = provider.mappers?.models;
  const raw = mapper ? mapper(payload) : payload;
  if (!Array.isArray(raw)) {
    return [];
  }

  return raw
    .map((model) => ({
      providerId: provider.id,
      providerName: provider.name,
      modelId: model.modelId ?? model.id ?? model.name,
      modelName: model.modelName ?? model.name ?? model.model ?? 'unknown',
      price: typeof model.price === 'number'
        ? model.price
        : Number(model.price ?? model.cost ?? null),
      currency: model.currency ?? DEFAULT_CURRENCY
    }))
    .filter((item) => !Number.isNaN(item.price));
};

export const refreshAllProviders = async () => {
  const timestamp = dayjs().toISOString();
  const priceAccumulator = [];

  const statuses = {};
  const balances = {};

  await Promise.all(providers.map(async (provider) => {
    const [
      healthData,
      balanceData,
      modelsData,
      newApiModelsData
    ] = await Promise.all([
      callEndpoint(provider, 'health'),
      callEndpoint(provider, 'balance'),
      callEndpoint(provider, 'models'),
      callNewApiModels(provider)
    ]);

    const online = !healthData?.error;
    const healthDetail = online ? 'reachable' : (healthData?.error ?? 'timeout');

    statuses[provider.id] = {
      providerName: provider.name,
      online,
      detail: healthDetail,
      lastChecked: timestamp
    };

    const balance = normaliseBalance(provider, balanceData);
    balances[provider.id] = {
      providerName: provider.name,
      amount: Number.isNaN(balance.amount) ? null : balance.amount,
      currency: balance.currency,
      lastChecked: timestamp
    };

    let models = normaliseModels(provider, modelsData);

    // Если доступен New API /api/models и он вернул success=true,
    // то оставляем только те модели, которые реально доступны текущему пользователю.
    if (newApiModelsData && !newApiModelsData.error && newApiModelsData.success && newApiModelsData.data) {
      const allowedNames = new Set(
        Object.values(newApiModelsData.data)
          .flat()
          .map((name) => String(name))
      );

      models = models.filter((model) => (
        allowedNames.has(model.modelId) || allowedNames.has(model.modelName)
      ));
    }

    priceAccumulator.push(...models);
  }));

  state.statuses = statuses;
  state.balances = balances;
  state.priceTable = priceAccumulator
    .sort((a, b) => a.price - b.price);
  state.summary = {
    lastRefreshed: timestamp,
    totalProviders: providers.length,
    cheapestModel: state.priceTable[0] ?? null
  };
};

export const getSummary = () => ({ ...state });

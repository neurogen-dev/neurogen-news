const defaultBalanceMapper = (resp) => {
  const payload = resp?.balance ?? resp?.data?.balance ?? resp?.data ?? resp;
  return {
    amount: typeof payload === 'number' ? payload : payload?.amount ?? null,
    currency: payload?.currency ?? 'USD'
  };
};

const defaultModelsMapper = (resp) => {
  const items = resp?.models ?? resp?.data?.models ?? resp?.data ?? [];
  return (Array.isArray(items) ? items : []).map((model) => {
    const modelId = model.id ?? model.name ?? 'unnamed';
    return {
      modelId,
      modelName: model.name ?? model.model ?? modelId,
      price: model.price ?? model.cost ?? null,
      currency: model.currency ?? 'USD'
    };
  });
};

export const providers = [
  {
    id: 'yunwu',
    name: 'Yunwu AI',
    baseUrl: 'https://yunwu.ai',
    apiKeyEnv: 'YUNWU_API_KEY',
    authHeader: 'Authorization',
    // Дополнительно поддерживаем стандартный New API фронтовый эндпоинт /api/models.
    // Для него нужен отдельный AccessToken пользователя (см. New API docs),
    // который, при желании, можно прописать в переменную окружения YUNWU_PANEL_TOKEN.
    panelAccessTokenEnv: 'YUNWU_PANEL_TOKEN',
    endpoints: {
      health: '/health',
      balance: '/api/balance',
      models: '/api/models'
    },
    mappers: {
      balance: defaultBalanceMapper,
      models: defaultModelsMapper
    }
  },
  {
    id: 'neko',
    name: 'Neko API',
    baseUrl: 'https://nekoapi.com',
    apiKeyEnv: 'NEKO_API_KEY',
    authHeader: 'X-API-KEY',
    panelAccessTokenEnv: 'NEKO_PANEL_TOKEN',
    endpoints: {
      health: '/v1/ping',
      balance: '/v1/account',
      models: '/v1/models'
    },
    mappers: {
      balance: (resp) => ({
        amount: resp?.data?.balance ?? resp?.balance ?? null,
        currency: resp?.data?.currency ?? 'USD'
      }),
      models: (resp) => (
        (resp?.data?.models ?? resp?.models ?? [])
          .map((model) => ({
            modelId: model.id ?? model.name,
            modelName: model.name ?? model.id,
            price: model.price_per_call ?? model.price ?? null,
            currency: model.currency ?? 'USD'
          }))
      )
    }
  },
  {
    id: 'aifuture',
    name: 'AI Future',
    baseUrl: 'https://api.aifuture.pw',
    apiKeyEnv: 'AIFUTURE_API_KEY',
    authHeader: 'Authorization',
    panelAccessTokenEnv: 'AIFUTURE_PANEL_TOKEN',
    endpoints: {
      health: '/v1/health',
      balance: '/v1/account',
      models: '/v1/models'
    },
    mappers: {
      balance: defaultBalanceMapper,
      models: defaultModelsMapper
    }
  },
  {
    id: 'mnapi',
    name: 'MNAPI',
    // According to MN API client docs, the OpenAI‑совместимый базовый URL:
    // https://api.mnapi.com или https://cf.mnapi.com (для зарубежных клиентов).
    // Для опроса служебных эндпоинтов мониторинга используем основной API‑хост.
    baseUrl: 'https://api.mnapi.com',
    apiKeyEnv: 'MNAPI_API_KEY',
    authHeader: 'Authorization',
    panelAccessTokenEnv: 'MNAPI_PANEL_TOKEN',
    endpoints: {
      health: '/status',
      balance: '/account/balance',
      models: '/models'
    },
    mappers: {
      balance: defaultBalanceMapper,
      models: defaultModelsMapper
    }
  }
];

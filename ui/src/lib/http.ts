const BASE_URL = import.meta.env.VITE_ADMIN_API_URL || 'http://localhost:9901';

export class HttpError extends Error {
  public readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.name = 'HttpError';
  }
}

async function request<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = endpoint.startsWith('http') ? endpoint : `${BASE_URL}${endpoint}`;
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const errorText = await response.text().catch(() => 'Unknown error');
    throw new HttpError(response.status, `API request failed (${response.status}): ${errorText}`);
  }

  return response.json() as Promise<T>;
}

export const http = {
  get: <T>(endpoint: string, options?: RequestInit): Promise<T> =>
    request<T>(endpoint, { ...options, method: 'GET' }),

  post: <T, B = unknown>(endpoint: string, body?: B, options?: RequestInit): Promise<T> =>
    request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),
};

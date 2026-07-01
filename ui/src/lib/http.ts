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
    let errorMsg = `API request failed (${response.status})`;
    try {
      const errorText = await response.text();
      try {
        const errorJson = JSON.parse(errorText);
        if (errorJson.error) {
          errorMsg = errorJson.error;
        } else if (errorJson.message) {
          errorMsg = errorJson.message;
        } else if (errorText) {
          errorMsg = errorText;
        }
      } catch {
        if (errorText) errorMsg = errorText;
      }
    } catch {
      // ignore
    }
    throw new HttpError(response.status, errorMsg);
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

  postRaw: <T>(
    endpoint: string,
    body: string,
    contentType = 'application/x-yaml',
    options?: RequestInit
  ): Promise<T> =>
    request<T>(endpoint, {
      ...options,
      method: 'POST',
      headers: {
        'Content-Type': contentType,
        ...options?.headers,
      },
      body,
    }),
};

// apps/web/src/api/client.ts
import { handleMockApi } from './mockService';

export const API_BASE = (import.meta.env.VITE_API_URL || '') + '/api/v1';

export interface ApiResponse<T = any> {
  data: T;
  meta: {
    request_id: string;
    server_time: string;
  };
}

export class ApiError extends Error {
  code: string;
  details?: any[];
  requestId?: string;

  constructor(code: string, message: string, details?: any[], requestId?: string) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.details = details;
    this.requestId = requestId;
  }
}

export async function apiFetch<T = any>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('bn_token');
  
  // If we are operating in mock mode or using a mock token, resolve with mock service immediately
  const isMockToken = token && token.startsWith('bn-mock');
  const hasNoBackend = !import.meta.env.VITE_API_URL;

  // For auth endpoints in static deployment without dedicated backend, serve mock instantly
  if (hasNoBackend && (endpoint.startsWith('/auth/login') || endpoint.startsWith('/auth/register'))) {
    return handleMockApi<T>(endpoint, options);
  }

  if (isMockToken) {
    return handleMockApi<T>(endpoint, options);
  }

  const headers = new Headers(options.headers || {});
  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  try {
    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    // If server responded with 404/405 (e.g. static server hosting SPA), fall back to mock
    if ((response.status === 404 || response.status === 405) && hasNoBackend) {
      console.warn(`[BN ONE API] Endpoint ${endpoint} returned ${response.status}, falling back to mock service.`);
      return handleMockApi<T>(endpoint, options);
    }

    const contentType = response.headers.get('content-type') || '';
    if (!contentType.includes('application/json')) {
      if (hasNoBackend) {
        return handleMockApi<T>(endpoint, options);
      }
      throw new ApiError('INVALID_RESPONSE', 'Server returned non-JSON response');
    }

    const json = await response.json();

    if (!response.ok) {
      const err = json.error || { code: 'UNKNOWN_ERROR', message: 'An unexpected error occurred' };
      throw new ApiError(err.code, err.message, err.details, err.request_id);
    }

    return json.data as T;
  } catch (err: any) {
    // If network failed or unexpected parse error, fall back to mock service when no backend configured
    if (hasNoBackend) {
      console.warn(`[BN ONE API] Network error for ${endpoint}, serving mock data.`, err);
      return handleMockApi<T>(endpoint, options);
    }
    throw err;
  }
}

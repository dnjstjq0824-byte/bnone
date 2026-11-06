// apps/web/src/api/client.ts
import { handleMockApi } from './mockService';

export const API_BASE = (import.meta.env?.VITE_API_URL || '') + '/api/v1';

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
  const token = localStorage.getItem('moa_token');

  // Mock demo data is served ONLY when explicitly enabled at build time
  // (VITE_MOCK_ENABLED=true). Production builds fail closed: without a
  // configured backend the call errors instead of masquerading demo
  // fixtures as real user data.
  const mockAllowed = import.meta.env?.VITE_MOCK_ENABLED === 'true';
  const hasNoBackend = !import.meta.env?.VITE_API_URL;
  const isMockToken = token && token.startsWith('moa-mock');

  if (mockAllowed && (isMockToken || hasNoBackend)) {
    return handleMockApi<T>(endpoint, options);
  }

  if (isMockToken && !mockAllowed) {
    localStorage.removeItem('moa_token');
    throw new ApiError('MOCK_DISABLED', '데모 세션은 이 환경에서 사용할 수 없습니다. 다시 로그인하세요.');
  }

  if (hasNoBackend) {
    throw new ApiError('BACKEND_NOT_CONFIGURED', '백엔드 API 주소(VITE_API_URL)가 설정되지 않았습니다. 관리자에게 문의하세요.');
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
    // only when demo mode is explicitly enabled.
    if ((response.status === 404 || response.status === 405) && mockAllowed) {
      console.warn(`[MOA API] Endpoint ${endpoint} returned ${response.status}, falling back to mock service.`);
      return handleMockApi<T>(endpoint, options);
    }

    const contentType = response.headers.get('content-type') || '';
    if (!contentType.includes('application/json')) {
      if (mockAllowed) {
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
    // If network failed or unexpected parse error, fall back to mock service
    // only when demo mode is explicitly enabled.
    if (mockAllowed && err?.code !== 'BACKEND_NOT_CONFIGURED' && err?.code !== 'MOCK_DISABLED') {
      console.warn(`[MOA API] Network error for ${endpoint}, serving mock data.`, err);
      return handleMockApi<T>(endpoint, options);
    }
    throw err;
  }
}

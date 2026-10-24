export const API_BASE = '/api/v1';

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
  const headers = new Headers(options.headers || {});
  
  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  const json = await response.json();

  if (!response.ok) {
    const err = json.error || { code: 'UNKNOWN_ERROR', message: 'An unexpected error occurred' };
    throw new ApiError(err.code, err.message, err.details, err.request_id);
  }

  return json.data as T;
}


export const API_BASE = (import.meta.env.VITE_API_URL || '') + '/api/v1';

export async function adminFetch<T = any>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('bn_admin_token');
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
    const err = json.error || { code: 'UNKNOWN_ERROR', message: 'An operational error occurred' };
    throw new Error(err.message || 'Operation failed');
  }

  return json.data as T;
}


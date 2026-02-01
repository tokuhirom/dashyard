import type { Dashboard, DashboardsResponse, LabelValuesResponse, PrometheusResponse } from '../types';

export interface OAuthProviderInfo {
  name: string;
  url: string;
}

export interface AuthInfo {
  password_enabled: boolean;
  oauth_providers: OAuthProviderInfo[];
}

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const resp = await fetch(url, options);
  if (resp.status === 401) {
    throw new ApiError(401, 'Unauthorized');
  }
  if (!resp.ok) {
    const body = await resp.text();
    throw new ApiError(resp.status, body);
  }
  return resp.json();
}

export async function fetchAuthInfo(): Promise<AuthInfo> {
  return request('/api/auth-info');
}

export async function login(userId: string, password: string): Promise<{ user_id: string }> {
  return request('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: userId, password }),
  });
}

export async function fetchDashboards(): Promise<DashboardsResponse> {
  return request('/api/dashboards');
}

export async function fetchDashboard(path: string): Promise<Dashboard> {
  return request(`/api/dashboards/${path}`);
}

export async function queryPrometheus(
  query: string,
  start: number,
  end: number,
  step: string,
): Promise<PrometheusResponse> {
  const params = new URLSearchParams({
    query,
    start: start.toString(),
    end: end.toString(),
    step,
  });
  return request(`/api/query?${params}`);
}

export async function fetchLabelValues(label: string, match?: string): Promise<LabelValuesResponse> {
  const params = new URLSearchParams({ label });
  if (match) {
    params.set('match', match);
  }
  return request(`/api/label-values?${params}`);
}

export async function fetchDashboardSource(path: string): Promise<string> {
  const resp = await fetch(`/api/dashboard-source/${path}`);
  if (resp.status === 401) {
    throw new ApiError(401, 'Unauthorized');
  }
  if (!resp.ok) {
    const body = await resp.text();
    throw new ApiError(resp.status, body);
  }
  return resp.text();
}

export { ApiError };

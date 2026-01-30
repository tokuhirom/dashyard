import { useState, useEffect, useCallback } from 'react';
import { fetchDashboards, fetchDashboard, ApiError } from '../api/client';
import type { Dashboard, DashboardsResponse } from '../types';

interface UseDashboardsResult {
  dashboardsData: DashboardsResponse | null;
  loading: boolean;
  error: string | null;
  reload: () => void;
}

export function useDashboards(onAuthError: () => void): UseDashboardsResult {
  const [dashboardsData, setDashboardsData] = useState<DashboardsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError(null);
    fetchDashboards()
      .then((data) => {
        setDashboardsData(data);
        setLoading(false);
      })
      .catch((err) => {
        if (err instanceof ApiError && err.status === 401) {
          onAuthError();
        } else {
          setError(err.message || 'Failed to load dashboards');
        }
        setLoading(false);
      });
  }, [onAuthError]);

  useEffect(() => {
    load();
  }, [load]);

  return { dashboardsData, loading, error, reload: load };
}

interface UseDashboardDetailResult {
  dashboard: Dashboard | null;
  loading: boolean;
  error: string | null;
}

export function useDashboardDetail(path: string | undefined, onAuthError: () => void): UseDashboardDetailResult {
  const [dashboard, setDashboard] = useState<Dashboard | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!path) return;

    let cancelled = false;
    setLoading(true);
    setError(null);

    fetchDashboard(path)
      .then((data) => {
        if (!cancelled) {
          setDashboard(data);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          if (err instanceof ApiError && err.status === 401) {
            onAuthError();
          } else {
            setError(err.message || 'Failed to load dashboard');
          }
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [path, onAuthError]);

  return { dashboard, loading, error };
}

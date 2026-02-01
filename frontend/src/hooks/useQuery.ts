import { useState, useEffect } from 'react';
import { queryPrometheus, ApiError } from '../api/client';
import type { PrometheusResponse, TimeRange } from '../types';
import { getTimeRangeParams } from '../utils/time';

interface UseQueryResult {
  data: PrometheusResponse | null;
  loading: boolean;
  error: string | null;
}

export function useQuery(query: string | undefined, timeRange: TimeRange, datasource?: string): UseQueryResult {
  const [data, setData] = useState<PrometheusResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!query) return;

    let cancelled = false;
    setLoading(true);
    setError(null);

    const { start, end, step } = getTimeRangeParams(timeRange);

    queryPrometheus(query, start, end, step, datasource)
      .then((result) => {
        if (!cancelled) {
          setData(result);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          if (err instanceof ApiError && err.status === 401) {
            setError('Session expired');
          } else {
            setError(err.message || 'Query failed');
          }
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [query, timeRange, datasource]);

  return { data, loading, error };
}

import { useState, useEffect } from 'react';
import { queryDatasource, ApiError } from '../api/client';
import type { QueryResponse, TimeRange } from '../types';
import { getTimeRangeParams } from '../utils/time';

interface UseQueryResult {
  data: QueryResponse | null;
  loading: boolean;
  error: string | null;
}

export function useQuery(query: string | undefined, timeRange: TimeRange, datasource?: string): UseQueryResult {
  const [data, setData] = useState<QueryResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!query) return;

    let cancelled = false;
    setLoading(true);
    setError(null);

    const { start, end, step } = getTimeRangeParams(timeRange);

    queryDatasource(query, start, end, step, datasource)
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

import { useState, useCallback, useEffect } from 'react';
import { LoginForm } from './components/LoginForm';
import { Layout } from './components/Layout';
import { DashboardView } from './components/DashboardView';
import { useDashboards } from './hooks/useDashboards';
import { DEFAULT_TIME_RANGE, TIME_RANGES, computeStep } from './utils/time';
import type { TimeRange } from './types';

function parseDashboardPath(): string | null {
  const path = window.location.pathname;
  if (path.startsWith('/d/')) {
    return path.slice(3);
  }
  return null;
}

function parseTimeRange(): TimeRange {
  const params = new URLSearchParams(window.location.search);

  // Check for absolute range (ISO 8601 from/to)
  const from = params.get('from');
  const to = params.get('to');
  if (from && to) {
    const startUnix = Math.floor(new Date(from).getTime() / 1000);
    const endUnix = Math.floor(new Date(to).getTime() / 1000);
    if (!isNaN(startUnix) && !isNaN(endUnix) && startUnix < endUnix) {
      const duration = endUnix - startUnix;
      const step = computeStep(duration);
      const fromDate = new Date(startUnix * 1000);
      const toDate = new Date(endUnix * 1000);
      const fmt = (d: Date) => `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
      return {
        type: 'absolute',
        label: `${fmt(fromDate)} – ${fmt(toDate)}`,
        start: startUnix,
        end: endUnix,
        step,
      };
    }
  }

  // Relative range
  const t = params.get('t');
  if (t) {
    const found = TIME_RANGES.find((r) => r.value === t);
    if (found) return found;
  }
  return DEFAULT_TIME_RANGE;
}

function parseVariableValues(): Record<string, string> {
  const params = new URLSearchParams(window.location.search);
  const values: Record<string, string> = {};
  params.forEach((value, key) => {
    if (key.startsWith('var-')) {
      values[key.slice(4)] = value;
    }
  });
  return values;
}

function buildUrl(dashboardPath: string, timeRange: TimeRange, varValues?: Record<string, string>): string {
  let url = `/d/${dashboardPath}`;
  const params = new URLSearchParams();
  if (timeRange.type === 'absolute') {
    const fromISO = new Date(timeRange.start * 1000).toISOString();
    const toISO = new Date(timeRange.end * 1000).toISOString();
    params.set('from', fromISO);
    params.set('to', toISO);
  } else if (timeRange.value !== DEFAULT_TIME_RANGE.value) {
    params.set('t', timeRange.value);
  }
  if (varValues) {
    for (const [name, value] of Object.entries(varValues)) {
      params.set(`var-${name}`, value);
    }
  }
  const qs = params.toString();
  return qs ? `${url}?${qs}` : url;
}

function App() {
  const [authenticated, setAuthenticated] = useState(true); // Optimistic; API calls will detect 401
  const [currentPath, setCurrentPath] = useState<string | null>(parseDashboardPath);
  const [timeRange, setTimeRange] = useState<TimeRange>(parseTimeRange);
  const [variableValues, setVariableValues] = useState<Record<string, string>>(parseVariableValues);
  const [refreshInterval, setRefreshInterval] = useState(0);

  useEffect(() => {
    if (refreshInterval <= 0) return;
    const id = setInterval(() => {
      setTimeRange((prev) => (prev.type === 'relative' ? { ...prev } : prev));
    }, refreshInterval);
    return () => clearInterval(id);
  }, [refreshInterval]);

  const handleAuthError = useCallback(() => {
    setAuthenticated(false);
  }, []);

  const handleLoginSuccess = useCallback(() => {
    setAuthenticated(true);
  }, []);

  const { dashboardsData, loading, error } = useDashboards(handleAuthError);

  const onNavigate = useCallback((path: string) => {
    setCurrentPath(path);
    setVariableValues({});
    setTimeRange((prev) => {
      const url = buildUrl(path, prev);
      window.history.pushState(null, '', url);
      return prev;
    });
  }, []);

  const onVariableValuesChange = useCallback((values: Record<string, string>) => {
    setVariableValues(values);
    setCurrentPath((prev) => {
      if (prev) {
        setTimeRange((tr) => {
          const url = buildUrl(prev, tr, values);
          window.history.replaceState(null, '', url);
          return tr;
        });
      }
      return prev;
    });
  }, []);

  const onTimeRangeChange = useCallback((range: TimeRange) => {
    setTimeRange(range);
    setCurrentPath((prev) => {
      if (prev) {
        setVariableValues((vars) => {
          const url = buildUrl(prev, range, vars);
          window.history.replaceState(null, '', url);
          return vars;
        });
      }
      return prev;
    });
  }, []);

  useEffect(() => {
    const handlePopState = () => {
      setCurrentPath(parseDashboardPath());
      setTimeRange(parseTimeRange());
      setVariableValues(parseVariableValues());
    };
    window.addEventListener('popstate', handlePopState);
    return () => window.removeEventListener('popstate', handlePopState);
  }, []);

  if (!authenticated) {
    return <LoginForm onLoginSuccess={handleLoginSuccess} />;
  }

  if (loading) {
    return <div className="app-loading">Loading...</div>;
  }

  if (error) {
    return <div className="app-error">Error: {error}</div>;
  }

  if (!dashboardsData) {
    return <div className="app-loading">Loading...</div>;
  }

  // Default to the first dashboard if none selected
  const activePath = currentPath || dashboardsData.dashboards[0]?.path;

  // Redirect root to first dashboard
  if (!currentPath && activePath && window.location.pathname === '/') {
    window.history.replaceState(null, '', buildUrl(activePath, timeRange, variableValues));
  }

  if (!activePath) {
    return <div className="app-empty">No dashboards configured</div>;
  }

  return (
    <Layout
      tree={dashboardsData.tree}
      currentPath={activePath}
      timeRange={timeRange}
      onTimeRangeChange={onTimeRangeChange}
      onNavigate={onNavigate}
      siteTitle={dashboardsData.site_title}
      headerColor={dashboardsData.header_color}
      refreshInterval={refreshInterval}
      onRefreshIntervalChange={setRefreshInterval}
    >
      <DashboardView
        path={activePath}
        timeRange={timeRange}
        onAuthError={handleAuthError}
        variableValues={variableValues}
        onVariableValuesChange={onVariableValuesChange}
      />
    </Layout>
  );
}

export default App;

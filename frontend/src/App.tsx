import { useState, useCallback, useEffect } from 'react';
import { LoginForm } from './components/LoginForm';
import { Layout } from './components/Layout';
import { DashboardView } from './components/DashboardView';
import { useDashboards } from './hooks/useDashboards';
import { DEFAULT_TIME_RANGE, TIME_RANGES } from './utils/time';
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
  const t = params.get('t');
  if (t) {
    const found = TIME_RANGES.find((r) => r.value === t);
    if (found) return found;
  }
  return DEFAULT_TIME_RANGE;
}

function buildUrl(dashboardPath: string, timeRange: TimeRange): string {
  let url = `/d/${dashboardPath}`;
  if (timeRange.value !== DEFAULT_TIME_RANGE.value) {
    url += `?t=${timeRange.value}`;
  }
  return url;
}

function App() {
  const [authenticated, setAuthenticated] = useState(true); // Optimistic; API calls will detect 401
  const [currentPath, setCurrentPath] = useState<string | null>(parseDashboardPath);
  const [timeRange, setTimeRange] = useState<TimeRange>(parseTimeRange);

  const handleAuthError = useCallback(() => {
    setAuthenticated(false);
  }, []);

  const handleLoginSuccess = useCallback(() => {
    setAuthenticated(true);
  }, []);

  const { dashboardsData, loading, error } = useDashboards(handleAuthError);

  const onNavigate = useCallback((path: string) => {
    setCurrentPath(path);
    setTimeRange((prev) => {
      const url = buildUrl(path, prev);
      window.history.pushState(null, '', url);
      return prev;
    });
  }, []);

  const onTimeRangeChange = useCallback((range: TimeRange) => {
    setTimeRange(range);
    setCurrentPath((prev) => {
      if (prev) {
        const url = buildUrl(prev, range);
        window.history.replaceState(null, '', url);
      }
      return prev;
    });
  }, []);

  useEffect(() => {
    const handlePopState = () => {
      setCurrentPath(parseDashboardPath());
      setTimeRange(parseTimeRange());
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
    window.history.replaceState(null, '', buildUrl(activePath, timeRange));
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
    >
      <DashboardView
        path={activePath}
        timeRange={timeRange}
        onAuthError={handleAuthError}
      />
    </Layout>
  );
}

export default App;

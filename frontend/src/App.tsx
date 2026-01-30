import { useState, useCallback } from 'react';
import { LoginForm } from './components/LoginForm';
import { Layout } from './components/Layout';
import { DashboardView } from './components/DashboardView';
import { useDashboards } from './hooks/useDashboards';
import { DEFAULT_TIME_RANGE } from './utils/time';
import type { TimeRange } from './types';

function App() {
  const [authenticated, setAuthenticated] = useState(true); // Optimistic; API calls will detect 401
  const [currentPath, setCurrentPath] = useState<string | null>(null);
  const [timeRange, setTimeRange] = useState<TimeRange>(DEFAULT_TIME_RANGE);

  const handleAuthError = useCallback(() => {
    setAuthenticated(false);
  }, []);

  const handleLoginSuccess = useCallback(() => {
    setAuthenticated(true);
  }, []);

  const { dashboardsData, loading, error } = useDashboards(handleAuthError);

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

  if (!activePath) {
    return <div className="app-empty">No dashboards configured</div>;
  }

  return (
    <Layout
      tree={dashboardsData.tree}
      currentPath={activePath}
      timeRange={timeRange}
      onTimeRangeChange={setTimeRange}
      onNavigate={setCurrentPath}
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

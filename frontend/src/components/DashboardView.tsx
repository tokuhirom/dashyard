import { useState, useEffect } from 'react';
import type { TimeRange } from '../types';
import { useDashboardDetail } from '../hooks/useDashboards';
import { fetchDashboardSource, ApiError } from '../api/client';
import { RowView } from './RowView';

interface DashboardViewProps {
  path: string;
  timeRange: TimeRange;
  onAuthError: () => void;
}

export function DashboardView({ path, timeRange, onAuthError }: DashboardViewProps) {
  const { dashboard, loading, error } = useDashboardDetail(path, onAuthError);
  const [showSource, setShowSource] = useState(false);
  const [source, setSource] = useState<string | null>(null);
  const [sourceLoading, setSourceLoading] = useState(false);

  useEffect(() => {
    setShowSource(false);
    setSource(null);
  }, [path]);

  useEffect(() => {
    if (!showSource || source !== null) return;
    setSourceLoading(true);
    fetchDashboardSource(path)
      .then(setSource)
      .catch((err) => {
        if (err instanceof ApiError && err.status === 401) {
          onAuthError();
        }
        setSource(`Error loading source: ${err.message}`);
      })
      .finally(() => setSourceLoading(false));
  }, [showSource, source, path, onAuthError]);

  if (loading) {
    return <div className="dashboard-loading">Loading dashboard...</div>;
  }

  if (error) {
    return <div className="dashboard-error">Error: {error}</div>;
  }

  if (!dashboard) {
    return <div className="dashboard-empty">Dashboard not found</div>;
  }

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <h1 className="dashboard-title">{dashboard.title}</h1>
        <button
          className={`dashboard-source-btn${showSource ? ' active' : ''}`}
          onClick={() => setShowSource(!showSource)}
        >
          {showSource ? 'Dashboard' : 'Source'}
        </button>
      </div>
      {showSource ? (
        sourceLoading ? (
          <div className="dashboard-loading">Loading source...</div>
        ) : (
          <pre className="dashboard-source">{source}</pre>
        )
      ) : (
        dashboard.rows.map((row, idx) => (
          <RowView key={idx} row={row} rowIndex={idx} timeRange={timeRange} />
        ))
      )}
    </div>
  );
}

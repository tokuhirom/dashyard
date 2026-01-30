import type { TimeRange } from '../types';
import { useDashboardDetail } from '../hooks/useDashboards';
import { RowView } from './RowView';

interface DashboardViewProps {
  path: string;
  timeRange: TimeRange;
  onAuthError: () => void;
}

export function DashboardView({ path, timeRange, onAuthError }: DashboardViewProps) {
  const { dashboard, loading, error } = useDashboardDetail(path, onAuthError);

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
      <h1 className="dashboard-title">{dashboard.title}</h1>
      {dashboard.rows.map((row, idx) => (
        <RowView key={idx} row={row} timeRange={timeRange} />
      ))}
    </div>
  );
}

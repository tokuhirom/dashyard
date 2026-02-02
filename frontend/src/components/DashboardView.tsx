import { useState, useEffect, useMemo, useCallback } from 'react';
import type { TimeRange } from '../types';
import { useDashboardDetail } from '../hooks/useDashboards';
import { useVariables } from '../hooks/useVariables';
import { fetchDashboardSource, ApiError } from '../api/client';
import { VariableBar } from './VariableBar';
import { RowView } from './RowView';

interface DashboardViewProps {
  path: string;
  timeRange: TimeRange;
  onAuthError: () => void;
  variableValues: Record<string, string>;
  onVariableValuesChange: (values: Record<string, string>) => void;
}

export function DashboardView({ path, timeRange, onAuthError, variableValues, onVariableValuesChange }: DashboardViewProps) {
  const { dashboard, loading, error } = useDashboardDetail(path, onAuthError);
  const { variables, selectedValues, allValues, setVariableValue, loading: varsLoading } =
    useVariables(dashboard?.variables, onAuthError, variableValues);

  const handleVariableChange = useCallback((name: string, value: string) => {
    setVariableValue(name, value);
    onVariableValuesChange({ ...variableValues, [name]: value });
  }, [setVariableValue, onVariableValuesChange, variableValues]);
  const repeatVarNames = useMemo(() => {
    if (!dashboard) return new Set<string>();
    return new Set(dashboard.rows.map((r) => r.repeat).filter(Boolean) as string[]);
  }, [dashboard]);

  // Only hide variables with explicit hide: true
  const hiddenVarNames = useMemo(() => {
    if (!dashboard?.variables) return new Set<string>();
    return new Set(dashboard.variables.filter((v) => v.hide).map((v) => v.name));
  }, [dashboard]);

  const visibleVariables = useMemo(
    () => variables.filter((v) => !hiddenVarNames.has(v.name)),
    [variables, hiddenVarNames],
  );
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
        <>
          {visibleVariables.length > 0 && (
            <VariableBar variables={visibleVariables} repeatVarNames={repeatVarNames} onValueChange={handleVariableChange} />
          )}
          {varsLoading ? (
            <div className="dashboard-loading">Loading variables...</div>
          ) : (
            dashboard.rows.map((row, idx) => {
              if (row.repeat && allValues[row.repeat]) {
                // Repeat this row for each value of the variable
                return allValues[row.repeat].map((value, repeatIdx) => {
                  const repeatValues = { ...selectedValues, [row.repeat!]: value };
                  return (
                    <RowView
                      key={`${idx}-${value}`}
                      row={row}
                      rowIndex={idx * 100 + repeatIdx}
                      timeRange={timeRange}
                      variableValues={repeatValues}
                    />
                  );
                });
              }
              return (
                <RowView
                  key={idx}
                  row={row}
                  rowIndex={idx}
                  timeRange={timeRange}
                  variableValues={selectedValues}
                />
              );
            })
          )}
        </>
      )}
    </div>
  );
}

import { useMemo } from 'react';
import type { Row, TimeRange } from '../types';
import { GraphPanel } from './GraphPanel';
import { MarkdownPanel } from './MarkdownPanel';
import { useQuery } from '../hooks/useQuery';
import { substituteVariables } from '../utils/variables';

interface RowViewProps {
  row: Row;
  rowIndex: number;
  timeRange: TimeRange;
  columns: number;
  variableValues?: Record<string, string>;
}

export function RowView({ row, rowIndex, timeRange, columns, variableValues }: RowViewProps) {
  const vars = variableValues || {};
  const title = substituteVariables(row.title, vars);

  return (
    <div className="row">
      <h2 className="row-title">{title}</h2>
      <div className="row-panels" style={{ gridTemplateColumns: `repeat(${columns}, 1fr)` }}>
        {row.panels.map((panel, idx) => {
          const panelId = `panel-${rowIndex}-${idx}`;
          return (
            <PanelRenderer key={idx} panel={panel} panelId={panelId} timeRange={timeRange} variableValues={vars} />
          );
        })}
      </div>
    </div>
  );
}

interface PanelRendererProps {
  panel: Row['panels'][0];
  panelId: string;
  timeRange: TimeRange;
  variableValues: Record<string, string>;
}

function PanelRenderer({ panel, panelId, timeRange, variableValues }: PanelRendererProps) {
  const substitutedTitle = substituteVariables(panel.title, variableValues);
  const substitutedQuery = useMemo(
    () => panel.query ? substituteVariables(panel.query, variableValues) : undefined,
    [panel.query, variableValues],
  );
  const substitutedContent = useMemo(
    () => panel.content ? substituteVariables(panel.content, variableValues) : undefined,
    [panel.content, variableValues],
  );

  const { data, loading, error } = useQuery(
    panel.type === 'graph' ? substitutedQuery : undefined,
    timeRange,
  );

  if (panel.type === 'markdown') {
    return <MarkdownPanel title={substitutedTitle} content={substitutedContent || ''} id={panelId} />;
  }

  return (
    <GraphPanel
      title={substitutedTitle}
      data={data}
      unit={panel.unit}
      yMin={panel.y_min}
      yMax={panel.y_max}
      legend={panel.legend}
      chartType={panel.chart_type}
      loading={loading}
      error={error}
      id={panelId}
    />
  );
}

import { useMemo } from 'react';
import type { Row, TimeRange } from '../types';
import { GraphPanel } from './GraphPanel';
import { MarkdownPanel } from './MarkdownPanel';
import { useQuery } from '../hooks/useQuery';
import { substituteVariables } from '../utils/variables';
import { getTimeRangeParams } from '../utils/time';

interface RowViewProps {
  row: Row;
  rowIndex: number;
  timeRange: TimeRange;
  variableValues?: Record<string, string>;
}

export function RowView({ row, rowIndex, timeRange, variableValues }: RowViewProps) {
  const vars = variableValues || {};
  const title = substituteVariables(row.title, vars);

  // Calculate default span for panels without explicit span.
  // Uses remaining grid space after explicit spans, clamped to 3–6.
  const explicitTotal = row.panels.reduce((sum, p) => sum + (p.span || 0), 0);
  const nonExplicitCount = row.panels.filter(p => !p.span).length;
  const defaultSpan = nonExplicitCount > 0
    ? Math.min(6, Math.max(3, Math.floor((12 - explicitTotal) / nonExplicitCount)))
    : 0;

  return (
    <div className="row">
      <h2 className="row-title">{title}</h2>
      <div className="row-panels" style={{ gridTemplateColumns: 'repeat(12, 1fr)' }}>
        {row.panels.map((panel, idx) => {
          const panelId = `panel-${rowIndex}-${idx}`;
          const span = panel.span || defaultSpan;
          return (
            <div key={idx} style={{ gridColumn: `span ${span}` }}>
              <PanelRenderer panel={panel} panelId={panelId} timeRange={timeRange} variableValues={vars} />
            </div>
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
  const substitutedDatasource = useMemo(
    () => panel.datasource ? substituteVariables(panel.datasource, variableValues) : undefined,
    [panel.datasource, variableValues],
  );

  const { data, loading, error } = useQuery(
    panel.type === 'graph' ? substitutedQuery : undefined,
    timeRange,
    substitutedDatasource,
  );

  const stepSeconds = useMemo(() => {
    const { step } = getTimeRangeParams(timeRange);
    return parseInt(step.replace('s', ''), 10);
  }, [timeRange]);

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
      legendDisplay={panel.legend_display}
      legendPosition={panel.legend_position}
      legendAlign={panel.legend_align}
      legendMaxHeight={panel.legend_max_height}
      legendMaxWidth={panel.legend_max_width}
      thresholds={panel.thresholds}
      chartType={panel.chart_type}
      stacked={panel.stacked}
      yScale={panel.y_scale}
      stepSeconds={stepSeconds}
      loading={loading}
      error={error}
      id={panelId}
    />
  );
}

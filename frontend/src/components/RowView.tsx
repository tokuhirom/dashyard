import type { Row, TimeRange } from '../types';
import { GraphPanel } from './GraphPanel';
import { MarkdownPanel } from './MarkdownPanel';
import { useQuery } from '../hooks/useQuery';

interface RowViewProps {
  row: Row;
  rowIndex: number;
  timeRange: TimeRange;
  columns: number;
}

export function RowView({ row, rowIndex, timeRange, columns }: RowViewProps) {
  return (
    <div className="row">
      <h2 className="row-title">{row.title}</h2>
      <div className="row-panels" style={{ gridTemplateColumns: `repeat(${columns}, 1fr)` }}>
        {row.panels.map((panel, idx) => {
          const panelId = `panel-${rowIndex}-${idx}`;
          return (
            <PanelRenderer key={idx} panel={panel} panelId={panelId} timeRange={timeRange} />
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
}

function PanelRenderer({ panel, panelId, timeRange }: PanelRendererProps) {
  const { data, loading, error } = useQuery(
    panel.type === 'graph' ? panel.query : undefined,
    timeRange,
  );

  if (panel.type === 'markdown') {
    return <MarkdownPanel title={panel.title} content={panel.content || ''} id={panelId} />;
  }

  return (
    <GraphPanel
      title={panel.title}
      data={data}
      unit={panel.unit}
      legend={panel.legend}
      chartType={panel.chart_type}
      loading={loading}
      error={error}
      id={panelId}
    />
  );
}

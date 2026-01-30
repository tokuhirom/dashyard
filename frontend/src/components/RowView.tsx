import type { Row, TimeRange } from '../types';
import { GraphPanel } from './GraphPanel';
import { MarkdownPanel } from './MarkdownPanel';
import { useQuery } from '../hooks/useQuery';

interface RowViewProps {
  row: Row;
  timeRange: TimeRange;
}

export function RowView({ row, timeRange }: RowViewProps) {
  return (
    <div className="row">
      <h2 className="row-title">{row.title}</h2>
      <div className="row-panels">
        {row.panels.map((panel, idx) => (
          <PanelRenderer key={idx} panel={panel} timeRange={timeRange} />
        ))}
      </div>
    </div>
  );
}

interface PanelRendererProps {
  panel: Row['panels'][0];
  timeRange: TimeRange;
}

function PanelRenderer({ panel, timeRange }: PanelRendererProps) {
  const { data, loading, error } = useQuery(
    panel.type === 'graph' ? panel.query : undefined,
    timeRange,
  );

  if (panel.type === 'markdown') {
    return <MarkdownPanel title={panel.title} content={panel.content || ''} />;
  }

  return (
    <GraphPanel
      title={panel.title}
      data={data}
      unit={panel.unit}
      loading={loading}
      error={error}
    />
  );
}

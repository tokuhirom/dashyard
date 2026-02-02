import { useState, useEffect, useCallback } from 'react';
import { createPortal } from 'react-dom';
import { Line, Bar, Scatter } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  LogarithmicScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  Filler,
} from 'chart.js';
import annotationPlugin from 'chartjs-plugin-annotation';
import 'chartjs-adapter-date-fns';
import type { QueryResponse, Threshold } from '../types';
import { getYAxisTickCallback } from '../utils/units';

ChartJS.register(
  CategoryScale,
  LinearScale,
  LogarithmicScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  Filler,
  annotationPlugin,
);

interface GraphPanelProps {
  title: string;
  data: QueryResponse | null;
  unit?: string;
  yMin?: number;
  yMax?: number;
  legend?: string;
  thresholds?: Threshold[];
  chartType?: 'line' | 'bar' | 'area' | 'scatter';
  stacked?: boolean;
  yScale?: 'linear' | 'log';
  loading: boolean;
  error: string | null;
  id?: string;
}

function buildAnnotations(thresholds?: Threshold[]) {
  if (!thresholds || thresholds.length === 0) return {};
  return {
    annotation: {
      annotations: thresholds.map((th) => {
        const color = th.color || '#ef4444';
        return {
          type: 'line' as const,
          scaleID: 'y',
          value: th.value,
          borderColor: color,
          borderWidth: 2,
          borderDash: [6, 3],
          drawTime: 'afterDatasetsDraw' as const,
          ...(th.label ? {
            label: {
              display: true,
              content: th.label,
              position: 'end' as const,
              backgroundColor: color,
              color: '#fff',
              font: { size: 11 },
            },
          } : {}),
        };
      }),
    },
  };
}

const COLORS = [
  '#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6',
  '#ec4899', '#06b6d4', '#84cc16', '#f97316', '#6366f1',
];

function buildLabel(metric: Record<string, string>, legend?: string): string {
  if (legend) {
    return legend.replace(/\{([^}]+)\}/g, (_, key) => metric[key] ?? '');
  }
  const entries = Object.entries(metric).filter(([k]) => k !== '__name__');
  if (entries.length === 0) {
    return metric['__name__'] || 'value';
  }
  return entries.map(([k, v]) => `${k}="${v}"`).join(', ');
}

export function GraphPanel({ title, data, unit, yMin, yMax, legend, thresholds, chartType, stacked, yScale, loading, error, id }: GraphPanelProps) {
  const [expanded, setExpanded] = useState(false);

  const close = useCallback(() => setExpanded(false), []);

  useEffect(() => {
    if (!expanded) return;
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') close();
    };
    document.addEventListener('keydown', onKeyDown);
    return () => document.removeEventListener('keydown', onKeyDown);
  }, [expanded, close]);

  const titleContent = (
    <h3 className="panel-title">
      {title}
      {id && <a href={`#${id}`} className="panel-anchor">#</a>}
      {!loading && !error && data?.data?.result?.length && (
        <button className="panel-expand-btn" onClick={() => setExpanded(true)} title="Expand">&#x2922;</button>
      )}
    </h3>
  );

  if (loading) {
    return (
      <div className="panel graph-panel" id={id}>
        {titleContent}
        <div className="panel-loading">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="panel graph-panel" id={id}>
        {titleContent}
        <div className="panel-error">{error}</div>
      </div>
    );
  }

  if (!data || !data.data?.result?.length) {
    return (
      <div className="panel graph-panel" id={id}>
        {titleContent}
        <div className="panel-empty">No data</div>
      </div>
    );
  }

  const effectiveType = chartType || 'line';

  const shouldStack = stacked && (effectiveType === 'line' || effectiveType === 'area' || effectiveType === 'bar');

  const datasets = data.data.result.map((result, idx) => ({
    label: buildLabel(result.metric, legend),
    data: result.values.map(([ts, val]) => ({
      x: ts * 1000,
      y: parseFloat(val),
    })),
    borderColor: COLORS[idx % COLORS.length],
    backgroundColor: shouldStack ? COLORS[idx % COLORS.length] + '80' : COLORS[idx % COLORS.length] + '20',
    borderWidth: 1.5,
    pointRadius: effectiveType === 'scatter' ? 3 : 0,
    tension: 0.1,
    fill: effectiveType === 'area' || (shouldStack && effectiveType === 'line') ? 'origin' : false,
  }));

  const tickCallback = getYAxisTickCallback(unit);

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const buildOptions = (isExpanded: boolean): any => ({
    responsive: true,
    maintainAspectRatio: !isExpanded,
    ...(isExpanded ? {} : { aspectRatio: 2.5 }),
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
    plugins: {
      legend: {
        position: 'bottom' as const,
        ...(isExpanded ? {} : { maxHeight: 60 }),
        labels: {
          boxWidth: 12,
          usePointStyle: true,
        },
      },
      ...buildAnnotations(thresholds),
    },
    scales: {
      x: {
        type: 'time' as const,
        ...(shouldStack ? { stacked: true } : {}),
        time: {
          tooltipFormat: 'HH:mm:ss',
        },
        ticks: {
          maxTicksLimit: isExpanded ? 16 : 8,
        },
      },
      y: {
        ...(yScale === 'log' ? { type: 'logarithmic' as const } : { beginAtZero: true }),
        ...(shouldStack ? { stacked: true } : {}),
        ...(unit === 'percent' ? { min: 0, max: 100 } : {}),
        ...(yMin !== undefined ? { min: yMin } : {}),
        ...(yMax !== undefined ? { max: yMax } : {}),
        ticks: {
          callback: tickCallback,
        },
      },
    },
  });

  const renderChart = (isExpanded: boolean) => {
    const opts = buildOptions(isExpanded);
    if (effectiveType === 'bar') {
      return <Bar data={{ datasets }} options={opts} />;
    } else if (effectiveType === 'scatter') {
      return <Scatter data={{ datasets }} options={opts} />;
    } else {
      return <Line data={{ datasets }} options={opts} />;
    }
  };

  return (
    <>
      <div className="panel graph-panel" id={id}>
        {titleContent}
        <div className="panel-chart">
          {renderChart(false)}
        </div>
      </div>
      {expanded && createPortal(
        <div className="modal-backdrop" onClick={close}>
          <div className="panel-modal" onClick={(e) => e.stopPropagation()}>
            <div className="panel-modal-header">
              <h3>{title}</h3>
              <button className="panel-modal-close" onClick={close}>&times;</button>
            </div>
            <div className="panel-chart" style={{ flex: 1 }}>
              {renderChart(true)}
            </div>
          </div>
        </div>,
        document.body
      )}
    </>
  );
}

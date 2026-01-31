import { Line, Bar, Scatter } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
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
import type { PrometheusResponse, Threshold } from '../types';
import { getYAxisTickCallback } from '../utils/units';

ChartJS.register(
  CategoryScale,
  LinearScale,
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
  data: PrometheusResponse | null;
  unit?: string;
  yMin?: number;
  yMax?: number;
  legend?: string;
  thresholds?: Threshold[];
  chartType?: 'line' | 'bar' | 'area' | 'scatter';
  stacked?: boolean;
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

export function GraphPanel({ title, data, unit, yMin, yMax, legend, thresholds, chartType, stacked, loading, error, id }: GraphPanelProps) {
  const titleContent = (
    <h3 className="panel-title">
      {title}
      {id && <a href={`#${id}`} className="panel-anchor">#</a>}
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
  const options: any = {
    responsive: true,
    maintainAspectRatio: true,
    aspectRatio: 2.5,
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
    plugins: {
      legend: {
        position: 'bottom' as const,
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
          maxTicksLimit: 8,
        },
      },
      y: {
        beginAtZero: true,
        ...(shouldStack ? { stacked: true } : {}),
        ...(unit === 'percent' ? { min: 0, max: 100 } : {}),
        ...(yMin !== undefined ? { min: yMin } : {}),
        ...(yMax !== undefined ? { max: yMax } : {}),
        ticks: {
          callback: tickCallback,
        },
      },
    },
  };

  let chart;
  if (effectiveType === 'bar') {
    chart = <Bar data={{ datasets }} options={options} />;
  } else if (effectiveType === 'scatter') {
    chart = <Scatter data={{ datasets }} options={options} />;
  } else {
    chart = <Line data={{ datasets }} options={options} />;
  }

  return (
    <div className="panel graph-panel" id={id}>
      {titleContent}
      <div className="panel-chart">
        {chart}
      </div>
    </div>
  );
}

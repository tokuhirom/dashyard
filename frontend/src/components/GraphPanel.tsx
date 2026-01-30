import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
} from 'chart.js';
import 'chartjs-adapter-date-fns';
import type { PrometheusResponse } from '../types';
import { getYAxisTickCallback } from '../utils/units';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
);

interface GraphPanelProps {
  title: string;
  data: PrometheusResponse | null;
  unit?: string;
  legend?: string;
  loading: boolean;
  error: string | null;
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

export function GraphPanel({ title, data, unit, legend, loading, error }: GraphPanelProps) {
  if (loading) {
    return (
      <div className="panel graph-panel">
        <h3 className="panel-title">{title}</h3>
        <div className="panel-loading">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="panel graph-panel">
        <h3 className="panel-title">{title}</h3>
        <div className="panel-error">{error}</div>
      </div>
    );
  }

  if (!data || !data.data?.result?.length) {
    return (
      <div className="panel graph-panel">
        <h3 className="panel-title">{title}</h3>
        <div className="panel-empty">No data</div>
      </div>
    );
  }

  const datasets = data.data.result.map((result, idx) => ({
    label: buildLabel(result.metric, legend),
    data: result.values.map(([ts, val]) => ({
      x: ts * 1000,
      y: parseFloat(val),
    })),
    borderColor: COLORS[idx % COLORS.length],
    backgroundColor: COLORS[idx % COLORS.length] + '20',
    borderWidth: 1.5,
    pointRadius: 0,
    tension: 0.1,
    fill: false,
  }));

  const tickCallback = getYAxisTickCallback(unit);

  const options = {
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
    },
    scales: {
      x: {
        type: 'time' as const,
        time: {
          tooltipFormat: 'HH:mm:ss',
        },
        ticks: {
          maxTicksLimit: 8,
        },
      },
      y: {
        beginAtZero: true,
        ticks: {
          callback: tickCallback,
        },
      },
    },
  };

  return (
    <div className="panel graph-panel">
      <h3 className="panel-title">{title}</h3>
      <div className="panel-chart">
        <Line data={{ datasets }} options={options} />
      </div>
    </div>
  );
}

import { Line, Bar, Scatter, Pie, Doughnut } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  Filler,
} from 'chart.js';
import 'chartjs-adapter-date-fns';
import type { PrometheusResponse } from '../types';
import { getYAxisTickCallback } from '../utils/units';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  Filler,
);

interface GraphPanelProps {
  title: string;
  data: PrometheusResponse | null;
  unit?: string;
  legend?: string;
  chartType?: 'line' | 'bar' | 'area' | 'scatter' | 'pie' | 'doughnut';
  loading: boolean;
  error: string | null;
  id?: string;
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

export function GraphPanel({ title, data, unit, legend, chartType, loading, error, id }: GraphPanelProps) {
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

  if (effectiveType === 'pie' || effectiveType === 'doughnut') {
    const labels = data.data.result.map((result) => buildLabel(result.metric, legend));
    const values = data.data.result.map((result) => {
      const lastValue = result.values[result.values.length - 1];
      return lastValue ? parseFloat(lastValue[1]) : 0;
    });
    const backgroundColors = data.data.result.map((_, idx) => COLORS[idx % COLORS.length]);

    const chartData = {
      labels,
      datasets: [{
        data: values,
        backgroundColor: backgroundColors,
        borderWidth: 1,
      }],
    };

    const options = {
      responsive: true,
      maintainAspectRatio: true,
      aspectRatio: 2.5,
      plugins: {
        legend: {
          position: 'bottom' as const,
          labels: {
            boxWidth: 12,
            usePointStyle: true,
          },
        },
      },
    };

    const ChartComponent = effectiveType === 'pie' ? Pie : Doughnut;

    return (
      <div className="panel graph-panel" id={id}>
        {titleContent}
        <div className="panel-chart">
          <ChartComponent data={chartData} options={options} />
        </div>
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
    pointRadius: effectiveType === 'scatter' ? 3 : 0,
    tension: 0.1,
    fill: effectiveType === 'area',
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

export function formatBytes(value: number): string {
  if (value === 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const k = 1024;
  const i = Math.floor(Math.log(Math.abs(value)) / Math.log(k));
  const idx = Math.max(0, Math.min(i, units.length - 1));
  return `${(value / Math.pow(k, idx)).toFixed(idx > 0 ? 1 : 0)} ${units[idx]}`;
}

export function formatPercent(value: number): string {
  return `${value.toFixed(1)}%`;
}

export function formatCount(value: number): string {
  return value.toLocaleString('en-US', { maximumFractionDigits: 2 });
}

export function formatSeconds(value: number): string {
  const abs = Math.abs(value);
  const sign = value < 0 ? '-' : '';
  if (abs === 0) return '0s';
  if (abs < 0.001) return `${sign}${(abs * 1_000_000).toFixed(0)}µs`;
  if (abs < 1) return `${sign}${(abs * 1000).toFixed(1)}ms`;
  if (abs < 60) return `${sign}${abs.toFixed(2)}s`;
  if (abs < 3600) return `${sign}${(abs / 60).toFixed(1)}m`;
  return `${sign}${(abs / 3600).toFixed(1)}h`;
}

export function formatValue(value: number, unit?: string): string {
  switch (unit) {
    case 'bytes':
      return formatBytes(value);
    case 'percent':
      return formatPercent(value);
    case 'seconds':
      return formatSeconds(value);
    case 'count':
      return formatCount(value);
    default:
      return formatCount(value);
  }
}

export function getYAxisTickCallback(unit?: string): (value: number | string) => string {
  return (value: number | string) => {
    const num = typeof value === 'string' ? parseFloat(value) : value;
    return formatValue(num, unit);
  };
}

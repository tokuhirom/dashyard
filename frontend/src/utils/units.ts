export function formatBytes(value: number): string {
  if (value === 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const k = 1024;
  const i = Math.floor(Math.log(Math.abs(value)) / Math.log(k));
  const idx = Math.min(i, units.length - 1);
  return `${(value / Math.pow(k, idx)).toFixed(idx > 0 ? 1 : 0)} ${units[idx]}`;
}

export function formatPercent(value: number): string {
  return `${value.toFixed(1)}%`;
}

export function formatCount(value: number): string {
  return value.toLocaleString('en-US', { maximumFractionDigits: 2 });
}

export function formatValue(value: number, unit?: string): string {
  switch (unit) {
    case 'bytes':
      return formatBytes(value);
    case 'percent':
      return formatPercent(value);
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

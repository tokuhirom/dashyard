import type { TimeRange } from '../types';

export const TIME_RANGES: TimeRange[] = [
  { label: 'Last 15 minutes', value: '15m', duration: 15 * 60, step: '15s' },
  { label: 'Last 30 minutes', value: '30m', duration: 30 * 60, step: '30s' },
  { label: 'Last 1 hour', value: '1h', duration: 60 * 60, step: '60s' },
  { label: 'Last 3 hours', value: '3h', duration: 3 * 60 * 60, step: '120s' },
  { label: 'Last 6 hours', value: '6h', duration: 6 * 60 * 60, step: '240s' },
  { label: 'Last 12 hours', value: '12h', duration: 12 * 60 * 60, step: '480s' },
  { label: 'Last 24 hours', value: '24h', duration: 24 * 60 * 60, step: '900s' },
  { label: 'Last 3 days', value: '3d', duration: 3 * 24 * 60 * 60, step: '3600s' },
  { label: 'Last 7 days', value: '7d', duration: 7 * 24 * 60 * 60, step: '7200s' },
];

export const DEFAULT_TIME_RANGE = TIME_RANGES[2]; // 1 hour

export function getTimeRangeParams(range_: TimeRange): { start: number; end: number; step: string } {
  const end = Math.floor(Date.now() / 1000);
  const start = end - range_.duration;
  return { start, end, step: range_.step };
}

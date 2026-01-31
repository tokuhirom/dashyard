import type { TimeRange, RelativeTimeRange } from '../types';

export const TIME_RANGES: RelativeTimeRange[] = [
  { type: 'relative', label: 'Last 15 minutes', value: '15m', duration: 15 * 60, step: '15s' },
  { type: 'relative', label: 'Last 30 minutes', value: '30m', duration: 30 * 60, step: '30s' },
  { type: 'relative', label: 'Last 1 hour', value: '1h', duration: 60 * 60, step: '60s' },
  { type: 'relative', label: 'Last 3 hours', value: '3h', duration: 3 * 60 * 60, step: '120s' },
  { type: 'relative', label: 'Last 6 hours', value: '6h', duration: 6 * 60 * 60, step: '240s' },
  { type: 'relative', label: 'Last 12 hours', value: '12h', duration: 12 * 60 * 60, step: '480s' },
  { type: 'relative', label: 'Last 24 hours', value: '24h', duration: 24 * 60 * 60, step: '900s' },
  { type: 'relative', label: 'Last 3 days', value: '3d', duration: 3 * 24 * 60 * 60, step: '3600s' },
  { type: 'relative', label: 'Last 7 days', value: '7d', duration: 7 * 24 * 60 * 60, step: '7200s' },
];

export const DEFAULT_TIME_RANGE = TIME_RANGES[2]; // 1 hour

export function computeStep(durationSeconds: number): string {
  if (durationSeconds <= 15 * 60) return '15s';
  if (durationSeconds <= 30 * 60) return '30s';
  if (durationSeconds <= 60 * 60) return '60s';
  if (durationSeconds <= 3 * 60 * 60) return '120s';
  if (durationSeconds <= 6 * 60 * 60) return '240s';
  if (durationSeconds <= 12 * 60 * 60) return '480s';
  if (durationSeconds <= 24 * 60 * 60) return '900s';
  if (durationSeconds <= 3 * 24 * 60 * 60) return '3600s';
  return '7200s';
}

export function getTimeRangeParams(range_: TimeRange): { start: number; end: number; step: string } {
  if (range_.type === 'absolute') {
    return { start: range_.start, end: range_.end, step: range_.step };
  }
  const end = Math.floor(Date.now() / 1000);
  const start = end - range_.duration;
  return { start, end, step: range_.step };
}

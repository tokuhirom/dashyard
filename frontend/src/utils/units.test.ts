import { describe, it, expect } from 'vitest';
import { formatBytes, formatPercent, formatCount, formatSeconds, formatValue, getYAxisTickCallback } from './units';

describe('formatBytes', () => {
  it('formats zero', () => {
    expect(formatBytes(0)).toBe('0 B');
  });

  it('formats bytes', () => {
    expect(formatBytes(500)).toBe('500 B');
  });

  it('formats kilobytes', () => {
    expect(formatBytes(1024)).toBe('1.0 KB');
    expect(formatBytes(1536)).toBe('1.5 KB');
  });

  it('formats megabytes', () => {
    expect(formatBytes(1048576)).toBe('1.0 MB');
  });

  it('formats gigabytes', () => {
    expect(formatBytes(1073741824)).toBe('1.0 GB');
  });

  it('formats terabytes', () => {
    expect(formatBytes(1099511627776)).toBe('1.0 TB');
  });

  it('handles negative values', () => {
    expect(formatBytes(-1024)).toBe('-1.0 KB');
  });
});

describe('formatPercent', () => {
  it('formats with one decimal place', () => {
    expect(formatPercent(75)).toBe('75.0%');
    expect(formatPercent(99.9)).toBe('99.9%');
    expect(formatPercent(0)).toBe('0.0%');
  });

  it('rounds to one decimal', () => {
    expect(formatPercent(33.333)).toBe('33.3%');
  });
});

describe('formatCount', () => {
  it('formats small numbers', () => {
    expect(formatCount(0)).toBe('0');
    expect(formatCount(42)).toBe('42');
  });

  it('formats with thousands separators', () => {
    expect(formatCount(1000)).toBe('1,000');
    expect(formatCount(1234567)).toBe('1,234,567');
  });

  it('formats decimals with max 2 fraction digits', () => {
    expect(formatCount(1.5)).toBe('1.5');
    expect(formatCount(1.999)).toBe('2');
  });
});

describe('formatSeconds', () => {
  it('formats zero', () => {
    expect(formatSeconds(0)).toBe('0s');
  });

  it('formats microseconds', () => {
    expect(formatSeconds(0.0005)).toBe('500µs');
  });

  it('formats milliseconds', () => {
    expect(formatSeconds(0.2)).toBe('200.0ms');
    expect(formatSeconds(0.999)).toBe('999.0ms');
  });

  it('formats seconds', () => {
    expect(formatSeconds(1)).toBe('1.00s');
    expect(formatSeconds(30.5)).toBe('30.50s');
  });

  it('formats minutes', () => {
    expect(formatSeconds(120)).toBe('2.0m');
    expect(formatSeconds(90)).toBe('1.5m');
  });

  it('formats hours', () => {
    expect(formatSeconds(3600)).toBe('1.0h');
    expect(formatSeconds(7200)).toBe('2.0h');
  });

  it('handles negative values', () => {
    expect(formatSeconds(-0.2)).toBe('-200.0ms');
    expect(formatSeconds(-60)).toBe('-1.0m');
  });
});

describe('formatValue', () => {
  it('dispatches to formatBytes', () => {
    expect(formatValue(1024, 'bytes')).toBe('1.0 KB');
  });

  it('dispatches to formatPercent', () => {
    expect(formatValue(50, 'percent')).toBe('50.0%');
  });

  it('dispatches to formatSeconds', () => {
    expect(formatValue(0.5, 'seconds')).toBe('500.0ms');
  });

  it('dispatches to formatCount', () => {
    expect(formatValue(1000, 'count')).toBe('1,000');
  });

  it('defaults to formatCount when unit is undefined', () => {
    expect(formatValue(1000)).toBe('1,000');
  });
});

describe('getYAxisTickCallback', () => {
  it('returns a callback that formats numbers', () => {
    const cb = getYAxisTickCallback('percent');
    expect(cb(50)).toBe('50.0%');
  });

  it('handles string values', () => {
    const cb = getYAxisTickCallback('bytes');
    expect(cb('1024')).toBe('1.0 KB');
  });
});

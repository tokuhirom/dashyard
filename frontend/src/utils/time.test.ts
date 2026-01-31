import { describe, it, expect, vi, afterEach } from 'vitest';
import { TIME_RANGES, DEFAULT_TIME_RANGE, computeStep, getTimeRangeParams } from './time';
import type { AbsoluteTimeRange } from '../types';

describe('TIME_RANGES', () => {
  it('contains 9 predefined ranges', () => {
    expect(TIME_RANGES).toHaveLength(9);
  });

  it('all have type relative', () => {
    for (const r of TIME_RANGES) {
      expect(r.type).toBe('relative');
    }
  });
});

describe('DEFAULT_TIME_RANGE', () => {
  it('is 1 hour', () => {
    expect(DEFAULT_TIME_RANGE.value).toBe('1h');
    expect(DEFAULT_TIME_RANGE.duration).toBe(3600);
  });
});

describe('computeStep', () => {
  it('returns 15s for up to 15 minutes', () => {
    expect(computeStep(15 * 60)).toBe('15s');
    expect(computeStep(60)).toBe('15s');
  });

  it('returns 30s for up to 30 minutes', () => {
    expect(computeStep(15 * 60 + 1)).toBe('30s');
    expect(computeStep(30 * 60)).toBe('30s');
  });

  it('returns 60s for up to 1 hour', () => {
    expect(computeStep(30 * 60 + 1)).toBe('60s');
    expect(computeStep(60 * 60)).toBe('60s');
  });

  it('returns 120s for up to 3 hours', () => {
    expect(computeStep(3 * 60 * 60)).toBe('120s');
  });

  it('returns 240s for up to 6 hours', () => {
    expect(computeStep(6 * 60 * 60)).toBe('240s');
  });

  it('returns 480s for up to 12 hours', () => {
    expect(computeStep(12 * 60 * 60)).toBe('480s');
  });

  it('returns 900s for up to 24 hours', () => {
    expect(computeStep(24 * 60 * 60)).toBe('900s');
  });

  it('returns 3600s for up to 3 days', () => {
    expect(computeStep(3 * 24 * 60 * 60)).toBe('3600s');
  });

  it('returns 7200s for longer durations', () => {
    expect(computeStep(3 * 24 * 60 * 60 + 1)).toBe('7200s');
    expect(computeStep(7 * 24 * 60 * 60)).toBe('7200s');
  });
});

describe('getTimeRangeParams', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns absolute range params as-is', () => {
    const range: AbsoluteTimeRange = {
      type: 'absolute',
      label: 'Custom',
      start: 1000,
      end: 2000,
      step: '15s',
    };
    const params = getTimeRangeParams(range);
    expect(params).toEqual({ start: 1000, end: 2000, step: '15s' });
  });

  it('computes start/end from relative range', () => {
    const now = 1700000000;
    vi.spyOn(Date, 'now').mockReturnValue(now * 1000);

    const params = getTimeRangeParams(DEFAULT_TIME_RANGE);
    expect(params.end).toBe(now);
    expect(params.start).toBe(now - 3600);
    expect(params.step).toBe('60s');
  });
});

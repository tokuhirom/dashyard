import { describe, it, expect } from 'vitest';
import { parseLabelValuesQuery, substituteVariables } from './variables';

describe('parseLabelValuesQuery', () => {
  it('parses a standard query', () => {
    const result = parseLabelValuesQuery('label_values(up, instance)');
    expect(result).toEqual({ metric: 'up', label: 'instance' });
  });

  it('handles extra whitespace', () => {
    const result = parseLabelValuesQuery('label_values(  up  ,  instance  )');
    expect(result).toEqual({ metric: 'up', label: 'instance' });
  });

  it('handles metric with braces', () => {
    const result = parseLabelValuesQuery('label_values(node_cpu_seconds_total{mode="idle"}, cpu)');
    expect(result).toEqual({ metric: 'node_cpu_seconds_total{mode="idle"}', label: 'cpu' });
  });

  it('returns null for invalid format', () => {
    expect(parseLabelValuesQuery('not a query')).toBeNull();
    expect(parseLabelValuesQuery('label_values()')).toBeNull();
    expect(parseLabelValuesQuery('label_values(metric)')).toBeNull();
    expect(parseLabelValuesQuery('')).toBeNull();
  });
});

describe('substituteVariables', () => {
  it('replaces ${var} form', () => {
    expect(substituteVariables('rate(http_requests{job="${job}"}[5m])', { job: 'api' }))
      .toBe('rate(http_requests{job="api"}[5m])');
  });

  it('replaces $var form', () => {
    expect(substituteVariables('up{instance=~"$instance"}', { instance: 'localhost:9090' }))
      .toBe('up{instance=~"localhost:9090"}');
  });

  it('does not replace $var when followed by word characters', () => {
    expect(substituteVariables('$device_total', { device: 'eth0' }))
      .toBe('$device_total');
  });

  it('replaces multiple variables', () => {
    const result = substituteVariables('up{job="$job", instance="$instance"}', {
      job: 'api',
      instance: 'localhost',
    });
    expect(result).toBe('up{job="api", instance="localhost"}');
  });

  it('handles longer variable names first to avoid prefix collisions', () => {
    const result = substituteVariables('$dev $device', {
      dev: 'A',
      device: 'B',
    });
    expect(result).toBe('A B');
  });

  it('returns template unchanged when no variables', () => {
    expect(substituteVariables('up{job="api"}', {})).toBe('up{job="api"}');
  });

  it('returns empty/undefined template as-is', () => {
    expect(substituteVariables('', { job: 'api' })).toBe('');
  });
});

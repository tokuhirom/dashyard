import { describe, it, expect } from 'vitest';
import { buildLabel } from './legend';

describe('buildLabel', () => {
  const metric = {
    __name__: 'http_requests_total',
    method: 'GET',
    status: '200',
    instance_id: '550e8400-e29b-41d4-a716-446655440000',
    request_id: '0192d2a4-b6c0-7a3a-8e1f-4a5b6c7d8e9f',
  };

  describe('basic interpolation', () => {
    it('replaces label name with value', () => {
      expect(buildLabel(metric, '{method}')).toBe('GET');
    });

    it('replaces multiple labels', () => {
      expect(buildLabel(metric, '{method} {status}')).toBe('GET 200');
    });

    it('returns empty string for missing label', () => {
      expect(buildLabel(metric, '{nonexistent}')).toBe('');
    });

    it('falls back to all labels when no legend', () => {
      const m = { __name__: 'up', job: 'api', instance: 'localhost' };
      expect(buildLabel(m)).toBe('job="api", instance="localhost"');
    });

    it('uses __name__ when no other labels', () => {
      expect(buildLabel({ __name__: 'up' })).toBe('up');
    });

    it('returns "value" when no labels at all', () => {
      expect(buildLabel({})).toBe('value');
    });
  });

  describe('trunc', () => {
    it('truncates to N chars with ...', () => {
      expect(buildLabel(metric, '{instance_id | trunc(8)}')).toBe('550e8400...');
    });

    it('does not truncate when value is shorter than N', () => {
      expect(buildLabel(metric, '{method | trunc(10)}')).toBe('GET');
    });

    it('does not truncate when value equals N', () => {
      expect(buildLabel(metric, '{method | trunc(3)}')).toBe('GET');
    });
  });

  describe('suffix', () => {
    it('shows last N chars with ...', () => {
      expect(buildLabel(metric, '{instance_id | suffix(8)}')).toBe('...55440000');
    });

    it('does not truncate when value is shorter than N', () => {
      expect(buildLabel(metric, '{method | suffix(10)}')).toBe('GET');
    });
  });

  describe('upper / lower', () => {
    it('converts to uppercase', () => {
      expect(buildLabel(metric, '{method | upper}')).toBe('GET');
    });

    it('converts to lowercase', () => {
      expect(buildLabel(metric, '{method | lower}')).toBe('get');
    });
  });

  describe('replace', () => {
    it('replaces substring', () => {
      expect(buildLabel(metric, '{method | replace("GET","POST")}')).toBe('POST');
    });

    it('replaces all occurrences', () => {
      const m = { val: 'a-b-c' };
      expect(buildLabel(m, '{val | replace("-","_")}')).toBe('a_b_c');
    });

    it('returns value unchanged on invalid args', () => {
      expect(buildLabel(metric, '{method | replace(bad)}')).toBe('GET');
    });
  });

  describe('pipe chaining', () => {
    it('applies multiple functions in order', () => {
      expect(buildLabel(metric, '{instance_id | trunc(8) | upper}')).toBe('550E8400...');
    });

    it('chains suffix and upper', () => {
      expect(buildLabel(metric, '{method | lower | replace("get","post")}')).toBe('post');
    });
  });

  describe('mixed template', () => {
    it('handles pipes in some placeholders and plain in others', () => {
      expect(buildLabel(metric, '{method} {instance_id | trunc(8)}')).toBe('GET 550e8400...');
    });
  });

  describe('edge cases', () => {
    it('handles unknown function gracefully', () => {
      expect(buildLabel(metric, '{method | nonexistent(5)}')).toBe('GET');
    });

    it('handles empty pipe expression', () => {
      expect(buildLabel(metric, '{method |}')).toBe('GET');
    });

    it('handles trunc(0)', () => {
      expect(buildLabel(metric, '{method | trunc(0)}')).toBe('GET');
    });

    it('handles negative trunc', () => {
      expect(buildLabel(metric, '{method | trunc(-1)}')).toBe('GET');
    });
  });
});

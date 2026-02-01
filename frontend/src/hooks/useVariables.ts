import { useState, useEffect } from 'react';
import type { Variable } from '../types';
import { fetchLabelValues, ApiError } from '../api/client';
import { parseLabelValuesQuery } from '../utils/variables';

export interface VariableState {
  name: string;
  label: string;
  values: string[];
  selected: string;
  loading: boolean;
  error: string | null;
}

interface UseVariablesResult {
  variables: VariableState[];
  selectedValues: Record<string, string>;
  allValues: Record<string, string[]>;
  setVariableValue: (name: string, value: string) => void;
  loading: boolean;
}

export function useVariables(
  definitions: Variable[] | undefined,
  onAuthError: () => void,
): UseVariablesResult {
  const [variables, setVariables] = useState<VariableState[]>([]);

  useEffect(() => {
    if (!definitions || definitions.length === 0) {
      setVariables([]);
      return;
    }

    // Initialize variable states
    const initial: VariableState[] = definitions.map((def) => ({
      name: def.name,
      label: def.label || def.name,
      values: [],
      selected: '',
      loading: true,
      error: null,
    }));
    setVariables(initial);

    // Fetch all variable values in parallel
    definitions.forEach((def, idx) => {
      const parsed = parseLabelValuesQuery(def.query);
      if (!parsed) {
        setVariables((prev) => {
          const next = [...prev];
          next[idx] = { ...next[idx], loading: false, error: 'Invalid query format' };
          return next;
        });
        return;
      }

      fetchLabelValues(parsed.label, parsed.metric, def.datasource)
        .then((resp) => {
          const values = resp.data || [];
          setVariables((prev) => {
            const next = [...prev];
            next[idx] = {
              ...next[idx],
              values,
              selected: values[0] || '',
              loading: false,
            };
            return next;
          });
        })
        .catch((err) => {
          if (err instanceof ApiError && err.status === 401) {
            onAuthError();
          }
          setVariables((prev) => {
            const next = [...prev];
            next[idx] = { ...next[idx], loading: false, error: err.message };
            return next;
          });
        });
    });
  }, [definitions, onAuthError]);

  const setVariableValue = (name: string, value: string) => {
    setVariables((prev) =>
      prev.map((v) => (v.name === name ? { ...v, selected: value } : v)),
    );
  };

  const selectedValues: Record<string, string> = {};
  const allValues: Record<string, string[]> = {};
  for (const v of variables) {
    if (v.selected) {
      selectedValues[v.name] = v.selected;
    }
    allValues[v.name] = v.values;
  }

  const loading = variables.some((v) => v.loading);

  return { variables, selectedValues, allValues, setVariableValue, loading };
}

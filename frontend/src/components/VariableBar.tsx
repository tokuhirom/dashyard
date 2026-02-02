import type { VariableState } from '../hooks/useVariables';

interface VariableBarProps {
  variables: VariableState[];
  repeatVarNames: Set<string>;
  onValueChange: (name: string, value: string) => void;
}

export function VariableBar({ variables, repeatVarNames, onValueChange }: VariableBarProps) {
  if (variables.length === 0) return null;

  return (
    <div className="variable-bar">
      {variables.map((variable) => {
        const isRepeat = repeatVarNames.has(variable.name);
        return (
          <div key={variable.name} className="variable-selector">
            <label className="variable-label">
              {variable.label}
              {isRepeat && <span className="variable-repeat-badge">row repeat</span>}
            </label>
            {variable.loading ? (
              <span className="variable-loading">Loading...</span>
            ) : variable.error ? (
              <span className="variable-error">{variable.error}</span>
            ) : (
              <select
                className="variable-select"
                value={variable.selected}
                onChange={(e) => onValueChange(variable.name, e.target.value)}
                disabled={isRepeat}
              >
                {isRepeat ? (
                  <option value={variable.selected}>All ({variable.values.length})</option>
                ) : (
                  variable.values.map((value) => (
                    <option key={value} value={value}>
                      {value}
                    </option>
                  ))
                )}
              </select>
            )}
          </div>
        );
      })}
    </div>
  );
}

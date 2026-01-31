const REFRESH_OPTIONS = [
  { label: 'Off', value: 0 },
  { label: '10s', value: 10_000 },
  { label: '30s', value: 30_000 },
  { label: '1m', value: 60_000 },
  { label: '5m', value: 300_000 },
];

interface RefreshIntervalSelectorProps {
  value: number;
  onChange: (interval: number) => void;
}

export function RefreshIntervalSelector({ value, onChange }: RefreshIntervalSelectorProps) {
  return (
    <select
      className="refresh-interval-selector"
      value={value}
      onChange={(e) => onChange(Number(e.target.value))}
    >
      {REFRESH_OPTIONS.map((opt) => (
        <option key={opt.value} value={opt.value}>
          {opt.value === 0 ? opt.label : `⟳ ${opt.label}`}
        </option>
      ))}
    </select>
  );
}

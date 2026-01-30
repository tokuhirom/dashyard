import { TimeRangeSelector } from './TimeRangeSelector';
import type { TimeRange } from '../types';

interface HeaderProps {
  timeRange: TimeRange;
  onTimeRangeChange: (range_: TimeRange) => void;
  siteTitle: string;
  headerColor: string;
  columns: number;
  onColumnsChange: (n: number) => void;
}

export function Header({ timeRange, onTimeRangeChange, siteTitle, headerColor, columns, onColumnsChange }: HeaderProps) {
  return (
    <header className="header" style={headerColor ? { background: headerColor } : undefined}>
      <h1 className="header-title">{siteTitle}</h1>
      <div className="header-controls">
        <select
          className="column-selector"
          value={columns}
          onChange={(e) => onColumnsChange(Number(e.target.value))}
        >
          {[2, 3, 4, 5, 6].map((n) => (
            <option key={n} value={n}>{n} columns</option>
          ))}
        </select>
        <TimeRangeSelector selected={timeRange} onChange={onTimeRangeChange} />
      </div>
    </header>
  );
}

import { TimeRangeSelector } from './TimeRangeSelector';
import type { TimeRange } from '../types';

interface HeaderProps {
  timeRange: TimeRange;
  onTimeRangeChange: (range_: TimeRange) => void;
}

export function Header({ timeRange, onTimeRangeChange }: HeaderProps) {
  return (
    <header className="header">
      <h1 className="header-title">Dashyard</h1>
      <TimeRangeSelector selected={timeRange} onChange={onTimeRangeChange} />
    </header>
  );
}

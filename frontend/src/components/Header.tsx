import { TimeRangeSelector } from './TimeRangeSelector';
import type { TimeRange } from '../types';

interface HeaderProps {
  timeRange: TimeRange;
  onTimeRangeChange: (range_: TimeRange) => void;
  siteTitle: string;
}

export function Header({ timeRange, onTimeRangeChange, siteTitle }: HeaderProps) {
  return (
    <header className="header">
      <h1 className="header-title">{siteTitle}</h1>
      <TimeRangeSelector selected={timeRange} onChange={onTimeRangeChange} />
    </header>
  );
}

import { TimeRangeSelector } from './TimeRangeSelector';
import { RefreshIntervalSelector } from './RefreshIntervalSelector';
import type { TimeRange } from '../types';

interface HeaderProps {
  timeRange: TimeRange;
  onTimeRangeChange: (range_: TimeRange) => void;
  siteTitle: string;
  headerColor: string;
  refreshInterval: number;
  onRefreshIntervalChange: (interval: number) => void;
}

export function Header({ timeRange, onTimeRangeChange, siteTitle, headerColor, refreshInterval, onRefreshIntervalChange }: HeaderProps) {
  return (
    <header className="header" style={headerColor ? { background: headerColor } : undefined}>
      <h1 className="header-title">{siteTitle}</h1>
      <div className="header-controls">
        <RefreshIntervalSelector value={refreshInterval} onChange={onRefreshIntervalChange} />
        <TimeRangeSelector selected={timeRange} onChange={onTimeRangeChange} />
      </div>
    </header>
  );
}

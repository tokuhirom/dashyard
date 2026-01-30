import { TIME_RANGES } from '../utils/time';
import type { TimeRange } from '../types';

interface TimeRangeSelectorProps {
  selected: TimeRange;
  onChange: (range_: TimeRange) => void;
}

export function TimeRangeSelector({ selected, onChange }: TimeRangeSelectorProps) {
  return (
    <select
      className="time-range-selector"
      value={selected.value}
      onChange={(e) => {
        const range_ = TIME_RANGES.find((r) => r.value === e.target.value);
        if (range_) onChange(range_);
      }}
    >
      {TIME_RANGES.map((range_) => (
        <option key={range_.value} value={range_.value}>
          {range_.label}
        </option>
      ))}
    </select>
  );
}

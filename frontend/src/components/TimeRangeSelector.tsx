import { useState } from 'react';
import { TIME_RANGES, computeStep } from '../utils/time';
import type { TimeRange, AbsoluteTimeRange } from '../types';

interface TimeRangeSelectorProps {
  selected: TimeRange;
  onChange: (range_: TimeRange) => void;
}

function toLocalDatetimeString(unixSeconds: number): string {
  const d = new Date(unixSeconds * 1000);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function fromLocalDatetimeString(s: string): number {
  return Math.floor(new Date(s).getTime() / 1000);
}

export function TimeRangeSelector({ selected, onChange }: TimeRangeSelectorProps) {
  const [showCustom, setShowCustom] = useState(selected.type === 'absolute');
  const [fromValue, setFromValue] = useState(() => {
    if (selected.type === 'absolute') {
      return toLocalDatetimeString(selected.start);
    }
    const now = Math.floor(Date.now() / 1000);
    return toLocalDatetimeString(now - 3600);
  });
  const [toValue, setToValue] = useState(() => {
    if (selected.type === 'absolute') {
      return toLocalDatetimeString(selected.end);
    }
    return toLocalDatetimeString(Math.floor(Date.now() / 1000));
  });

  const selectValue = selected.type === 'relative' ? selected.value : '__custom__';

  const handleSelectChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const val = e.target.value;
    if (val === '__custom__') {
      setShowCustom(true);
      return;
    }
    setShowCustom(false);
    const range_ = TIME_RANGES.find((r) => r.value === val);
    if (range_) onChange(range_);
  };

  const handleApply = () => {
    const startUnix = fromLocalDatetimeString(fromValue);
    const endUnix = fromLocalDatetimeString(toValue);
    if (isNaN(startUnix) || isNaN(endUnix) || startUnix >= endUnix) return;
    const duration = endUnix - startUnix;
    const step = computeStep(duration);
    const fromDate = new Date(startUnix * 1000);
    const toDate = new Date(endUnix * 1000);
    const fmt = (d: Date) => `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
    const range_: AbsoluteTimeRange = {
      type: 'absolute',
      label: `${fmt(fromDate)} – ${fmt(toDate)}`,
      start: startUnix,
      end: endUnix,
      step,
    };
    onChange(range_);
  };

  return (
    <div className="time-range-selector-container">
      <select
        className="time-range-selector"
        value={selectValue}
        onChange={handleSelectChange}
      >
        {TIME_RANGES.map((range_) => (
          <option key={range_.value} value={range_.value}>
            {range_.label}
          </option>
        ))}
        <option value="__custom__">Custom...</option>
      </select>
      {showCustom && (
        <div className="time-range-custom">
          <input
            type="datetime-local"
            className="time-range-input"
            value={fromValue}
            onChange={(e) => setFromValue(e.target.value)}
          />
          <span className="time-range-separator">to</span>
          <input
            type="datetime-local"
            className="time-range-input"
            value={toValue}
            onChange={(e) => setToValue(e.target.value)}
          />
          <button className="time-range-apply" onClick={handleApply}>
            Apply
          </button>
        </div>
      )}
    </div>
  );
}

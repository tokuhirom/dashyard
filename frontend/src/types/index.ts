export interface Threshold {
  value: number;
  color?: string;
  label?: string;
}

export interface Panel {
  title: string;
  type: 'graph' | 'markdown';
  chart_type?: 'line' | 'bar' | 'area' | 'scatter';
  query?: string;
  unit?: 'bytes' | 'percent' | 'count';
  y_min?: number;
  y_max?: number;
  legend?: string;
  thresholds?: Threshold[];
  stacked?: boolean;
  content?: string;
}

export interface Row {
  title: string;
  repeat?: string;
  panels: Panel[];
}

export interface Variable {
  name: string;
  label?: string;
  query: string;
}

export interface Dashboard {
  title: string;
  variables?: Variable[];
  rows: Row[];
  path: string;
}

export interface LabelValuesResponse {
  status: string;
  data: string[];
}

export interface DashboardListItem {
  path: string;
  title: string;
}

export interface DashboardTreeNode {
  name: string;
  path?: string;
  children?: DashboardTreeNode[];
}

export interface DashboardsResponse {
  dashboards: DashboardListItem[];
  tree: DashboardTreeNode[];
  site_title: string;
  header_color: string;
}

export interface PrometheusResult {
  metric: Record<string, string>;
  values: [number, string][];
}

export interface PrometheusResponse {
  status: string;
  data: {
    resultType: string;
    result: PrometheusResult[];
  };
}

export interface RelativeTimeRange {
  type: 'relative';
  label: string;
  value: string;
  duration: number; // seconds
  step: string;
}

export interface AbsoluteTimeRange {
  type: 'absolute';
  label: string;
  start: number; // Unix seconds
  end: number;   // Unix seconds
  step: string;
}

export type TimeRange = RelativeTimeRange | AbsoluteTimeRange;

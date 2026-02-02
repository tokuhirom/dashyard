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
  datasource?: string;
  unit?: 'bytes' | 'percent' | 'count' | 'seconds';
  y_min?: number;
  y_max?: number;
  legend?: string;
  legend_display?: boolean;
  legend_position?: 'top' | 'bottom' | 'left' | 'right';
  legend_align?: 'start' | 'center' | 'end';
  legend_max_height?: number;
  legend_max_width?: number;
  thresholds?: Threshold[];
  stacked?: boolean;
  y_scale?: 'linear' | 'log';
  content?: string;
  span?: number;
}

export interface Row {
  title: string;
  repeat?: string;
  panels: Panel[];
}

export interface Variable {
  name: string;
  type?: 'query' | 'datasource';
  label?: string;
  query?: string;
  datasource?: string;
  hide?: boolean;
}

export interface DatasourcesResponse {
  datasources: string[];
  default: string;
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

export interface QueryResult {
  metric: Record<string, string>;
  values: [number, string][];
}

export interface QueryResponse {
  status: string;
  data: {
    resultType: string;
    result: QueryResult[];
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

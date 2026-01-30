export interface Panel {
  title: string;
  type: 'graph' | 'markdown';
  query?: string;
  unit?: 'bytes' | 'percent' | 'count';
  legend?: string;
  content?: string;
}

export interface Row {
  title: string;
  panels: Panel[];
}

export interface Dashboard {
  title: string;
  rows: Row[];
  path: string;
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

export interface TimeRange {
  label: string;
  value: string;
  duration: number; // seconds
  step: string;
}

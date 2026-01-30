import { Header } from './Header';
import { Sidebar } from './Sidebar';
import type { DashboardTreeNode, TimeRange } from '../types';

interface LayoutProps {
  tree: DashboardTreeNode[];
  currentPath: string;
  timeRange: TimeRange;
  onTimeRangeChange: (range_: TimeRange) => void;
  onNavigate: (path: string) => void;
  siteTitle: string;
  headerColor: string;
  children: React.ReactNode;
}

export function Layout({ tree, currentPath, timeRange, onTimeRangeChange, onNavigate, siteTitle, headerColor, children }: LayoutProps) {
  return (
    <div className="layout">
      <Header timeRange={timeRange} onTimeRangeChange={onTimeRangeChange} siteTitle={siteTitle} headerColor={headerColor} />
      <div className="layout-body">
        <Sidebar tree={tree} currentPath={currentPath} onNavigate={onNavigate} />
        <main className="layout-main">
          {children}
        </main>
      </div>
    </div>
  );
}

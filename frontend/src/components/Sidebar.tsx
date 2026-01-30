import { useState } from 'react';
import type { DashboardTreeNode } from '../types';

interface SidebarProps {
  tree: DashboardTreeNode[];
  currentPath: string;
  onNavigate: (path: string) => void;
}

export function Sidebar({ tree, currentPath, onNavigate }: SidebarProps) {
  return (
    <nav className="sidebar">
      <div className="sidebar-content">
        {tree.map((node) => (
          <TreeNode
            key={node.name}
            node={node}
            currentPath={currentPath}
            onNavigate={onNavigate}
          />
        ))}
      </div>
    </nav>
  );
}

interface TreeNodeProps {
  node: DashboardTreeNode;
  currentPath: string;
  onNavigate: (path: string) => void;
  depth?: number;
}

function TreeNode({ node, currentPath, onNavigate, depth = 0 }: TreeNodeProps) {
  const [expanded, setExpanded] = useState(true);
  const isLeaf = !!node.path;
  const isActive = node.path === currentPath;

  if (isLeaf) {
    return (
      <a
        href={`/d/${node.path}`}
        className={`sidebar-item ${isActive ? 'active' : ''}`}
        style={{ paddingLeft: `${(depth + 1) * 12}px` }}
        onClick={(e) => {
          e.preventDefault();
          onNavigate(node.path!);
        }}
      >
        {node.name}
      </a>
    );
  }

  return (
    <div className="sidebar-group">
      <div
        className="sidebar-group-header"
        style={{ paddingLeft: `${depth * 12}px` }}
        onClick={() => setExpanded(!expanded)}
      >
        <span className={`sidebar-arrow ${expanded ? 'expanded' : ''}`}>&#9656;</span>
        {node.name}
      </div>
      {expanded && node.children?.map((child) => (
        <TreeNode
          key={child.name}
          node={child}
          currentPath={currentPath}
          onNavigate={onNavigate}
          depth={depth + 1}
        />
      ))}
    </div>
  );
}

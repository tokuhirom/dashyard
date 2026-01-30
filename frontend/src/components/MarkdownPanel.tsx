import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

interface MarkdownPanelProps {
  title: string;
  content: string;
}

export function MarkdownPanel({ title, content }: MarkdownPanelProps) {
  return (
    <div className="panel markdown-panel">
      <h3 className="panel-title">{title}</h3>
      <div className="panel-content">
        <ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown>
      </div>
    </div>
  );
}

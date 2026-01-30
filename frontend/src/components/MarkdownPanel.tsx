import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

interface MarkdownPanelProps {
  title: string;
  content: string;
  id?: string;
}

export function MarkdownPanel({ title, content, id }: MarkdownPanelProps) {
  return (
    <div className="panel markdown-panel" id={id}>
      <h3 className="panel-title">
        {title}
        {id && <a href={`#${id}`} className="panel-anchor">#</a>}
      </h3>
      <div className="panel-content">
        <ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown>
      </div>
    </div>
  );
}

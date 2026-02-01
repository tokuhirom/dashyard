import { Component } from 'react';
import type { ErrorInfo, ReactNode } from 'react';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, info: ErrorInfo): void {
    console.error('Uncaught error:', error, info.componentStack);
  }

  render(): ReactNode {
    if (this.state.hasError) {
      return (
        <div style={{ padding: '2rem', textAlign: 'center' }}>
          <h1>Something went wrong</h1>
          <p style={{ color: '#666', marginBottom: '1rem' }}>
            An unexpected error occurred. Try reloading the page.
          </p>
          <pre
            style={{
              textAlign: 'left',
              background: '#f5f5f5',
              padding: '1rem',
              borderRadius: '4px',
              maxWidth: '600px',
              margin: '0 auto',
              overflow: 'auto',
              fontSize: '0.85rem',
            }}
          >
            {this.state.error?.message}
          </pre>
          <button
            onClick={() => window.location.reload()}
            style={{
              marginTop: '1rem',
              padding: '0.5rem 1rem',
              cursor: 'pointer',
            }}
          >
            Reload
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}

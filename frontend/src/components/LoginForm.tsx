import { useState, useEffect } from 'react';
import { login, fetchAuthInfo } from '../api/client';
import type { AuthInfo } from '../api/client';

interface LoginFormProps {
  onLoginSuccess: () => void;
}

export function LoginForm({ onLoginSuccess }: LoginFormProps) {
  const [userId, setUserId] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [authInfo, setAuthInfo] = useState<AuthInfo | null>(null);

  useEffect(() => {
    fetchAuthInfo()
      .then(setAuthInfo)
      .catch(() => {
        // Fallback: assume password-only if auth-info fails
        setAuthInfo({ password_enabled: true, oauth_providers: [] });
      });
  }, []);

  // Check for OAuth error in URL
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const oauthError = params.get('error');
    if (oauthError) {
      const messages: Record<string, string> = {
        oauth_failed: 'OAuth authentication failed',
        access_denied: 'Access denied. You are not in the allowed users list.',
        unknown_provider: 'Unknown OAuth provider',
        session_failed: 'Session creation failed',
      };
      setError(messages[oauthError] || 'Authentication failed');
      // Clean up URL
      window.history.replaceState({}, '', '/');
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(userId, password);
      onLoginSuccess();
    } catch {
      setError('Invalid credentials');
    } finally {
      setLoading(false);
    }
  };

  if (!authInfo) {
    return (
      <div className="login-container">
        <div className="login-form">
          <h1>Dashyard</h1>
          <div style={{ textAlign: 'center', color: '#6b7280' }}>Loading...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="login-container">
      <div className="login-form">
        <h1>Dashyard</h1>
        {error && <div className="login-error">{error}</div>}

        {authInfo.oauth_providers.length > 0 && (
          <div className="oauth-buttons">
            {authInfo.oauth_providers.map((provider) => (
              <a
                key={provider.name}
                href={provider.url}
                className={`oauth-button oauth-button-${provider.name}`}
              >
                Sign in with {provider.name.charAt(0).toUpperCase() + provider.name.slice(1)}
              </a>
            ))}
          </div>
        )}

        {authInfo.password_enabled && authInfo.oauth_providers.length > 0 && (
          <div className="login-divider">
            <span>or</span>
          </div>
        )}

        {authInfo.password_enabled && (
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="userId">User ID</label>
              <input
                id="userId"
                type="text"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                required
                autoFocus
              />
            </div>
            <div className="form-group">
              <label htmlFor="password">Password</label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <button type="submit" disabled={loading}>
              {loading ? 'Logging in...' : 'Log in'}
            </button>
          </form>
        )}
      </div>
    </div>
  );
}

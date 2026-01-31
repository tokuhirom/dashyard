import { useState, useEffect } from 'react';
import { login, fetchAuthInfo } from '../api/client';
import type { AuthInfo } from '../api/client';

interface LoginFormProps {
  onLoginSuccess: () => void;
}

const providerLabels: Record<string, string> = {
  github: 'GitHub',
  google: 'Google',
  oidc: 'SSO',
};

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
        setAuthInfo({ password_enabled: true, oauth_enabled: false });
      });
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
          <div className="login-loading">Loading...</div>
        </div>
      </div>
    );
  }

  const providerName = authInfo.oauth_provider
    ? providerLabels[authInfo.oauth_provider] || authInfo.oauth_provider
    : '';

  return (
    <div className="login-container">
      <div className="login-form">
        <h1>Dashyard</h1>
        {error && <div className="login-error">{error}</div>}

        {authInfo.oauth_enabled && authInfo.oauth_login_url && (
          <a href={authInfo.oauth_login_url} className="oauth-login-btn">
            Log in with {providerName}
          </a>
        )}

        {authInfo.oauth_enabled && authInfo.password_enabled && (
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

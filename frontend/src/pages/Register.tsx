import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { RegisterCredentials } from '../types';

interface RegisterProps {
  onRegister: (credentials: RegisterCredentials) => Promise<void>;
  loading?: boolean;
}

const Register: React.FC<RegisterProps> = ({ onRegister, loading = false }) => {
  const [credentials, setCredentials] = useState<RegisterCredentials>({
    username: '',
    email: '',
    password: '',
    displayName: '',
  });
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCredentials(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    
    try {
      await onRegister(credentials);
      navigate('/');
    } catch (err: any) {
      setError(err.response?.data?.error || '登録に失敗しました');
    }
  };

  return (
    <div style={containerStyle}>
      <div style={formContainerStyle}>
        <h1 style={titleStyle}>Digeon</h1>
        <h2 style={subtitleStyle}>新規登録</h2>
        
        {error && (
          <div style={errorStyle}>
            {error}
          </div>
        )}
        
        <form onSubmit={handleSubmit} style={formStyle}>
          <div style={fieldStyle}>
            <label style={labelStyle}>ユーザー名</label>
            <input
              type="text"
              name="username"
              value={credentials.username}
              onChange={handleChange}
              style={inputStyle}
              required
              placeholder="username"
            />
          </div>
          
          <div style={fieldStyle}>
            <label style={labelStyle}>表示名</label>
            <input
              type="text"
              name="displayName"
              value={credentials.displayName}
              onChange={handleChange}
              style={inputStyle}
              required
              placeholder="表示名"
            />
          </div>
          
          <div style={fieldStyle}>
            <label style={labelStyle}>メールアドレス</label>
            <input
              type="email"
              name="email"
              value={credentials.email}
              onChange={handleChange}
              style={inputStyle}
              required
              placeholder="email@example.com"
            />
          </div>
          
          <div style={fieldStyle}>
            <label style={labelStyle}>パスワード</label>
            <input
              type="password"
              name="password"
              value={credentials.password}
              onChange={handleChange}
              style={inputStyle}
              required
              placeholder="8文字以上"
            />
          </div>
          
          <button 
            type="submit" 
            disabled={loading}
            style={{
              ...buttonStyle,
              opacity: loading ? 0.6 : 1,
            }}
          >
            {loading ? '登録中...' : '新規登録'}
          </button>
        </form>
        
        <p style={linkTextStyle}>
          すでにアカウントをお持ちの方は{' '}
          <Link to="/login" style={linkStyle}>
            ログイン
          </Link>
        </p>
      </div>
    </div>
  );
};

const containerStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  minHeight: '100vh',
  backgroundColor: '#f5f8fa',
};

const formContainerStyle: React.CSSProperties = {
  backgroundColor: 'white',
  padding: '2rem',
  borderRadius: '8px',
  boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
  width: '100%',
  maxWidth: '400px',
};

const titleStyle: React.CSSProperties = {
  fontSize: '2rem',
  fontWeight: 'bold',
  color: '#1da1f2',
  textAlign: 'center',
  marginBottom: '0.5rem',
};

const subtitleStyle: React.CSSProperties = {
  fontSize: '1.5rem',
  fontWeight: 'bold',
  color: '#14171a',
  textAlign: 'center',
  marginBottom: '1.5rem',
};

const errorStyle: React.CSSProperties = {
  backgroundColor: '#ffebee',
  color: '#c62828',
  padding: '0.75rem',
  borderRadius: '4px',
  marginBottom: '1rem',
  border: '1px solid #e57373',
};

const formStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '1rem',
};

const fieldStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
};

const labelStyle: React.CSSProperties = {
  fontSize: '0.9rem',
  fontWeight: 'bold',
  color: '#14171a',
  marginBottom: '0.5rem',
};

const inputStyle: React.CSSProperties = {
  padding: '0.75rem',
  border: '1px solid #e1e8ed',
  borderRadius: '4px',
  fontSize: '1rem',
  outline: 'none',
  transition: 'border-color 0.2s',
};

const buttonStyle: React.CSSProperties = {
  backgroundColor: '#1da1f2',
  color: 'white',
  border: 'none',
  borderRadius: '20px',
  padding: '0.75rem 1.5rem',
  fontSize: '1rem',
  fontWeight: 'bold',
  cursor: 'pointer',
  marginTop: '1rem',
};

const linkTextStyle: React.CSSProperties = {
  textAlign: 'center',
  marginTop: '1.5rem',
  color: '#657786',
};

const linkStyle: React.CSSProperties = {
  color: '#1da1f2',
  textDecoration: 'none',
  fontWeight: 'bold',
};

export default Register;
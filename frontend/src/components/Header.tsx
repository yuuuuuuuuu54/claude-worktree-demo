import React from 'react';
import { Link } from 'react-router-dom';
import { User } from '../types';

interface HeaderProps {
  user: User | null;
  onLogout: () => void;
}

const Header: React.FC<HeaderProps> = ({ user, onLogout }) => {
  return (
    <header style={headerStyle}>
      <div style={containerStyle}>
        <Link to="/" style={logoStyle}>
          <h1>Digeon</h1>
        </Link>
        
        <nav style={navStyle}>
          {user ? (
            <div style={userMenuStyle}>
              <span style={usernameStyle}>@{user.username}</span>
              <button onClick={onLogout} style={logoutButtonStyle}>
                ログアウト
              </button>
            </div>
          ) : (
            <div style={authLinksStyle}>
              <Link to="/login" style={linkStyle}>
                ログイン
              </Link>
              <Link to="/register" style={registerLinkStyle}>
                新規登録
              </Link>
            </div>
          )}
        </nav>
      </div>
    </header>
  );
};

const headerStyle: React.CSSProperties = {
  backgroundColor: '#1da1f2',
  color: 'white',
  padding: '1rem 0',
  borderBottom: '1px solid #e1e8ed',
};

const containerStyle: React.CSSProperties = {
  maxWidth: '1200px',
  margin: '0 auto',
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  padding: '0 1rem',
};

const logoStyle: React.CSSProperties = {
  textDecoration: 'none',
  color: 'white',
};

const navStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
};

const userMenuStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '1rem',
};

const usernameStyle: React.CSSProperties = {
  fontWeight: 'bold',
};

const logoutButtonStyle: React.CSSProperties = {
  backgroundColor: 'transparent',
  border: '1px solid white',
  color: 'white',
  padding: '0.5rem 1rem',
  borderRadius: '20px',
  cursor: 'pointer',
};

const authLinksStyle: React.CSSProperties = {
  display: 'flex',
  gap: '1rem',
};

const linkStyle: React.CSSProperties = {
  color: 'white',
  textDecoration: 'none',
  padding: '0.5rem 1rem',
  borderRadius: '20px',
  border: '1px solid white',
};

const registerLinkStyle: React.CSSProperties = {
  color: '#1da1f2',
  textDecoration: 'none',
  padding: '0.5rem 1rem',
  borderRadius: '20px',
  backgroundColor: 'white',
  fontWeight: 'bold',
};

export default Header;
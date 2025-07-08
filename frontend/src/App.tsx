import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { User, LoginCredentials, RegisterCredentials } from './types';
import { apiClient } from './api/client';
import Header from './components/Header';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import './App.css';

function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      getCurrentUser();
    } else {
      setLoading(false);
    }
  }, []);

  const getCurrentUser = async () => {
    try {
      const response = await apiClient.getCurrentUser();
      setUser(response.data);
    } catch (error) {
      console.error('Failed to get current user:', error);
      localStorage.removeItem('token');
    } finally {
      setLoading(false);
    }
  };

  const handleLogin = async (credentials: LoginCredentials) => {
    const response = await apiClient.login(credentials);
    localStorage.setItem('token', response.data.token);
    setUser(response.data.user);
  };

  const handleRegister = async (credentials: RegisterCredentials) => {
    const response = await apiClient.register(credentials);
    localStorage.setItem('token', response.data.token);
    setUser(response.data.user);
  };

  const handleLogout = async () => {
    try {
      await apiClient.logout();
    } catch (error) {
      console.error('Failed to logout:', error);
    } finally {
      localStorage.removeItem('token');
      setUser(null);
    }
  };

  if (loading) {
    return (
      <div style={loadingStyle}>
        <p>読み込み中...</p>
      </div>
    );
  }

  return (
    <Router>
      <div className="App">
        <Header user={user} onLogout={handleLogout} />
        <main>
          <Routes>
            <Route path="/" element={<Home user={user} />} />
            <Route 
              path="/login" 
              element={
                user ? <Navigate to="/" /> : <Login onLogin={handleLogin} />
              } 
            />
            <Route 
              path="/register" 
              element={
                user ? <Navigate to="/" /> : <Register onRegister={handleRegister} />
              } 
            />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

const loadingStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  height: '100vh',
  fontSize: '1.2rem',
  color: '#657786',
};

export default App;

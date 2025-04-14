import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Navbar = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav style={styles.navbar}>
      <div style={styles.logo}>
        <Link to="/" style={styles.link}>Status Page</Link>
      </div>
      <div style={styles.links}>
        {isAuthenticated ? (
          <>
            <Link to="/dashboard" style={styles.link}>Dashboard</Link>
            <Link to="/services" style={styles.link}>Services</Link>
            <Link to="/incidents" style={styles.link}>Incidents</Link>
            <Link to={`/public/${user.orgId}`} style={styles.link} target="_blank">Public Page</Link>
            <span style={styles.user}>{user.email}</span>
            <button onClick={handleLogout} style={styles.logoutButton}>Logout</button>
          </>
        ) : (
          <>
            <Link to="/login" style={styles.link}>Login</Link>
            <Link to="/signup" style={styles.link}>Signup</Link>
          </>
        )}
      </div>
    </nav>
  );
};

const styles = {
  navbar: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '1rem 2rem',
    backgroundColor: '#f8f9fa',
    borderBottom: '1px solid #ddd',
  },
  logo: {
    fontWeight: 'bold',
    fontSize: '1.5rem',
  },
  links: {
    display: 'flex',
    gap: '1.5rem',
    alignItems: 'center',
  },
  link: {
    textDecoration: 'none',
    color: '#333',
  },
  user: {
    marginLeft: '1rem',
    color: '#666',
  },
  logoutButton: {
    backgroundColor: '#f0f0f0',
    border: '1px solid #ddd',
    padding: '0.25rem 0.75rem',
    borderRadius: '4px',
    cursor: 'pointer',
  },
};

export default Navbar; 
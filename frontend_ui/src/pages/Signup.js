import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { signup } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

const Signup = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [orgName, setOrgName] = useState('');
  const [orgId, setOrgId] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login: authLogin } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await signup(email, password, orgName, orgId);
      authLogin(response.data.user, response.data.token);
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create account');
    } finally {
      setLoading(false);
    }
  };

  // Generate a slug from the organization name
  const handleOrgNameChange = (e) => {
    const name = e.target.value;
    setOrgName(name);
    if (name) {
      setOrgId(name.toLowerCase().replace(/[^a-z0-9]/g, '-'));
    } else {
      setOrgId('');
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.formContainer}>
        <h2 style={styles.title}>Sign Up</h2>
        {error && <div style={styles.error}>{error}</div>}
        <form onSubmit={handleSubmit} style={styles.form}>
          <div style={styles.formGroup}>
            <label htmlFor="email" style={styles.label}>Email</label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              style={styles.input}
            />
          </div>
          <div style={styles.formGroup}>
            <label htmlFor="password" style={styles.label}>Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={6}
              style={styles.input}
            />
          </div>
          <div style={styles.formGroup}>
            <label htmlFor="orgName" style={styles.label}>Organization Name</label>
            <input
              id="orgName"
              type="text"
              value={orgName}
              onChange={handleOrgNameChange}
              required
              style={styles.input}
            />
          </div>
          <div style={styles.formGroup}>
            <label htmlFor="orgId" style={styles.label}>Organization ID</label>
            <input
              id="orgId"
              type="text"
              value={orgId}
              onChange={(e) => setOrgId(e.target.value)}
              required
              pattern="[a-z0-9-]+"
              title="Only lowercase letters, numbers, and hyphens"
              style={styles.input}
            />
            <small style={styles.hint}>
              Used for your public status page URL. Only lowercase letters, numbers, and hyphens.
            </small>
          </div>
          <button type="submit" disabled={loading} style={styles.button}>
            {loading ? 'Creating account...' : 'Sign Up'}
          </button>
        </form>
        <div style={styles.footer}>
          Already have an account? <Link to="/login" style={styles.link}>Log in</Link>
        </div>
      </div>
    </div>
  );
};

const styles = {
  container: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: 'calc(100vh - 70px)',
    padding: '2rem',
  },
  formContainer: {
    width: '100%',
    maxWidth: '400px',
    padding: '2rem',
    backgroundColor: '#fff',
    borderRadius: '8px',
    boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
  },
  title: {
    marginBottom: '1.5rem',
    textAlign: 'center',
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: '1rem',
  },
  formGroup: {
    display: 'flex',
    flexDirection: 'column',
    gap: '0.5rem',
  },
  label: {
    fontWeight: '500',
  },
  input: {
    padding: '0.75rem',
    borderRadius: '4px',
    border: '1px solid #ddd',
    fontSize: '1rem',
  },
  button: {
    padding: '0.75rem',
    backgroundColor: '#4A90E2',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    fontSize: '1rem',
    cursor: 'pointer',
    marginTop: '1rem',
  },
  error: {
    color: 'red',
    marginBottom: '1rem',
    textAlign: 'center',
  },
  hint: {
    color: '#666',
    fontSize: '0.875rem',
  },
  footer: {
    marginTop: '1.5rem',
    textAlign: 'center',
  },
  link: {
    color: '#4A90E2',
    textDecoration: 'none',
  },
};

export default Signup; 
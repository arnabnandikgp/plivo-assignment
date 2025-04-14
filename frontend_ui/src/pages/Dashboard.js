import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getServices, getIncidents } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import { useWebSocket } from '../contexts/WebSocketContext';

const Dashboard = () => {
  const [services, setServices] = useState([]);
  const [incidents, setIncidents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user } = useAuth();
  const { connected } = useWebSocket();

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [servicesResponse, incidentsResponse] = await Promise.all([
          getServices(),
          getIncidents(),
        ]);
        setServices(servicesResponse.data.services || []);
        setIncidents(incidentsResponse.data.incidents || []);
      } catch (err) {
        setError('Failed to fetch data');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <div style={styles.loading}>Loading...</div>;
  }

  if (error) {
    return <div style={styles.error}>{error}</div>;
  }

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.title}>Dashboard</h1>
        <div style={styles.status}>
          <span>
            WebSocket: {connected ? (
              <span style={styles.connected}>Connected</span>
            ) : (
              <span style={styles.disconnected}>Disconnected</span>
            )}
          </span>
        </div>
      </div>

      <div style={styles.section}>
        <div style={styles.sectionHeader}>
          <h2 style={styles.sectionTitle}>Services</h2>
          <Link to="/services/new" style={styles.button}>
            Add Service
          </Link>
        </div>
        {services.length === 0 ? (
          <p style={styles.empty}>No services yet. Add your first service to start.</p>
        ) : (
          <div style={styles.grid}>
            {services.map((service) => (
              <div key={service.id} style={styles.card}>
                <h3 style={styles.cardTitle}>{service.name}</h3>
                <div style={getStatusStyle(service.status)}>{service.status}</div>
                <Link to={`/services/${service.id}`} style={styles.link}>
                  Manage
                </Link>
              </div>
            ))}
          </div>
        )}
      </div>

      <div style={styles.section}>
        <div style={styles.sectionHeader}>
          <h2 style={styles.sectionTitle}>Active Incidents</h2>
          <Link to="/incidents/new" style={styles.button}>
            Create Incident
          </Link>
        </div>
        {incidents.length === 0 ? (
          <p style={styles.empty}>No active incidents.</p>
        ) : (
          <div style={styles.incidentsList}>
            {incidents.map((incident) => (
              <div key={incident.incident.id} style={styles.incidentCard}>
                <div style={styles.incidentHeader}>
                  <h3 style={styles.incidentTitle}>{incident.incident.title}</h3>
                  <div style={getIncidentStatusStyle(incident.incident.status)}>
                    {incident.incident.status}
                  </div>
                </div>
                <p style={styles.incidentDescription}>{incident.incident.description}</p>
                <div style={styles.incidentMeta}>
                  <div style={styles.affectedServices}>
                    Services: {incident.services.map(s => s.name).join(', ')}
                  </div>
                  <div style={styles.updatesCount}>
                    Updates: {incident.updates.length}
                  </div>
                </div>
                <Link to={`/incidents/${incident.incident.id}`} style={styles.link}>
                  Manage
                </Link>
              </div>
            ))}
          </div>
        )}
      </div>

      <div style={styles.section}>
        <div style={styles.sectionHeader}>
          <h2 style={styles.sectionTitle}>Status Page</h2>
        </div>
        <p style={styles.statusPageInfo}>
          Your public status page is available at:{' '}
          <Link to={`/public/${user.orgId}`} target="_blank" style={styles.link}>
            {window.location.origin}/public/{user.orgId}
          </Link>
        </p>
      </div>
    </div>
  );
};

const getStatusStyle = (status) => {
  const baseStyle = {
    padding: '0.25rem 0.5rem',
    borderRadius: '4px',
    display: 'inline-block',
    marginBottom: '0.5rem',
  };

  switch (status) {
    case 'Operational':
      return { ...baseStyle, backgroundColor: '#4caf50', color: 'white' };
    case 'Degraded':
      return { ...baseStyle, backgroundColor: '#ff9800', color: 'white' };
    case 'Outage':
      return { ...baseStyle, backgroundColor: '#f44336', color: 'white' };
    default:
      return { ...baseStyle, backgroundColor: '#9e9e9e', color: 'white' };
  }
};

const getIncidentStatusStyle = (status) => {
  const baseStyle = {
    padding: '0.25rem 0.5rem',
    borderRadius: '4px',
    display: 'inline-block',
  };

  switch (status) {
    case 'Investigating':
      return { ...baseStyle, backgroundColor: '#f44336', color: 'white' };
    case 'Identified':
      return { ...baseStyle, backgroundColor: '#ff9800', color: 'white' };
    case 'Monitoring':
      return { ...baseStyle, backgroundColor: '#2196f3', color: 'white' };
    case 'Resolved':
      return { ...baseStyle, backgroundColor: '#4caf50', color: 'white' };
    default:
      return { ...baseStyle, backgroundColor: '#9e9e9e', color: 'white' };
  }
};

const styles = {
  container: {
    padding: '2rem',
    maxWidth: '1200px',
    margin: '0 auto',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '2rem',
  },
  title: {
    margin: 0,
  },
  status: {
    fontSize: '0.875rem',
  },
  connected: {
    color: 'green',
    fontWeight: 'bold',
  },
  disconnected: {
    color: 'red',
    fontWeight: 'bold',
  },
  section: {
    marginBottom: '3rem',
  },
  sectionHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '1rem',
  },
  sectionTitle: {
    margin: 0,
  },
  button: {
    padding: '0.5rem 1rem',
    backgroundColor: '#4A90E2',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    textDecoration: 'none',
    fontSize: '0.875rem',
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))',
    gap: '1rem',
  },
  card: {
    padding: '1.5rem',
    backgroundColor: '#fff',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.05)',
  },
  cardTitle: {
    marginTop: 0,
    marginBottom: '0.5rem',
  },
  incidentsList: {
    display: 'flex',
    flexDirection: 'column',
    gap: '1rem',
  },
  incidentCard: {
    padding: '1.5rem',
    backgroundColor: '#fff',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.05)',
  },
  incidentHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '0.5rem',
  },
  incidentTitle: {
    margin: 0,
  },
  incidentDescription: {
    marginBottom: '1rem',
    color: '#555',
  },
  incidentMeta: {
    display: 'flex',
    justifyContent: 'space-between',
    fontSize: '0.875rem',
    color: '#666',
    marginBottom: '1rem',
  },
  link: {
    color: '#4A90E2',
    textDecoration: 'none',
  },
  empty: {
    color: '#666',
    textAlign: 'center',
    padding: '2rem',
    backgroundColor: '#f9f9f9',
    borderRadius: '8px',
  },
  loading: {
    textAlign: 'center',
    padding: '3rem',
    fontSize: '1.25rem',
    color: '#666',
  },
  error: {
    textAlign: 'center',
    padding: '3rem',
    fontSize: '1.25rem',
    color: 'red',
  },
  statusPageInfo: {
    padding: '1.5rem',
    backgroundColor: '#f9f9f9',
    borderRadius: '8px',
  },
};

export default Dashboard; 
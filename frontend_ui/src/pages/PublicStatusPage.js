import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getPublicServices, getPublicIncidents, getWebSocketUrl } from '../services/api';

const PublicStatusPage = () => {
  const { orgId } = useParams();
  const [services, setServices] = useState([]);
  const [incidents, setIncidents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [servicesResponse, incidentsResponse] = await Promise.all([
          getPublicServices(orgId),
          getPublicIncidents(orgId),
        ]);
        setServices(servicesResponse.data.services || []);
        setIncidents(incidentsResponse.data.incidents || []);
      } catch (err) {
        setError('Failed to fetch status data');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();

    // Set up WebSocket
    const ws = new WebSocket(getWebSocketUrl(orgId));
    
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        
        // Handle different types of events
        if (data.event === 'SERVICE_UPDATED') {
          setServices(prevServices => {
            const updated = [...prevServices];
            const index = updated.findIndex(s => s.id === data.data.id);
            if (index >= 0) {
              updated[index] = data.data;
            }
            return updated;
          });
        } else if (data.event === 'INCIDENT_CREATED' || data.event === 'INCIDENT_UPDATED') {
          // Refresh incidents
          getPublicIncidents(orgId)
            .then(response => setIncidents(response.data.incidents || []));
        }
      } catch (error) {
        console.error('Error handling WebSocket message', error);
      }
    };

    setSocket(ws);

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [orgId]);

  if (loading) {
    return <div style={styles.loading}>Loading...</div>;
  }

  if (error) {
    return <div style={styles.error}>{error}</div>;
  }

  // Count services by status
  const serviceStats = services.reduce(
    (acc, service) => {
      acc[service.status] = (acc[service.status] || 0) + 1;
      return acc;
    },
    { Operational: 0, Degraded: 0, Outage: 0 }
  );

  // Calculate overall system status
  const getSystemStatus = () => {
    if (serviceStats.Outage > 0) return 'Major Outage';
    if (serviceStats.Degraded > 0) return 'Partial Outage';
    return 'All Systems Operational';
  };

  const systemStatus = getSystemStatus();

  return (
    <div style={styles.container}>
      <header style={styles.header}>
        <h1 style={styles.title}>System Status</h1>
        <div style={getSystemStatusStyle(systemStatus)}>
          {systemStatus}
        </div>
      </header>

      <section style={styles.section}>
        <h2 style={styles.sectionTitle}>Services</h2>
        <div style={styles.servicesList}>
          {services.map((service) => (
            <div key={service.id} style={styles.serviceItem}>
              <div style={styles.serviceName}>{service.name}</div>
              <div style={getStatusStyle(service.status)}>{service.status}</div>
            </div>
          ))}
        </div>
      </section>

      {incidents.length > 0 && (
        <section style={styles.section}>
          <h2 style={styles.sectionTitle}>Active Incidents</h2>
          <div style={styles.incidentsList}>
            {incidents.map((incident) => (
              <div key={incident.id} style={styles.incidentCard}>
                <div style={styles.incidentHeader}>
                  <h3 style={styles.incidentTitle}>{incident.title}</h3>
                  <div style={getIncidentStatusStyle(incident.status)}>
                    {incident.status}
                  </div>
                </div>
                {incident.description && (
                  <p style={styles.incidentDescription}>{incident.description}</p>
                )}
                <div style={styles.affectedServices}>
                  Affected services: {incident.services.map(s => s.name).join(', ')}
                </div>
                <div style={styles.updates}>
                  <h4 style={styles.updatesTitle}>Updates</h4>
                  {incident.updates.map((update) => (
                    <div key={update.id} style={styles.update}>
                      <div style={styles.updateTime}>
                        {new Date(update.createdAt).toLocaleString()}
                      </div>
                      <div style={styles.updateMessage}>{update.message}</div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
};

const getSystemStatusStyle = (status) => {
  const baseStyle = {
    padding: '0.5rem 1rem',
    borderRadius: '4px',
    fontSize: '1.125rem',
    fontWeight: 'bold',
  };

  switch (status) {
    case 'All Systems Operational':
      return { ...baseStyle, backgroundColor: '#4caf50', color: 'white' };
    case 'Partial Outage':
      return { ...baseStyle, backgroundColor: '#ff9800', color: 'white' };
    case 'Major Outage':
      return { ...baseStyle, backgroundColor: '#f44336', color: 'white' };
    default:
      return { ...baseStyle, backgroundColor: '#9e9e9e', color: 'white' };
  }
};

const getStatusStyle = (status) => {
  const baseStyle = {
    padding: '0.25rem 0.5rem',
    borderRadius: '4px',
    display: 'inline-block',
    fontSize: '0.875rem',
    fontWeight: '500',
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
    fontSize: '0.875rem',
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
    maxWidth: '900px',
    margin: '0 auto',
    padding: '2rem',
    fontFamily: 'sans-serif',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '2rem',
  },
  title: {
    margin: 0,
    fontSize: '2rem',
  },
  section: {
    marginBottom: '3rem',
  },
  sectionTitle: {
    borderBottom: '1px solid #eee',
    paddingBottom: '0.5rem',
    marginBottom: '1.5rem',
  },
  servicesList: {
    display: 'flex',
    flexDirection: 'column',
    gap: '0.75rem',
  },
  serviceItem: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '1rem',
    backgroundColor: '#f9f9f9',
    borderRadius: '4px',
  },
  serviceName: {
    fontWeight: '500',
  },
  incidentsList: {
    display: 'flex',
    flexDirection: 'column',
    gap: '1.5rem',
  },
  incidentCard: {
    padding: '1.5rem',
    backgroundColor: '#f9f9f9',
    borderRadius: '8px',
    boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)',
  },
  incidentHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '0.75rem',
  },
  incidentTitle: {
    margin: 0,
    fontSize: '1.25rem',
  },
  incidentDescription: {
    margin: '0 0 1rem',
    color: '#555',
  },
  affectedServices: {
    marginBottom: '1rem',
    fontSize: '0.875rem',
    color: '#666',
  },
  updates: {
    borderTop: '1px solid #ddd',
    paddingTop: '1rem',
  },
  updatesTitle: {
    fontSize: '1rem',
    marginTop: 0,
    marginBottom: '0.75rem',
  },
  update: {
    marginBottom: '1rem',
    paddingBottom: '1rem',
    borderBottom: '1px solid #eee',
  },
  updateTime: {
    fontSize: '0.75rem',
    color: '#666',
    marginBottom: '0.25rem',
  },
  updateMessage: {
    color: '#333',
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
};

export default PublicStatusPage; 
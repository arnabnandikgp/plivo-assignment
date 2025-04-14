import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

// Create axios instance
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Auth
export const login = (email, password) => {
  return api.post('/auth/login', { email, password });
};

export const signup = (email, password, orgName, orgId) => {
  return api.post('/auth/signup', { email, password, org_name: orgName, org_id: orgId });
};

// Services
export const getServices = () => {
  return api.get('/services');
};

export const getService = (id) => {
  return api.get(`/services/${id}`);
};

export const createService = (name, status) => {
  return api.post('/services', { name, status });
};

export const updateService = (id, name, status) => {
  return api.put(`/services/${id}`, { name, status });
};

export const deleteService = (id) => {
  return api.delete(`/services/${id}`);
};

// Incidents
export const getIncidents = () => {
  return api.get('/incidents');
};

export const getIncident = (id) => {
  return api.get(`/incidents/${id}`);
};

export const createIncident = (title, description, status, serviceIds) => {
  return api.post('/incidents', { title, description, status, serviceIds });
};

export const updateIncident = (id, title, description, status, serviceIds) => {
  return api.put(`/incidents/${id}`, { title, description, status, serviceIds });
};

export const deleteIncident = (id) => {
  return api.delete(`/incidents/${id}`);
};

export const addIncidentUpdate = (id, message) => {
  return api.post(`/incidents/${id}/updates`, { message });
};

// Public status page
export const getPublicServices = (orgId) => {
  return api.get(`/public/${orgId}/services`);
};

export const getPublicIncidents = (orgId) => {
  return api.get(`/public/${orgId}/incidents`);
};

// Websocket setup (to be used in context)
export const getWebSocketUrl = (orgId) => {
  return `ws://localhost:8080/api/ws/${orgId}`;
}; 
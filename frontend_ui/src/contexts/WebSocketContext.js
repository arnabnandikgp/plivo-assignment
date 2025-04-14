import React, { createContext, useContext, useEffect, useState } from 'react';
import { getWebSocketUrl } from '../services/api';
import { useAuth } from './AuthContext';

const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }) => {
  const [socket, setSocket] = useState(null);
  const [connected, setConnected] = useState(false);
  const [events, setEvents] = useState([]);
  const { user, isAuthenticated } = useAuth();

  useEffect(() => {
    let ws = null;

    if (isAuthenticated && user?.orgId) {
      const url = getWebSocketUrl(user.orgId);
      ws = new WebSocket(url);

      ws.onopen = () => {
        console.log('WebSocket connected');
        setConnected(true);
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          setEvents((prevEvents) => [...prevEvents, data]);
        } catch (error) {
          console.error('Error parsing WebSocket message', error);
        }
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
        setConnected(false);
      };

      ws.onerror = (error) => {
        console.error('WebSocket error', error);
        setConnected(false);
      };

      setSocket(ws);
    }

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [isAuthenticated, user]);

  const value = {
    socket,
    connected,
    events,
  };

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
};

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
}; 
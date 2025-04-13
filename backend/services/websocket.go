package services

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketService manages WebSocket connections
type WebSocketService struct {
	// Maps organization ID to a list of connections
	clients map[string][]*websocket.Conn
	mutex   sync.RWMutex
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[string][]*websocket.Conn),
	}
}

// RegisterClient adds a WebSocket connection for an organization
func (s *WebSocketService) RegisterClient(orgID string, conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.clients[orgID]; !ok {
		s.clients[orgID] = make([]*websocket.Conn, 0)
	}
	s.clients[orgID] = append(s.clients[orgID], conn)
	log.Printf("Client registered for organization %s. Total clients: %d", orgID, len(s.clients[orgID]))
}

// UnregisterClient removes a WebSocket connection
func (s *WebSocketService) UnregisterClient(orgID string, conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if connections, ok := s.clients[orgID]; ok {
		for i, c := range connections {
			if c == conn {
				// Remove connection from list
				s.clients[orgID] = append(connections[:i], connections[i+1:]...)
				log.Printf("Client unregistered for organization %s. Remaining clients: %d", orgID, len(s.clients[orgID]))
				break
			}
		}
	}
}

// BroadcastToOrganization sends a message to all clients for an organization
func (s *WebSocketService) BroadcastToOrganization(orgID string, event string, data interface{}) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if connections, ok := s.clients[orgID]; ok {
		message := map[string]interface{}{
			"event": event,
			"data":  data,
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal WebSocket message: %v", err)
			return
		}

		for _, conn := range connections {
			err := conn.WriteMessage(websocket.TextMessage, messageBytes)
			if err != nil {
				log.Printf("Failed to send WebSocket message: %v", err)
				// We don't remove the connection here, as it might be a temporary failure
				// Connections will be cleaned up when they're properly closed
			}
		}
	}
}

// WebSocketEvent represents different types of events
const (
	ServiceUpdated  = "SERVICE_UPDATED"
	IncidentCreated = "INCIDENT_CREATED"
	IncidentUpdated = "INCIDENT_UPDATED"
	UpdateAdded     = "UPDATE_ADDED"
)

// --->>here<<--- WebSocket service for real-time updates

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/status_page/backend/services"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for testing, restrict in production
		},
	}
	// WebsocketService is a global instance of the WebSocket service
	WebsocketService = services.NewWebSocketService()
)

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *gin.Context) {
	orgID := c.Param("orgId")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket connection"})
		return
	}

	// Register client
	WebsocketService.RegisterClient(orgID, conn)

	// Handle disconnection
	go func() {
		defer func() {
			WebsocketService.UnregisterClient(orgID, conn)
			conn.Close()
		}()

		// Keep reading messages to detect disconnection
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// BroadcastServiceUpdate broadcasts a service update to all clients
func BroadcastServiceUpdate(orgID string, service interface{}) {
	WebsocketService.BroadcastToOrganization(orgID, services.ServiceUpdated, service)
}

// BroadcastIncidentCreated broadcasts incident creation to all clients
func BroadcastIncidentCreated(orgID string, incident interface{}) {
	WebsocketService.BroadcastToOrganization(orgID, services.IncidentCreated, incident)
}

// BroadcastIncidentUpdated broadcasts incident update to all clients
func BroadcastIncidentUpdated(orgID string, incident interface{}) {
	WebsocketService.BroadcastToOrganization(orgID, services.IncidentUpdated, incident)
}

// BroadcastUpdateAdded broadcasts an incident update to all clients
func BroadcastUpdateAdded(orgID string, update interface{}) {
	WebsocketService.BroadcastToOrganization(orgID, services.UpdateAdded, update)
}

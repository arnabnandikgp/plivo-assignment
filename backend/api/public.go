package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/status_page/backend/db"
	"github.com/status_page/backend/models"
)

// GetPublicServices returns services for a public status page
func GetPublicServices(c *gin.Context) {
	orgID := c.Param("orgId")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	// --->>here<<--- Database query to get all services for an organization for the public page
	var services []models.Service
	if err := db.DB.Where("org_id = ?", orgID).Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}

// PublicIncidentResponse represents an incident with its services and updates for the public API
type PublicIncidentResponse struct {
	ID          string                  `json:"id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Status      string                  `json:"status"`
	CreatedAt   string                  `json:"createdAt"`
	UpdatedAt   string                  `json:"updatedAt"`
	Services    []models.Service        `json:"services"`
	Updates     []models.IncidentUpdate `json:"updates"`
}

// GetPublicIncidents returns active incidents for a public status page
func GetPublicIncidents(c *gin.Context) {
	orgID := c.Param("orgId")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	// --->>here<<--- Database query to get active incidents (non-resolved) for the public page
	var incidents []models.Incident
	if err := db.DB.Where("org_id = ? AND status != ?", orgID, "Resolved").
		Order("created_at DESC").
		Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incidents"})
		return
	}

	// For each incident, get the associated services and updates
	var responses []PublicIncidentResponse
	for _, incident := range incidents {
		var incidentServices []models.Service
		var incidentServiceIDs []string

		// Get service IDs for this incident
		if err := db.DB.Model(&models.IncidentService{}).
			Where("incident_id = ?", incident.ID).
			Pluck("service_id", &incidentServiceIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident services"})
			return
		}

		// Get the actual services
		if len(incidentServiceIDs) > 0 {
			if err := db.DB.Where("id IN ?", incidentServiceIDs).Find(&incidentServices).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
				return
			}
		}

		// Get updates for this incident
		var updates []models.IncidentUpdate
		if err := db.DB.Where("incident_id = ?", incident.ID).
			Order("created_at DESC").
			Find(&updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident updates"})
			return
		}

		responses = append(responses, PublicIncidentResponse{
			ID:          incident.ID,
			Title:       incident.Title,
			Description: incident.Description,
			Status:      incident.Status,
			CreatedAt:   incident.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   incident.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Services:    incidentServices,
			Updates:     updates,
		})
	}

	c.JSON(http.StatusOK, gin.H{"incidents": responses})
}

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/status_page/backend/db"
	"github.com/status_page/backend/models"
	"github.com/status_page/backend/utils"
	"gorm.io/gorm"
)

// IncidentRequest represents the request for creating/updating an incident
type IncidentRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Status      string   `json:"status" binding:"required"`
	ServiceIDs  []string `json:"serviceIds" binding:"required"`
}

// IncidentUpdateRequest represents the request for adding an update to an incident
type IncidentUpdateRequest struct {
	Message string `json:"message" binding:"required"`
}

// IncidentResponse represents an incident with its services and updates
type IncidentResponse struct {
	Incident models.Incident         `json:"incident"`
	Services []models.Service        `json:"services"`
	Updates  []models.IncidentUpdate `json:"updates"`
}

// GetIncidents returns all incidents for the user's organization
func GetIncidents(c *gin.Context) {
	orgID, _ := c.Get("org_id")

	// --->>here<<--- Database query to get all incidents for an organization
	var incidents []models.Incident
	if err := db.DB.Where("org_id = ?", orgID).Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incidents"})
		return
	}

	// For each incident, get the associated services and updates
	var responses []IncidentResponse
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
		if err := db.DB.Where("incident_id = ?", incident.ID).Find(&updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident updates"})
			return
		}

		responses = append(responses, IncidentResponse{
			Incident: incident,
			Services: incidentServices,
			Updates:  updates,
		})
	}

	c.JSON(http.StatusOK, gin.H{"incidents": responses})
}

// GetIncident returns a specific incident
func GetIncident(c *gin.Context) {
	incidentID := c.Param("id")
	orgID, _ := c.Get("org_id")

	// --->>here<<--- Database query to get a specific incident
	var incident models.Incident
	if err := db.DB.Where("id = ? AND org_id = ?", incidentID, orgID).First(&incident).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident"})
		}
		return
	}

	// Get service IDs for this incident
	var incidentServiceIDs []string
	if err := db.DB.Model(&models.IncidentService{}).
		Where("incident_id = ?", incident.ID).
		Pluck("service_id", &incidentServiceIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident services"})
		return
	}

	// Get the actual services
	var incidentServices []models.Service
	if len(incidentServiceIDs) > 0 {
		if err := db.DB.Where("id IN ?", incidentServiceIDs).Find(&incidentServices).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
			return
		}
	}

	// Get updates for this incident
	var updates []models.IncidentUpdate
	if err := db.DB.Where("incident_id = ?", incident.ID).Find(&updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident updates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"incident": IncidentResponse{
			Incident: incident,
			Services: incidentServices,
			Updates:  updates,
		},
	})
}

// CreateIncident creates a new incident
func CreateIncident(c *gin.Context) {
	var req IncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, _ := c.Get("org_id")

	// Validate status
	validStatuses := map[string]bool{
		"Investigating": true,
		"Identified":    true,
		"Monitoring":    true,
		"Resolved":      true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// Validate that all services exist and belong to the organization
	var count int64
	if err := db.DB.Model(&models.Service{}).
		Where("id IN ? AND org_id = ?", req.ServiceIDs, orgID).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate services"})
		return
	}

	if int(count) != len(req.ServiceIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more services not found"})
		return
	}

	// Create incident
	// --->>here<<--- Database transaction to create an incident and associate services
	tx := db.DB.Begin()

	incidentID := utils.GenerateUUID()
	incident := models.Incident{
		ID:          incidentID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		OrgID:       orgID.(string),
	}

	if err := tx.Create(&incident).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create incident"})
		return
	}

	// Associate services with the incident
	for _, serviceID := range req.ServiceIDs {
		incidentService := models.IncidentService{
			IncidentID: incident.ID,
			ServiceID:  serviceID,
		}

		if err := tx.Create(&incidentService).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate service with incident"})
			return
		}
	}

	// Create first update for the incident
	firstUpdate := models.IncidentUpdate{
		ID:         utils.GenerateUUID(),
		Message:    "Incident reported",
		IncidentID: incident.ID,
	}

	if err := tx.Create(&firstUpdate).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create incident update"})
		return
	}

	tx.Commit()

	// Return the created incident with services and updates
	var incidentServices []models.Service
	if err := db.DB.Where("id IN ?", req.ServiceIDs).Find(&incidentServices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	var updates []models.IncidentUpdate
	if err := db.DB.Where("incident_id = ?", incident.ID).Find(&updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident updates"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"incident": IncidentResponse{
			Incident: incident,
			Services: incidentServices,
			Updates:  updates,
		},
	})
}

// UpdateIncident updates an existing incident
func UpdateIncident(c *gin.Context) {
	var req IncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incidentID := c.Param("id")
	orgID, _ := c.Get("org_id")

	// Validate status
	validStatuses := map[string]bool{
		"Investigating": true,
		"Identified":    true,
		"Monitoring":    true,
		"Resolved":      true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// Validate that all services exist and belong to the organization
	var count int64
	if err := db.DB.Model(&models.Service{}).
		Where("id IN ? AND org_id = ?", req.ServiceIDs, orgID).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate services"})
		return
	}

	if int(count) != len(req.ServiceIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more services not found"})
		return
	}

	// Update incident
	// --->>here<<--- Database transaction to update an incident and its associated services
	tx := db.DB.Begin()

	var incident models.Incident
	if err := tx.Where("id = ? AND org_id = ?", incidentID, orgID).First(&incident).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident"})
		}
		return
	}

	// Update incident fields
	incident.Title = req.Title
	incident.Description = req.Description
	incident.Status = req.Status

	if err := tx.Save(&incident).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update incident"})
		return
	}

	// Remove existing service associations
	if err := tx.Where("incident_id = ?", incidentID).Delete(&models.IncidentService{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update incident services"})
		return
	}

	// Create new service associations
	for _, serviceID := range req.ServiceIDs {
		incidentService := models.IncidentService{
			IncidentID: incident.ID,
			ServiceID:  serviceID,
		}

		if err := tx.Create(&incidentService).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate service with incident"})
			return
		}
	}

	tx.Commit()

	// Get updated services and updates
	var incidentServices []models.Service
	if err := db.DB.Where("id IN ?", req.ServiceIDs).Find(&incidentServices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	var updates []models.IncidentUpdate
	if err := db.DB.Where("incident_id = ?", incident.ID).Find(&updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident updates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"incident": IncidentResponse{
			Incident: incident,
			Services: incidentServices,
			Updates:  updates,
		},
	})
}

// AddIncidentUpdate adds an update to an incident
func AddIncidentUpdate(c *gin.Context) {
	var req IncidentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incidentID := c.Param("id")
	orgID, _ := c.Get("org_id")

	// Check if incident exists and belongs to organization
	var incident models.Incident
	if err := db.DB.Where("id = ? AND org_id = ?", incidentID, orgID).First(&incident).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident"})
		}
		return
	}

	// Create update
	// --->>here<<--- Database operation to add an update to an incident
	update := models.IncidentUpdate{
		ID:         utils.GenerateUUID(),
		Message:    req.Message,
		IncidentID: incidentID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.DB.Create(&update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create incident update"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"update": update})
}

// DeleteIncident deletes an incident
func DeleteIncident(c *gin.Context) {
	incidentID := c.Param("id")
	orgID, _ := c.Get("org_id")

	// --->>here<<--- Database transaction to delete an incident and related records
	tx := db.DB.Begin()

	var incident models.Incident
	if err := tx.Where("id = ? AND org_id = ?", incidentID, orgID).First(&incident).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incident"})
		}
		return
	}

	// Delete incident updates
	if err := tx.Where("incident_id = ?", incidentID).Delete(&models.IncidentUpdate{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete incident updates"})
		return
	}

	// Delete incident-service associations
	if err := tx.Where("incident_id = ?", incidentID).Delete(&models.IncidentService{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete incident services"})
		return
	}

	// Delete the incident
	if err := tx.Delete(&incident).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete incident"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Incident deleted successfully"})
}

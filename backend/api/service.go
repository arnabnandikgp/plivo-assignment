package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/status_page/backend/db"
	"github.com/status_page/backend/models"
	"github.com/status_page/backend/utils"
	"gorm.io/gorm"
)

// ServiceRequest represents the request for creating/updating a service
type ServiceRequest struct {
	Name   string `json:"name" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// GetServices returns all services for the user's organization
func GetServices(c *gin.Context) {
	orgID, _ := c.Get("org_id")

	// --->>here<<--- Database query to get all services for an organization
	var services []models.Service
	if err := db.DB.Where("org_id = ?", orgID).Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}

// GetService returns a specific service
func GetService(c *gin.Context) {
	serviceID := c.Param("id")
	orgID, _ := c.Get("org_id")

	// --->>here<<--- Database query to get a specific service
	var service models.Service
	if err := db.DB.Where("id = ? AND org_id = ?", serviceID, orgID).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service})
}

// CreateService creates a new service
func CreateService(c *gin.Context) {
	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, _ := c.Get("org_id")

	// Validate status
	validStatuses := map[string]bool{
		"Operational": true,
		"Degraded":    true,
		"Outage":      true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// --->>here<<--- Database operation to create a new service
	service := models.Service{
		ID:     utils.GenerateUUID(),
		Name:   req.Name,
		Status: req.Status,
		OrgID:  orgID.(string),
	}

	if err := db.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"service": service})
}

// UpdateService updates an existing service
func UpdateService(c *gin.Context) {
	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceID := c.Param("id")
	orgID, _ := c.Get("org_id")

	// Validate status
	validStatuses := map[string]bool{
		"Operational": true,
		"Degraded":    true,
		"Outage":      true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// --->>here<<--- Database operation to update a service
	var service models.Service
	if err := db.DB.Where("id = ? AND org_id = ?", serviceID, orgID).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service"})
		}
		return
	}

	// Update service
	service.Name = req.Name
	service.Status = req.Status

	if err := db.DB.Save(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service})
}

// DeleteService deletes a service
func DeleteService(c *gin.Context) {
	serviceID := c.Param("id")
	orgID, _ := c.Get("org_id")

	// --->>here<<--- Database operation to delete a service
	var service models.Service
	if err := db.DB.Where("id = ? AND org_id = ?", serviceID, orgID).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service"})
		}
		return
	}

	if err := db.DB.Delete(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

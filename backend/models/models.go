package models

import (
	"time"

	"gorm.io/gorm"
)

// Organization represents a tenant in the system
type Organization struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Users     []User         `gorm:"foreignKey:OrgID"`
	Services  []Service      `gorm:"foreignKey:OrgID"`
}

// User represents a user in the system
type User struct {
	ID        string `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"` // Stored as hashed
	Role      string `gorm:"not null"` // admin, member
	OrgID     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Service represents a service that is monitored
type Service struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Status    string `gorm:"not null"` // Operational, Degraded, Outage
	OrgID     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Incident represents an incident affecting one or more services
type Incident struct {
	ID          string `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Status      string `gorm:"not null"` // Investigating, Identified, Monitoring, Resolved
	OrgID       string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt   `gorm:"index"`
	Updates     []IncidentUpdate `gorm:"foreignKey:IncidentID"`
}

// IncidentUpdate represents an update to an incident
type IncidentUpdate struct {
	ID         string `gorm:"primaryKey"`
	Message    string `gorm:"not null"`
	IncidentID string `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// IncidentService represents the many-to-many relationship between incidents and services
type IncidentService struct {
	IncidentID string `gorm:"primaryKey"`
	ServiceID  string `gorm:"primaryKey"`
}

// --->>here<<--- This is where the database models are defined for PostgreSQL integration with GORM

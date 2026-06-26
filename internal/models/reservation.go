package models

import "time"

// Reservation represents the reservations table in the database.
type Reservation struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserID       uint        `gorm:"not null" json:"user_id"`
	ZoneID       uint        `gorm:"not null" json:"zone_id"`
	LicensePlate string      `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       ReservationStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, completed, cancelled
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	User         User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone         ParkingZone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
}

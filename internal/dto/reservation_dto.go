package dto

import "time"

// CreateReservationRequest is the payload for POST /api/v1/reservations.
type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

// ReservationResponse is returned when a reservation is created.
type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ReservationZoneInfo is a trimmed zone object nested in reservation responses.
type ReservationZoneInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ReservationUserInfo is a trimmed user object nested in admin reservation responses.
type ReservationUserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// MyReservationResponse is used for GET /reservations/my-reservations.
type MyReservationResponse struct {
	ID           uint                `json:"id"`
	LicensePlate string              `json:"license_plate"`
	Status       string              `json:"status"`
	Zone         ReservationZoneInfo `json:"zone"`
	CreatedAt    time.Time           `json:"created_at"`
}

// AdminReservationResponse is used for GET /reservations (admin).
type AdminReservationResponse struct {
	ID           uint                `json:"id"`
	LicensePlate string              `json:"license_plate"`
	Status       string              `json:"status"`
	User         ReservationUserInfo `json:"user"`
	Zone         ReservationZoneInfo `json:"zone"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

package dto

import "time"

// CreateZoneRequest is the payload for POST /api/v1/zones.
type CreateZoneRequest struct {
	Name          string  `json:"name" validate:"required"`
	Type          string  `json:"type" validate:"required,oneof=general ev_charging covered"`
	TotalCapacity int     `json:"total_capacity" validate:"required,gt=0"`
	PricePerHour  float64 `json:"price_per_hour" validate:"required,gt=0"`
}

// UpdateZoneRequest is the payload for PUT /api/v1/zones/:id.
type UpdateZoneRequest struct {
	Name          *string  `json:"name" validate:"omitempty,min=1"`
	Type          *string  `json:"type" validate:"omitempty,oneof=general ev_charging covered"`
	TotalCapacity *int     `json:"total_capacity" validate:"omitempty,gt=0"`
	PricePerHour  *float64 `json:"price_per_hour" validate:"omitempty,gt=0"`
}

// ZoneResponse is used for GET /zones and GET /zones/:id — includes available_spots.
type ZoneResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	TotalCapacity  int       `json:"total_capacity"`
	AvailableSpots int       `json:"available_spots"`
	PricePerHour   float64   `json:"price_per_hour"`
	CreatedAt      time.Time `json:"created_at"`
}

// ZoneDetailResponse is used for POST and PUT responses — includes updated_at.
type ZoneDetailResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	TotalCapacity int       `json:"total_capacity"`
	PricePerHour  float64   `json:"price_per_hour"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

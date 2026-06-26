package models

// Role represents the user role enum type.
type Role string

const (
	RoleDriver Role = "driver"
	RoleAdmin  Role = "admin"
)

// String returns the string representation of the role.
func (r Role) String() string {
	return string(r)
}

// IsValid checks if the role value is one of the allowed enum values.
func (r Role) IsValid() bool {
	switch r {
	case RoleDriver, RoleAdmin:
		return true
	}
	return false
}

// ZoneType represents the parking zone type enum.
type ZoneType string

const (
	ZoneTypeGeneral    ZoneType = "general"
	ZoneTypeEVCharging ZoneType = "ev_charging"
	ZoneTypeCovered    ZoneType = "covered"
)

// String returns the string representation of the zone type.
func (zt ZoneType) String() string {
	return string(zt)
}

// IsValid checks if the zone type is valid.
func (zt ZoneType) IsValid() bool {
	switch zt {
	case ZoneTypeGeneral, ZoneTypeEVCharging, ZoneTypeCovered:
		return true
	}
	return false
}

// ReservationStatus represents the reservation status enum.
type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusCompleted ReservationStatus = "completed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

// String returns the string representation of the reservation status.
func (rs ReservationStatus) String() string {
	return string(rs)
}

// IsValid checks if the reservation status is valid.
func (rs ReservationStatus) IsValid() bool {
	switch rs {
	case ReservationStatusActive, ReservationStatusCompleted, ReservationStatusCancelled:
		return true
	}
	return false
}

package repository

import (
	"errors"
	"spot-sync/internal/models"
	"spot-sync/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrZoneFull is returned when a parking zone has no available spots.
var ErrZoneFull = errors.New("parking zone is full")

// ReservationRepository defines all database operations for reservations.
type ReservationRepository interface {
	// CreateWithLock atomically checks capacity and creates a reservation using FOR UPDATE.
	CreateWithLock(reservation *models.Reservation) error
	FindByID(id uint) (*models.Reservation, error)
	FindByUserID(userID uint, qb *utils.QueryBuilder) ([]models.Reservation, int64, error)
	FindAll(qb *utils.QueryBuilder) ([]models.Reservation, int64, error)
	CancelReservation(id uint) error
}

type reservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a new reservation repository instance.
func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// CreateWithLock uses a GORM transaction with row-level locking (FOR UPDATE)
// to prevent race conditions when reserving parking spots.
// This solves the "EV Spot Bottleneck" concurrency problem.
func (r *reservationRepository) CreateWithLock(reservation *models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the parking zone row to prevent concurrent reads
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, reservation.ZoneID).Error; err != nil {
			return err
		}

		// 2. Count current active reservations for this zone
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", reservation.ZoneID, models.ReservationStatusActive).
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Check if there is available capacity
		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Create the reservation atomically
		return tx.Create(reservation).Error
	})
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Preload("User").Preload("Zone").First(&reservation, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &reservation, err
}

// FindByUserID returns all reservations for a specific user with zone details, paginated.
func (r *reservationRepository) FindByUserID(userID uint, qb *utils.QueryBuilder) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	// Base query with user filter
	baseQuery := r.db.Model(&models.Reservation{}).Where("user_id = ?", userID)

	// Apply search
	query := qb.ApplySearch(baseQuery, []string{"license_plate", "status"})

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination, sorting, preloads, then fetch
	err := qb.ApplyPaginationAndSort(query).
		Preload("Zone").
		Find(&reservations).Error
	return reservations, total, err
}

// FindAll returns all reservations with user and zone details, paginated (admin use).
func (r *reservationRepository) FindAll(qb *utils.QueryBuilder) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	// Apply search
	query := qb.ApplySearch(r.db.Model(&models.Reservation{}), []string{"license_plate", "status"})

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination, sorting, preloads, then fetch
	err := qb.ApplyPaginationAndSort(query).
		Preload("User").Preload("Zone").
		Find(&reservations).Error
	return reservations, total, err
}

// CancelReservation sets the reservation status to "cancelled".
func (r *reservationRepository) CancelReservation(id uint) error {
	return r.db.Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", models.ReservationStatusCancelled).Error
}

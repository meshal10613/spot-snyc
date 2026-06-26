package repository

import (
	"errors"
	"spot-sync/internal/models"
	"spot-sync/pkg/utils"

	"gorm.io/gorm"
)

// ZoneRepository defines all database operations for parking zones.
type ZoneRepository interface {
	Create(zone *models.ParkingZone) error
	FindAll(qb *utils.QueryBuilder) ([]models.ParkingZone, int64, error)
	FindByID(id uint) (*models.ParkingZone, error)
	Update(zone *models.ParkingZone) error
	Delete(id uint) error
	CountActiveReservations(zoneID uint) (int64, error)
}

type zoneRepository struct {
	db *gorm.DB
}

// NewZoneRepository creates a new zone repository instance.
func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

// FindAll retrieves zones with pagination, sorting, and search.
// Returns the list of zones, total count (before pagination), and any error.
func (r *zoneRepository) FindAll(qb *utils.QueryBuilder) ([]models.ParkingZone, int64, error) {
	var zones []models.ParkingZone
	var total int64

	// Apply search filters to base query
	query := qb.ApplySearch(r.db.Model(&models.ParkingZone{}), []string{"name", "type"})

	// Count total matching records (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting, then fetch
	err := qb.ApplyPaginationAndSort(query).Find(&zones).Error
	return zones, total, err
}

func (r *zoneRepository) FindByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	err := r.db.First(&zone, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &zone, err
}

func (r *zoneRepository) Update(zone *models.ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *zoneRepository) Delete(id uint) error {
	return r.db.Delete(&models.ParkingZone{}, id).Error
}

// CountActiveReservations returns the number of active reservations for a zone.
func (r *zoneRepository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, models.ReservationStatusActive).
		Count(&count).Error
	return count, err
}

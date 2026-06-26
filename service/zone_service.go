package service

import (
	"errors"
	"fmt"
	"spot-sync/dto"
	"spot-sync/models"
	"spot-sync/repository"
	"spot-sync/utils"
)

// ZoneService defines the parking zone business logic contract.
type ZoneService interface {
	Create(req *dto.CreateZoneRequest) (*dto.ZoneDetailResponse, error)
	GetAll(qb *utils.QueryBuilder) ([]dto.ZoneResponse, int64, error)
	GetByID(id uint) (*dto.ZoneResponse, error)
	Update(id uint, req *dto.UpdateZoneRequest) (*dto.ZoneDetailResponse, error)
	Delete(id uint) error
}

type zoneService struct {
	repo repository.ZoneRepository
}

// NewZoneService creates a new zone service with injected dependencies.
func NewZoneService(repo repository.ZoneRepository) ZoneService {
	return &zoneService{repo: repo}
}

func (s *zoneService) Create(req *dto.CreateZoneRequest) (*dto.ZoneDetailResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          models.ZoneType(req.Type),
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, fmt.Errorf("failed to create parking zone: %w", err)
	}

	return &dto.ZoneDetailResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          string(zone.Type),
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}, nil
}

// GetAll retrieves all zones with dynamically calculated available_spots, paginated.
func (s *zoneService) GetAll(qb *utils.QueryBuilder) ([]dto.ZoneResponse, int64, error) {
	zones, total, err := s.repo.FindAll(qb)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch parking zones: %w", err)
	}

	responses := make([]dto.ZoneResponse, 0, len(zones))
	for _, z := range zones {
		// Dynamically calculate available spots
		activeCount, err := s.repo.CountActiveReservations(z.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count reservations for zone %d: %w", z.ID, err)
		}

		responses = append(responses, dto.ZoneResponse{
			ID:             z.ID,
			Name:           z.Name,
			Type:           string(z.Type),
			TotalCapacity:  z.TotalCapacity,
			AvailableSpots: z.TotalCapacity - int(activeCount),
			PricePerHour:   z.PricePerHour,
			CreatedAt:      z.CreatedAt,
		})
	}

	return responses, total, nil
}

// GetByID retrieves a single zone with dynamically calculated available_spots.
func (s *zoneService) GetByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if zone == nil {
		return nil, errors.New("parking zone not found")
	}

	activeCount, err := s.repo.CountActiveReservations(id)
	if err != nil {
		return nil, fmt.Errorf("failed to count reservations: %w", err)
	}

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           string(zone.Type),
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity - int(activeCount),
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}, nil
}

func (s *zoneService) Update(id uint, req *dto.UpdateZoneRequest) (*dto.ZoneDetailResponse, error) {
	zone, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if zone == nil {
		return nil, errors.New("parking zone not found")
	}

	// Apply partial updates
	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Type != nil {
		zone.Type = models.ZoneType(*req.Type)
	}
	if req.TotalCapacity != nil {
		zone.TotalCapacity = *req.TotalCapacity
	}
	if req.PricePerHour != nil {
		zone.PricePerHour = *req.PricePerHour
	}

	if err := s.repo.Update(zone); err != nil {
		return nil, fmt.Errorf("failed to update parking zone: %w", err)
	}

	return &dto.ZoneDetailResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          string(zone.Type),
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}, nil
}

func (s *zoneService) Delete(id uint) error {
	zone, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if zone == nil {
		return errors.New("parking zone not found")
	}

	return s.repo.Delete(id)
}

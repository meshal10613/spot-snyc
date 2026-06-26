package service

import (
	"errors"
	"fmt"
	"spot-sync/dto"
	"spot-sync/models"
	"spot-sync/repository"
	"spot-sync/utils"
)

// ReservationService defines the reservation business logic contract.
type ReservationService interface {
	Create(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint, qb *utils.QueryBuilder) ([]dto.MyReservationResponse, int64, error)
	Cancel(reservationID, userID uint, userRole string) error
	GetAll(qb *utils.QueryBuilder) ([]dto.AdminReservationResponse, int64, error)
}

type reservationService struct {
	repo     repository.ReservationRepository
	zoneRepo repository.ZoneRepository
}

// NewReservationService creates a new reservation service with injected dependencies.
func NewReservationService(repo repository.ReservationRepository, zoneRepo repository.ZoneRepository) ReservationService {
	return &reservationService{
		repo:     repo,
		zoneRepo: zoneRepo,
	}
}

// Create reserves a parking spot using a concurrency-safe transaction with row locks.
func (s *reservationService) Create(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	// Verify that the zone exists
	zone, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if zone == nil {
		return nil, errors.New("parking zone not found")
	}

	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       models.ReservationStatusActive,
	}

	// Atomically check capacity and create reservation (FOR UPDATE lock)
	if err := s.repo.CreateWithLock(reservation); err != nil {
		if errors.Is(err, repository.ErrZoneFull) {
			return nil, repository.ErrZoneFull
		}
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       string(reservation.Status),
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}, nil
}

// GetMyReservations returns all reservations for the authenticated user with zone details, paginated.
func (s *reservationService) GetMyReservations(userID uint, qb *utils.QueryBuilder) ([]dto.MyReservationResponse, int64, error) {
	reservations, total, err := s.repo.FindByUserID(userID, qb)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch reservations: %w", err)
	}

	responses := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       string(r.Status),
			Zone: dto.ReservationZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: string(r.Zone.Type),
			},
			CreatedAt: r.CreatedAt,
		})
	}

	return responses, total, nil
}

// Cancel changes a reservation status to "cancelled".
// Drivers can only cancel their own reservations (403 Forbidden otherwise).
func (s *reservationService) Cancel(reservationID, userID uint, userRole string) error {
	reservation, err := s.repo.FindByID(reservationID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}

	// Ownership check: drivers can only cancel their own reservations
	if userRole != string(models.RoleAdmin) && reservation.UserID != userID {
		return errors.New("forbidden")
	}

	if reservation.Status == models.ReservationStatusCancelled {
		return errors.New("reservation is already cancelled")
	}

	return s.repo.CancelReservation(reservationID)
}

// GetAll returns all reservations with user and zone details, paginated (admin only).
func (s *reservationService) GetAll(qb *utils.QueryBuilder) ([]dto.AdminReservationResponse, int64, error) {
	reservations, total, err := s.repo.FindAll(qb)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch reservations: %w", err)
	}

	responses := make([]dto.AdminReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       string(r.Status),
			User: dto.ReservationUserInfo{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: dto.ReservationZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: string(r.Zone.Type),
			},
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	return responses, total, nil
}

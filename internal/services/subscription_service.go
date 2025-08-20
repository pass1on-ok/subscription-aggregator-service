package services

import (
	"context"
	"errors"
	"time"

	"subscription-service/internal/models"
	"subscription-service/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionService interface {
	Create(ctx context.Context, in CreateDTO) (*models.Subscription, error)
	Get(ctx context.Context, id string) (*models.Subscription, error)
	Update(ctx context.Context, id string, in UpdateDTO) (*models.Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter ListFilter) ([]models.Subscription, error)
	Total(ctx context.Context, agg TotalFilter) (int64, error)
}

type service struct {
	repo   repository.SubscriptionRepository
	logger *zap.SugaredLogger
}

func NewSubscriptionService(repo repository.SubscriptionRepository, logger *zap.SugaredLogger) SubscriptionService {
	return &service{repo: repo, logger: logger}
}

type CreateDTO struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=0"`
	UserID      string  `json:"user_id" binding:"required,uuid"`
	StartMonth  string  `json:"start_date" binding:"required"`
	EndMonth    *string `json:"end_date"`
}

type UpdateDTO struct {
	ServiceName *string `json:"service_name"`
	Price       *int    `json:"price"`
	StartMonth  *string `json:"start_date"`
	EndMonth    *string `json:"end_date"`
}

type ListFilter struct {
	UserID      string
	ServiceName string
	From        *time.Time
	To          *time.Time
	Limit       int
	Offset      int
}

type TotalFilter struct {
	UserID      string
	ServiceName string
	From        time.Time
	To          time.Time
}

const monthLayout = "01-2006"

func parseMonth(m string) (time.Time, error) {
	t, err := time.Parse(monthLayout, m)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func (s *service) Create(ctx context.Context, in CreateDTO) (*models.Subscription, error) {
	start, err := parseMonth(in.StartMonth)
	if err != nil {
		return nil, errors.New("invalid start_date, expected MM-YYYY")
	}
	var endPtr *time.Time
	if in.EndMonth != nil && *in.EndMonth != "" {
		e, err := parseMonth(*in.EndMonth)
		if err != nil {
			return nil, errors.New("invalid end_date, expected MM-YYYY")
		}
		endPtr = &e
	}

	uid, err := uuidFromString(in.UserID)
	if err != nil {
		return nil, errors.New("invalid user_id uuid")
	}

	sub := &models.Subscription{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      uid,
		StartDate:   start,
		EndDate:     endPtr,
	}
	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *service) Get(ctx context.Context, id string) (*models.Subscription, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id string, in UpdateDTO) (*models.Subscription, error) {
	ex, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.ServiceName != nil {
		ex.ServiceName = *in.ServiceName
	}
	if in.Price != nil {
		if *in.Price < 0 {
			return nil, errors.New("price must be >= 0")
		}
		ex.Price = *in.Price
	}
	if in.StartMonth != nil {
		st, err := parseMonth(*in.StartMonth)
		if err != nil {
			return nil, errors.New("invalid start_date, expected MM-YYYY")
		}
		ex.StartDate = st
	}
	if in.EndMonth != nil {
		if *in.EndMonth == "" {
			ex.EndDate = nil
		} else {
			ed, err := parseMonth(*in.EndMonth)
			if err != nil {
				return nil, errors.New("invalid end_date, expected MM-YYYY")
			}
			ex.EndDate = &ed
		}
	}
	if err := s.repo.Update(ctx, ex); err != nil {
		return nil, err
	}
	return ex, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, f ListFilter) ([]models.Subscription, error) {
	return s.repo.List(ctx, f.UserID, f.ServiceName, f.From, f.To, f.Limit, f.Offset)
}

func (s *service) Total(ctx context.Context, agg TotalFilter) (int64, error) {
	agg.From = time.Date(agg.From.Year(), agg.From.Month(), 1, 0, 0, 0, 0, time.UTC)
	agg.To = time.Date(agg.To.Year(), agg.To.Month(), 1, 0, 0, 0, 0, time.UTC)
	return s.repo.TotalForPeriod(ctx, agg.UserID, agg.ServiceName, agg.From, agg.To)
}

func uuidFromString(sid string) (uuid.UUID, error) {
	return uuid.Parse(sid)
}

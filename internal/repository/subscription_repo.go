package repository

import (
	"context"
	"time"

	"subscription-service/internal/models"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *models.Subscription) error
	Get(ctx context.Context, id string) (*models.Subscription, error)
	Update(ctx context.Context, s *models.Subscription) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID, serviceName string, from, to *time.Time, limit, offset int) ([]models.Subscription, error)
	TotalForPeriod(ctx context.Context, userID, serviceName string, from, to time.Time) (int64, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, s *models.Subscription) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *subscriptionRepository) Get(ctx context.Context, id string) (*models.Subscription, error) {
	var sub models.Subscription
	if err := r.db.WithContext(ctx).First(&sub, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, s *models.Subscription) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *subscriptionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Subscription{}, "id = ?", id).Error
}

func (r *subscriptionRepository) List(ctx context.Context, userID, serviceName string, from, to *time.Time, limit, offset int) ([]models.Subscription, error) {
	q := r.db.WithContext(ctx).Model(&models.Subscription{})

	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		q = q.Where("service_name = ?", serviceName)
	}
	if from != nil {
		q = q.Where("start_date >= ?", *from)
	}
	if to != nil {
		q = q.Where("start_date <= ?", *to)
	}

	if limit == 0 {
		limit = 100
	}
	var out []models.Subscription
	if err := q.Order("start_date desc, created_at desc NULLS LAST").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *subscriptionRepository) TotalForPeriod(ctx context.Context, userID, serviceName string, from, to time.Time) (int64, error) {
	from = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC)

	type result struct {
		Total int64
	}

	var whereUser, whereService string
	if userID != "" {
		whereUser = " AND s.user_id = ?"
	}
	if serviceName != "" {
		whereService = " AND s.service_name = ?"
	}

	sql := `
WITH period AS (
  SELECT date_trunc('month', ?::date) AS from_m,
         date_trunc('month', ?::date) AS to_m
),
subs AS (
  SELECT s.*
  FROM subscriptions s, period p
  WHERE date_trunc('month', s.start_date) <= (SELECT to_m FROM period)
    AND date_trunc('month', COALESCE(s.end_date, (SELECT to_m FROM period))) >= (SELECT from_m FROM period)
    ` + whereUser + whereService + `
),
exp AS (
  SELECT s.id,
         s.price,
         COUNT(*) AS months_cnt
  FROM subs s, period p,
       generate_series(
         GREATEST(date_trunc('month', s.start_date), p.from_m),
         LEAST( date_trunc('month', COALESCE(s.end_date, p.to_m)), p.to_m),
         interval '1 month'
       ) m
  GROUP BY s.id, s.price
)
SELECT COALESCE(SUM(price * months_cnt), 0) AS total
FROM exp;
`

	var res result
	args := []any{from, to}
	if userID != "" {
		args = append(args, userID)
	}
	if serviceName != "" {
		args = append(args, serviceName)
	}

	tx := r.db.WithContext(ctx).
		Raw(sql, args...).
		Scan(&res)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return res.Total, nil
}

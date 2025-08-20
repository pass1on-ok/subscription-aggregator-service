package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string     `gorm:"type:text;not null;index:idx_service_user"       json:"service_name" example:"Yandex Plus"`
	Price       int        `gorm:"type:int;not null"                               json:"price" example:"400"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_service_user"       json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `gorm:"type:date;not null;index:idx_start_end"          json:"start_date" example:"2025-07-01"`
	EndDate     *time.Time `gorm:"type:date;index:idx_start_end"                   json:"end_date,omitempty" example:"2025-10-01"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

func (Subscription) TableName() string { return "subscriptions" }

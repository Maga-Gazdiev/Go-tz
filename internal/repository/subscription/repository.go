package subscription

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("subscription not found")

type Subscription struct {
	ID          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type ListFilter struct {
	UserID, ServiceName *string
	Limit, Offset       int
}
type CreateParams struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
}
type UpdateParams struct {
	ServiceName *string
	Price       *int
	UserID      *string
	StartDate   *string
	EndDate     *string
}

type Repository interface {
	Create(context.Context, CreateParams) (*Subscription, error)
	GetByID(context.Context, string) (*Subscription, error)
	Update(context.Context, string, UpdateParams) (*Subscription, error)
	Delete(context.Context, string) error
	List(context.Context, ListFilter) ([]Subscription, error)
}
type repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

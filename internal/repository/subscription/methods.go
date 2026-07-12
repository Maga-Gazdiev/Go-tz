package subscription

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (r *repository) Create(ctx context.Context, params CreateParams) (*Subscription, error) {
	start, err := parseMonthYear(params.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := optionalMonthYear(params.EndDate)
	if err != nil {
		return nil, err
	}

	sub := &Subscription{}
	var startDate time.Time
	var endDate *time.Time
	err = r.db.QueryRow(ctx, createSubscription, params.ServiceName, params.Price, params.UserID, start, end).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&startDate,
		&endDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	sub.StartDate = startDate.Format("01-2006")
	if endDate != nil {
		value := endDate.Format("01-2006")
		sub.EndDate = &value
	}
	return sub, nil
}
func (r *repository) GetByID(ctx context.Context, id string) (*Subscription, error) {
	sub := &Subscription{}
	var startDate time.Time
	var endDate *time.Time
	err := r.db.QueryRow(ctx, getSubscription, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&startDate,
		&endDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	sub.StartDate = startDate.Format("01-2006")
	if endDate != nil {
		value := endDate.Format("01-2006")
		sub.EndDate = &value
	}
	return sub, nil
}
func (r *repository) Update(ctx context.Context, id string, params UpdateParams) (*Subscription, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if params.ServiceName != nil {
		existing.ServiceName = *params.ServiceName
	}
	if params.Price != nil {
		existing.Price = *params.Price
	}
	if params.UserID != nil {
		existing.UserID = *params.UserID
	}
	if params.StartDate != nil {
		existing.StartDate = *params.StartDate
	}
	if params.EndDate != nil {
		existing.EndDate = params.EndDate
	}
	start, err := parseMonthYear(existing.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := optionalMonthYear(existing.EndDate)
	if err != nil {
		return nil, err
	}

	sub := &Subscription{}
	var startDate time.Time
	var endDate *time.Time
	err = r.db.QueryRow(ctx, updateSubscription, id, existing.ServiceName, existing.Price, existing.UserID, start, end).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&startDate,
		&endDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	sub.StartDate = startDate.Format("01-2006")
	if endDate != nil {
		value := endDate.Format("01-2006")
		sub.EndDate = &value
	}
	return sub, nil
}
func (r *repository) Delete(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, deleteSubscription, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *repository) List(ctx context.Context, filter ListFilter) ([]Subscription, error) {
	if filter.Limit < 1 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	var userID, serviceName any
	if filter.UserID != nil {
		userID = *filter.UserID
	}
	if filter.ServiceName != nil && *filter.ServiceName != "" {
		serviceName = *filter.ServiceName
	}
	rows, err := r.db.Query(ctx, listSubscriptions, userID, serviceName, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	subs := make([]Subscription, 0)
	for rows.Next() {
		var sub Subscription
		var startDate time.Time
		var endDate *time.Time
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&startDate,
			&endDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sub.StartDate = startDate.Format("01-2006")
		if endDate != nil {
			value := endDate.Format("01-2006")
			sub.EndDate = &value
		}
		subs = append(subs, sub)
	}
	return subs, rows.Err()
}
func parseMonthYear(value string) (time.Time, error) {
	date, err := time.Parse("01-2006", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный формат даты %q, ожидается MM-YYYY", value)
	}
	return date, nil
}
func optionalMonthYear(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	date, err := parseMonthYear(*value)
	return &date, err
}

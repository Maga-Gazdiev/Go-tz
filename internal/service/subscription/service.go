package subscription

import (
	"context"
	"fmt"
	"time"

	repo "subscriptions/internal/repository/subscription"
)

type TotalParams struct {
	UserID, ServiceName *string
	From, To            string
}

type Service struct{ repository repo.Repository }

func NewService(repository repo.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, p repo.CreateParams) (*repo.Subscription, error) {
	if p.ServiceName == "" {
		return nil, fmt.Errorf("service_name обязателен")
	}
	if p.Price <= 0 {
		return nil, fmt.Errorf("price должен быть больше 0")
	}
	return s.repository.Create(ctx, p)
}

func (s *Service) Get(ctx context.Context, id string) (*repo.Subscription, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, p repo.UpdateParams) (*repo.Subscription, error) {
	if p.Price != nil && *p.Price <= 0 {
		return nil, fmt.Errorf("price должен быть больше 0")
	}
	return s.repository.Update(ctx, id, p)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, f repo.ListFilter) ([]repo.Subscription, error) {
	if f.Limit == 0 {
		f.Limit = 50
	}
	return s.repository.List(ctx, f)
}

func (s *Service) CalcTotal(ctx context.Context, p TotalParams) (int, error) {
	from, err := monthYear(p.From)
	if err != nil {
		return 0, err
	}
	to, err := monthYear(p.To)
	if err != nil {
		return 0, err
	}
	if from.After(to) {
		return 0, fmt.Errorf("from не может быть позже to")
	}
	subs, err := s.repository.List(ctx, repo.ListFilter{UserID: p.UserID, ServiceName: p.ServiceName, Limit: 10000})
	if err != nil {
		return 0, err
	}
	total := 0
	for _, sub := range subs {
		start, _ := monthYear(sub.StartDate)
		var end *time.Time
		if sub.EndDate != nil {
			v, _ := monthYear(*sub.EndDate)
			end = &v
		}
		total += sub.Price * monthsInPeriod(&start, end, from, to)
	}
	return total, nil
}

func monthYear(value string) (time.Time, error) {
	t, err := time.Parse("01-2006", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный формат даты %q, ожидается MM-YYYY", value)
	}
	return t, nil
}

func monthsInPeriod(start, end *time.Time, from, to time.Time) int {
	if start.Before(from) {
		start = &from
	}
	if end != nil && end.Before(to) {
		to = *end
	}
	if start.After(to) {
		return 0
	}
	months := 0
	for date := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC); !date.After(to); date = date.AddDate(0, 1, 0) {
		months++
	}
	return months
}

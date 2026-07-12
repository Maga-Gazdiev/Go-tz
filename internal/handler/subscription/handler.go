package subscription

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	repo "subscriptions/internal/repository/subscription"
	service "subscriptions/internal/service/subscription"
)

type Handler struct {
	service *service.Service
	log     *slog.Logger
}

func NewHandler(service *service.Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /api/v1/subscriptions", h.Create)
	mux.HandleFunc("GET /api/v1/subscriptions", h.List)
	mux.HandleFunc("GET /api/v1/subscriptions/total", h.Total)
	mux.HandleFunc("GET /api/v1/subscriptions/{id}", h.Get)
	mux.HandleFunc("PUT /api/v1/subscriptions/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/subscriptions/{id}", h.Delete)
}

type CreateRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}
type UpdateRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	UserID      *string `json:"user_id,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		h.error(w, http.StatusBadRequest, "некорректный json")
		return
	}
	if req.UserID = normalizeID(req.UserID); req.UserID == "" {
		h.error(w, http.StatusBadRequest, "невалидный user_id")
		return
	}
	sub, err := h.service.Create(r.Context(), repo.CreateParams{ServiceName: req.ServiceName, Price: req.Price, UserID: req.UserID, StartDate: req.StartDate, EndDate: req.EndDate})
	if err != nil {
		h.error(w, http.StatusBadRequest, err.Error())
		return
	}
	h.json(w, http.StatusCreated, subscriptionResponse(sub))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := normalizeID(r.PathValue("id"))
	if id == "" {
		h.error(w, http.StatusBadRequest, "невалидный id")
		return
	}
	sub, err := h.service.Get(r.Context(), id)
	if errors.Is(err, repo.ErrNotFound) {
		h.error(w, http.StatusNotFound, "не найдено")
		return
	}
	if err != nil {
		h.error(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	h.json(w, http.StatusOK, subscriptionResponse(sub))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := normalizeID(r.PathValue("id"))
	if id == "" {
		h.error(w, http.StatusBadRequest, "невалидный id")
		return
	}
	var req UpdateRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		h.error(w, http.StatusBadRequest, "некорректный json")
		return
	}
	if req.UserID != nil {
		normalized := normalizeID(*req.UserID)
		if normalized == "" {
			h.error(w, http.StatusBadRequest, "невалидный user_id")
			return
		}
		req.UserID = &normalized
	}
	sub, err := h.service.Update(r.Context(), id, repo.UpdateParams{ServiceName: req.ServiceName, Price: req.Price, UserID: req.UserID, StartDate: req.StartDate, EndDate: req.EndDate})
	if errors.Is(err, repo.ErrNotFound) {
		h.error(w, http.StatusNotFound, "не найдено")
		return
	}
	if err != nil {
		h.error(w, http.StatusBadRequest, err.Error())
		return
	}
	h.json(w, http.StatusOK, subscriptionResponse(sub))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := normalizeID(r.PathValue("id"))
	if id == "" {
		h.error(w, http.StatusBadRequest, "невалидный id")
		return
	}
	err := h.service.Delete(r.Context(), id)
	if errors.Is(err, repo.ErrNotFound) {
		h.error(w, http.StatusNotFound, "не найдено")
		return
	}
	if err != nil {
		h.error(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := repo.ListFilter{}
	if value := q.Get("user_id"); value != "" {
		value = normalizeID(value)
		if value == "" {
			h.error(w, http.StatusBadRequest, "невалидный user_id")
			return
		}
		filter.UserID = &value
	}
	if value := q.Get("service_name"); value != "" {
		filter.ServiceName = &value
	}
	filter.Limit, _ = strconv.Atoi(q.Get("limit"))
	filter.Offset, _ = strconv.Atoi(q.Get("offset"))
	subs, err := h.service.List(r.Context(), filter)
	if err != nil {
		h.error(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	responses := make([]SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		responses = append(responses, subscriptionResponse(&sub))
	}
	h.json(w, http.StatusOK, responses)
}

func (h *Handler) Total(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if q.Get("from") == "" || q.Get("to") == "" {
		h.error(w, http.StatusBadRequest, "параметры from и to обязательны")
		return
	}
	params := service.TotalParams{From: q.Get("from"), To: q.Get("to")}
	if value := q.Get("user_id"); value != "" {
		value = normalizeID(value)
		if value == "" {
			h.error(w, http.StatusBadRequest, "невалидный user_id")
			return
		}
		params.UserID = &value
	}
	if value := q.Get("service_name"); value != "" {
		params.ServiceName = &value
	}
	total, err := h.service.CalcTotal(r.Context(), params)
	if err != nil {
		h.error(w, http.StatusBadRequest, err.Error())
		return
	}
	h.json(w, http.StatusOK, map[string]int{"total": total})
}

func (h *Handler) json(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (h *Handler) error(w http.ResponseWriter, status int, message string) {
	h.json(w, status, map[string]string{"error": message})
}

func normalizeID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return ""
	}
	for i := range value {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !((value[i] >= '0' && value[i] <= '9') || (value[i] >= 'a' && value[i] <= 'f')) {
			return ""
		}
	}
	return value
}

func subscriptionResponse(sub *repo.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}
}

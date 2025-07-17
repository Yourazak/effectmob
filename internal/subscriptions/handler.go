package subscriptions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createSubscription(w, r)
	case http.MethodGet:
		if r.URL.Query().Get("sum") == "true" {
			h.sumPriceHandler(w, r)
		} else if idStr := r.URL.Query().Get("id"); idStr != "" {
			h.getByIDHandler(w, r)
		} else {
			h.listSubscriptions(w, r)
		}
	case http.MethodPut:
		h.updateSubscription(w, r)
	case http.MethodDelete:
		h.deleteSubscription(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var sub Subscription
	if err := json.Unmarshal(body, &sub); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	err = h.service.Create(r.Context(), sub)
	if err != nil {
		http.Error(w, "Ошибка создания подписки", http.StatusInternalServerError)
		log.Println("Ошибка создания подписки:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, "Ошибка получения подписок", http.StatusInternalServerError)
		log.Println("Ошибка получения подписок:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Подписка не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var sub Subscription
	if err := json.Unmarshal(body, &sub); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if sub.ID == 0 {
		http.Error(w, "ID обязателен для обновления", http.StatusBadRequest)
		return
	}

	err = h.service.Update(r.Context(), sub)
	if err != nil {
		http.Error(w, "Ошибка обновления подписки", http.StatusInternalServerError)
		log.Println("Ошибка обновления подписки:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "Ошибка удаления подписки", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) sumPriceHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")
	start := query.Get("start")
	end := query.Get("end")

	if userID == "" || serviceName == "" || start == "" || end == "" {
		http.Error(w, "Отсутствуют обязательные параметры", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01", start)
	if err != nil {
		http.Error(w, "Некорректный формат даты начала", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01", end)
	if err != nil {
		http.Error(w, "Некорректный формат даты конца", http.StatusBadRequest)
		return
	}
	// Чтобы включить весь последний месяц
	endDate = endDate.AddDate(0, 1, -1)

	total, err := h.service.SumPrice(r.Context(), userID, serviceName, startDate, endDate)
	if err != nil {
		http.Error(w, "Ошибка подсчета стоимости", http.StatusInternalServerError)
		log.Println("Ошибка подсчета стоимости:", err)
		return
	}

	resp := map[string]int{"total_price": total}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	}
}

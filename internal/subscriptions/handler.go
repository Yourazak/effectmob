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

// createSubscription создает новую подписку
// @Summary Создать подписку
// @Description Создает новую запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body subscriptions.Subscription true "Данные подписки"
// @Success 201 "Подписка создана"
// @Failure 400 "Некорректный запрос"
// @Failure 500 "Ошибка сервера"
// @Router /subscriptions [post]
func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Handler: ошибка чтения тела запроса:", err)
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var sub Subscription
	if err := json.Unmarshal(body, &sub); err != nil {
		log.Printf("Handler: некорректный JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	err = h.service.Create(r.Context(), sub)
	if err != nil {
		log.Println("Handler: ошибка создания подписки:", err)
		http.Error(w, "Ошибка создания подписки", http.StatusInternalServerError)
		return
	}

	log.Printf("Handler: подписка создана — %+v", sub)
	w.WriteHeader(http.StatusCreated)
}

// listSubscriptions возвращает список всех подписок
// @Summary Получить список подписок
// @Description Возвращает список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} subscriptions.Subscription
// @Failure 500 "Ошибка сервера"
// @Router /subscriptions [get]
func (h *Handler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.service.List(r.Context())
	if err != nil {
		log.Println("Handler: ошибка получения подписок:", err)
		http.Error(w, "Ошибка получения подписок", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		log.Println("Handler: ошибка кодирования ответа:", err)
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	}
}

// getByIDHandler возвращает подписку по ID
// @Summary Получить подписку по ID
// @Description Возвращает подписку с указанным ID
// @Tags subscriptions
// @Produce json
// @Param id query int true "ID подписки"
// @Success 200 {object} subscriptions.Subscription
// @Failure 400 "Некорректный ID"
// @Failure 404 "Подписка не найдена"
// @Router /subscriptions [get]
func (h *Handler) getByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Handler: некорректный ID: %v", err)
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Handler: подписка с ID %d не найдена: %v", id, err)
		http.Error(w, "Подписка не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		log.Println("Handler: ошибка кодирования ответа:", err)
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	}
}

// updateSubscription обновляет данные подписки
// @Summary Обновить подписку
// @Description Обновляет подписку с указанным ID
// @Tags subscriptions
// @Accept json
// @Param subscription body subscriptions.Subscription true "Данные подписки"
// @Success 200 "Подписка обновлена"
// @Failure 400 "Некорректный запрос"
// @Failure 500 "Ошибка сервера"
// @Router /subscriptions [put]
func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Handler: ошибка чтения тела запроса:", err)
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var sub Subscription
	if err := json.Unmarshal(body, &sub); err != nil {
		log.Printf("Handler: некорректный JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if sub.ID == 0 {
		log.Println("Handler: ID обязателен для обновления, но не передан")
		http.Error(w, "ID обязателен для обновления", http.StatusBadRequest)
		return
	}

	err = h.service.Update(r.Context(), sub)
	if err != nil {
		log.Println("Handler: ошибка обновления подписки:", err)
		http.Error(w, "Ошибка обновления подписки", http.StatusInternalServerError)
		return
	}

	log.Printf("Handler: подписка с ID %d обновлена", sub.ID)
	w.WriteHeader(http.StatusOK)
}

// deleteSubscription удаляет подписку по ID
// @Summary Удалить подписку
// @Description Удаляет подписку с указанным ID
// @Tags subscriptions
// @Param id query int true "ID подписки"
// @Success 200 "Подписка удалена"
// @Failure 400 "Некорректный ID"
// @Failure 500 "Ошибка сервера"
// @Router /subscriptions [delete]
func (h *Handler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	log.Printf("Handler: удаление подписки, получен ID = %s", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Handler: некорректный ID: %s, ошибка: %v", idStr, err)
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		log.Printf("Handler: ошибка при удалении подписки с ID %d: %v", id, err)
		http.Error(w, "Ошибка удаления подписки", http.StatusInternalServerError)
		return
	}

	log.Printf("Handler: подписка с ID %d успешно удалена", id)
	w.WriteHeader(http.StatusOK)
}

// sumPriceHandler считает суммарную стоимость подписок за период
// @Summary Суммировать стоимость подписок
// @Description Считает суммарную стоимость всех подписок пользователя с фильтрацией по имени сервиса и дате
// @Tags subscriptions
// @Produce json
// @Param sum query bool true "Активировать подсчет суммы"
// @Param user_id query string true "UUID пользователя"
// @Param service_name query string true "Название сервиса"
// @Param start query string true "Дата начала периода в формате YYYY-MM"
// @Param end query string true "Дата конца периода в формате YYYY-MM"
// @Success 200 {object} map[string]int "Объект с total_price"
// @Failure 400 "Некорректные параметры"
// @Failure 500 "Ошибка сервера"
// @Router /subscriptions [get]
func (h *Handler) sumPriceHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")
	start := query.Get("start")
	end := query.Get("end")

	if userID == "" || serviceName == "" || start == "" || end == "" {
		log.Println("Handler: отсутствуют обязательные параметры для подсчета")
		http.Error(w, "Отсутствуют обязательные параметры", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01", start)
	if err != nil {
		log.Println("Handler: некорректный формат даты начала:", err)
		http.Error(w, "Некорректный формат даты начала", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01", end)
	if err != nil {
		log.Println("Handler: некорректный формат даты конца:", err)
		http.Error(w, "Некорректный формат даты конца", http.StatusBadRequest)
		return
	}
	endDate = endDate.AddDate(0, 1, -1) // включаем весь месяц

	log.Printf("Handler: подсчет суммы, user_id=%s, service_name=%s, start=%s, end=%s", userID, serviceName, start, end)

	total, err := h.service.SumPrice(r.Context(), userID, serviceName, startDate, endDate)
	if err != nil {
		log.Println("Handler: ошибка подсчета стоимости:", err)
		http.Error(w, "Ошибка подсчета стоимости", http.StatusInternalServerError)
		return
	}

	log.Printf("Handler: итоговая сумма подписок = %d", total)

	resp := map[string]int{"total_price": total}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Handler: ошибка кодирования ответа:", err)
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	}
}

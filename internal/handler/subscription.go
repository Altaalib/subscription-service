package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"subscription-service/internal/model"
	"subscription-service/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// Create godoc
// @Summary Создать подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.CreateSubscriptionInput true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c echo.Context) error {
	var input model.CreateSubscriptionInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	sub, err := h.service.Create(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, sub)
}

// GetByID godoc
// @Summary Получить подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	sub, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, sub)
}

// GetAll godoc
// @Summary Получить все подписки
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c echo.Context) error {
	subs, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, subs)
}

// Update godoc
// @Summary Обновить подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param input body model.UpdateSubscriptionInput true "Данные для обновления"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var input model.UpdateSubscriptionInput
	if bindErr := c.Bind(&input); bindErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": bindErr.Error()})
	}

	sub, err := h.service.Update(c.Request().Context(), id, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, sub)
}

// Delete godoc
// @Summary Удалить подписку
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusNoContent, nil)
}

// GetTotalCost godoc
// @Summary Получить суммарную стоимость подписок
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "ID пользователя"
// @Param service_name query string true "Название сервиса"
// @Success 200 {object} model.TotalCostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) GetTotalCost(c echo.Context) error {
	userIDStr := c.QueryParam("user_id")
	serviceName := c.QueryParam("service_name")

	if userIDStr == "" || serviceName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id and service_name are required"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	result, err := h.service.GetTotalCost(c.Request().Context(), userID, serviceName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

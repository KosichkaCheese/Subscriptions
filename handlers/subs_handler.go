package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"subscriptions/models"
	"subscriptions/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	service services.SubscriptionServiceInterface
}

func NewSubscriptionHandler(service services.SubscriptionServiceInterface) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// @Summary Получить список подписок
// @Schemes
// @Description Возвращает список всех подписок
// @Tags Subscription
// @Accept json
// @Produce json
// @Success 200 {array} models.Subscription
// @Router /subs [get]
func (handler *SubscriptionHandler) GetAll(c *gin.Context) {
	subscriptions, err := handler.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subscriptions)
}

// @Summary Добавить новую подписку
// @Schemes
// @Description Добавляет новую подписку
// @Tags Subscription
// @Accept json
// @Produce json
// @Param subscription body models.CreateSubscription true "Subscription"
// @Success 200 {object} models.Subscription
// @Router /subs [post]
func (handler *SubscriptionHandler) Create(c *gin.Context) {
	var subscription models.CreateSubscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newSubscription, err := handler.service.Create(c.Request.Context(), &subscription)
	if err != nil {
		if err == services.ErrInvalidDate {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newSubscription)
}

// @Summary Получить подписку по ID
// @Schemes
// @Description Возвращает подписку по ID
// @Tags Subscription
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Subscription
// @Router /subs/{id} [get]
func (handler *SubscriptionHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subscription, err := handler.service.GetById(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subscription)
}

// @Summary Обновить подписку
// @Schemes
// @Description Обновляет существующую подписку
// @Tags Subscription
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param subscription body models.UpdateSubscription true "Subscription"
// @Success 200 {object} models.Subscription
// @Router /subs/{id} [put]
func (handler *SubscriptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var upd models.UpdateSubscription
	if err := c.ShouldBindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub, err := handler.service.Update(c.Request.Context(), uint(id), &upd)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		if err == services.ErrInvalidDate {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// @Summary Удалить подписку
// @Schemes
// @Description Удаляет существующую подписку
// @Tags Subscription
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} map[string]interface{}
// @Router /subs/{id} [delete]
func (handler *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = handler.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

// @Summary Получить сумму подписок по фильтрам
// @Schemes
// @Description Возвращает сумму подписок по фильтрам
// @Tags Subscription
// @Accept json
// @Produce json
// @Param filters query models.SumFilter true "Filters"
// @Success 200 {object} map[string]interface{}
// @Router /subs/sum [get]
func (handler *SubscriptionHandler) SumByFilters(c *gin.Context) {
	var filters models.SumFilter

	err := c.ShouldBindQuery(&filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sum, err := handler.service.SumByFilters(c.Request.Context(), &filters)
	if err != nil {
		if err == services.ErrInvalidDate {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sum": sum})
}

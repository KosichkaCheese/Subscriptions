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

type ServiceHandler struct {
	service services.ServiceServiceInterface
}

func NewServiceHandler(service services.ServiceServiceInterface) *ServiceHandler {
	return &ServiceHandler{service: service}
}

// @Summary Получить список сервисов
// @Schemes
// @Description Возвращает список всех сервисов
// @Tags Service
// @Accept json
// @Produce json
// @Success 200 {array} models.Service
// @Router /services [get]
func (handler *ServiceHandler) GetAll(c *gin.Context) {
	services, err := handler.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

// @Summary Создать новый сервис
// @Schemes
// @Description Добавляет новый сервис
// @Tags Service
// @Accept json
// @Produce json
// @Param service body models.CreateService true "Service"
// @Success 201 {object} models.Service
// @Router /services [post]
func (handler *ServiceHandler) Create(c *gin.Context) {
	var service models.CreateService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newService, err := handler.service.Create(c.Request.Context(), &service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newService)
}

// @Summary Удалить сервис
// @Schemes
// @Description Удаляет существующий сервис
// @Tags Service
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} map[string]interface{}
// @Router /services/{id} [delete]
func (handler *ServiceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = handler.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

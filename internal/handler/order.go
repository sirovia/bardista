package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirovia/bardista/internal/repository"
	"github.com/sirovia/bardista/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type orderItemInput struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

type createOrderRequest struct {
	Items []orderItemInput `json:"items" binding:"required,min=1"`
}

type updateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_INPUT", "message": err.Error()},
		})
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var items []service.OrderItemInput
	for _, item := range req.Items {
		items = append(items, service.OrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), userID, items)
	if err != nil {
		if err == service.ErrProductUnavailable {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{"code": "UNPROCESSABLE", "message": "one or more products are unavailable"},
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_INPUT", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))
	role := c.GetString("role")

	orders, err := h.orderService.GetOrders(c.Request.Context(), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "something went wrong"},
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_INPUT", "message": "invalid order id"},
		})
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))
	role := c.GetString("role")

	order, err := h.orderService.GetByID(c.Request.Context(), id, userID, role)
	if err == repository.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "order not found"},
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "something went wrong"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_INPUT", "message": "invalid order id"},
		})
		return
	}

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_INPUT", "message": err.Error()},
		})
		return
	}

	order, err := h.orderService.UpdateStatus(c.Request.Context(), id, req.Status)
	if err == repository.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "order not found"},
		})
		return
	}
	if err == service.ErrInvalidTransition {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "UNPROCESSABLE", "message": "invalid status transition"},
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "something went wrong"},
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

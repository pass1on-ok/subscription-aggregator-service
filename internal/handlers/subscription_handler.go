package handlers

import (
	"net/http"
	"strconv"
	"time"

	"subscription-service/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	svc    services.SubscriptionService
	logger *zap.SugaredLogger
}

func NewSubscriptionHandler(svc services.SubscriptionService, logger *zap.SugaredLogger) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, logger: logger}
}

type createReq struct {
	ServiceName string  `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int     `json:"price"        binding:"required,min=0" example:"400"`
	UserID      string  `json:"user_id"      binding:"required,uuid"   example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string  `json:"start_date"   binding:"required"        example:"07-2025"`
	EndDate     *string `json:"end_date"                               example:"09-2025"`
}

type updateReq struct {
	ServiceName *string `json:"service_name" example:"Yandex Plus Premium"`
	Price       *int    `json:"price"        example:"500"`
	StartDate   *string `json:"start_date"   example:"08-2025"`
	EndDate     *string `json:"end_date"     example:"10-2025"`
}

// Create godoc
// @Summary      Create subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        body body createReq true "subscription"
// @Success      201 {object} models.Subscription
// @Failure      400 {object} map[string]string
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var in createReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub, err := h.svc.Create(c, services.CreateDTO{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      in.UserID,
		StartMonth:  in.StartDate,
		EndMonth:    in.EndDate,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

// Get godoc
// @Summary      Get subscription by ID
// @Tags         subscriptions
// @Produce      json
// @Param        id path string true "subscription id"
// @Success      200 {object} models.Subscription
// @Failure      404 {object} map[string]string
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHandler) Get(c *gin.Context) {
	id := c.Param("id")
	sub, err := h.svc.Get(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Update godoc
// @Summary      Update subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "subscription id"
// @Param        body body updateReq true "update"
// @Success      200 {object} models.Subscription
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var in updateReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub, err := h.svc.Update(c, id, services.UpdateDTO{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		StartMonth:  in.StartDate,
		EndMonth:    in.EndDate,
	})
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Delete godoc
// @Summary      Delete subscription
// @Tags         subscriptions
// @Param        id path string true "subscription id"
// @Success      200 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// List godoc
// @Summary      List subscriptions
// @Tags         subscriptions
// @Produce      json
// @Param        user_id query string false "UUID"
// @Param        service_name query string false "filter by service"
// @Param        from query string false "MM-YYYY"
// @Param        to query string false "MM-YYYY"
// @Param        limit query int false "limit"
// @Param        offset query int false "offset"
// @Success      200 {array} models.Subscription
// @Router       /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	var (
		userID      = c.Query("user_id")
		serviceName = c.Query("service_name")
		limit       = parseIntDefault(c.Query("limit"), 100)
		offset      = parseIntDefault(c.Query("offset"), 0)
	)

	var fromPtr, toPtr *time.Time
	if fromQ := c.Query("from"); fromQ != "" {
		if m, err := time.Parse("01-2006", fromQ); err == nil {
			m = time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.UTC)
			fromPtr = &m
		}
	}
	if toQ := c.Query("to"); toQ != "" {
		if m, err := time.Parse("01-2006", toQ); err == nil {
			m = time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.UTC)
			toPtr = &m
		}
	}

	items, err := h.svc.List(c, services.ListFilter{
		UserID:      userID,
		ServiceName: serviceName,
		From:        fromPtr,
		To:          toPtr,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Total godoc
// @Summary      Sum total price over period (by months overlap)
// @Tags         subscriptions
// @Produce      json
// @Param        user_id query string false "UUID"
// @Param        service_name query string false "filter by service"
// @Param        from query string true "MM-YYYY"
// @Param        to query string true "MM-YYYY"
// @Success      200 {object} map[string]int64
// @Failure      400 {object} map[string]string
// @Router       /subscriptions/total [get]
func (h *SubscriptionHandler) Total(c *gin.Context) {
	fromQ := c.Query("from")
	toQ := c.Query("to")
	if fromQ == "" || toQ == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required (MM-YYYY)"})
		return
	}

	from, err1 := time.Parse("01-2006", fromQ)
	to, err2 := time.Parse("01-2006", toQ)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected MM-YYYY"})
		return
	}

	total, err := h.svc.Total(c, services.TotalFilter{
		UserID:      c.Query("user_id"),
		ServiceName: c.Query("service_name"),
		From:        from,
		To:          to,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total})
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}

package http

import (
	"net/http"
	"strconv"

	"github.com/mihazzz123/m3zold-server/internal/usecase/device"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	CreateUC       *device.CreateUseCase
	ListUC         *device.ListUseCase
	FindUC         *device.FindUseCase
	UpdateStatusUC *device.UpdateStatusUseCase
	DeleteUC       *device.DeleteUseCase
}

func NewDeviceHandler(
	createUC *device.CreateUseCase,
	listUC *device.ListUseCase,
	findUC *device.FindUseCase,
	updateStatusUC *device.UpdateStatusUseCase,
	deleteUC *device.DeleteUseCase,
) *DeviceHandler {
	return &DeviceHandler{
		CreateUC:       createUC,
		ListUC:         listUC,
		FindUC:         findUC,
		UpdateStatusUC: updateStatusUC,
		DeleteUC:       deleteUC,
	}
}

// POST /devices
func (h *DeviceHandler) Create(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	var input device.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	input.UserID = userID

	if err := h.CreateUC.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Устройство добавлено"})
}

// GET /devices
func (h *DeviceHandler) List(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	devices, err := h.ListUC.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, devices)
}

// GET /devices/:id
func (h *DeviceHandler) Find(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	d, err := h.FindUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Устройство не найдено"})
		return
	}

	c.JSON(http.StatusOK, d)
}

// PATCH /devices/:id/status
func (h *DeviceHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный статус"})
		return
	}

	if err := h.UpdateStatusUC.Execute(c.Request.Context(), id, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус обновлён"})
}

// DELETE /devices/:id
func (h *DeviceHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	if err := h.DeleteUC.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Устройство удалено"})
}

// Временная функция — позже заменим на JWT
func extractUserID(c *gin.Context) (int, error) {
	return strconv.Atoi(c.GetHeader("X-User-ID"))
}

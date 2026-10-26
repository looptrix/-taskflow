package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"taskflow/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{service: taskService}
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type updateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

func currentUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	id, ok := userID.(uint)
	if !ok {
		return 0
	}

	return id
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req createTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные"})
		return
	}

	task, err := h.service.Create(
		c.Request.Context(),
		currentUserID(c),
		req.Title,
		req.Description,
		req.Priority,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	tasks, err := h.service.GetAll(c.Request.Context(), currentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить задачи"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID задачи"})
		return
	}

	task, err := h.service.GetByID(c.Request.Context(), currentUserID(c), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "задача не найдена"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить задачу"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID задачи"})
		return
	}

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные"})
		return
	}

	task, err := h.service.Update(
		c.Request.Context(),
		currentUserID(c),
		uint(id),
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
	)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "задача не найдена"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID задачи"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), currentUserID(c), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить задачу"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "задача удалена"})
}

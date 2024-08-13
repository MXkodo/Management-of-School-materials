package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MXkodo/Management-of-School-materials/internal/service"

	"github.com/MXkodo/Management-of-School-materials/model"
	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	service service.MaterialService
}

func NewMaterialHandler(service service.MaterialService) *MaterialHandler {
	return &MaterialHandler{service: service}
}

func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
	var mat model.Material
	if err := c.BindJSON(&mat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuid, err := h.service.CreateMaterial(c, &mat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"uuid": uuid})
}

func (h *MaterialHandler) GetMaterial(c *gin.Context) {
	uuid := c.Param("uuid")
	mat, err := h.service.GetMaterial(c, uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}
	c.JSON(http.StatusOK, mat)
}

func (h *MaterialHandler) UpdateMaterial(c *gin.Context) {
	uuid := c.Param("uuid")
	var mat model.Material
	if err := c.BindJSON(&mat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mat.UUID = uuid
	if err := h.service.UpdateMaterial(c, &mat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *MaterialHandler) GetAllMaterials(c *gin.Context) {
	filter := model.MaterialFilter{
		Type:   c.Query("type"),
		Status: c.Query("status"),
	}

	if from := c.Query("date_from"); from != "" {
		if dateFrom, err := time.Parse(time.RFC3339, from); err == nil {
			filter.DateFrom = dateFrom
		}
	}
	if to := c.Query("date_to"); to != "" {
		if dateTo, err := time.Parse(time.RFC3339, to); err == nil {
			filter.DateTo = dateTo
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	materials, err := h.service.GetAllMaterials(c, filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, materials)
}

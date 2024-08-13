package handler

import (
	"github.com/MXkodo/Management-of-School-materials/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, service service.MaterialService) {
	materialHandler := NewMaterialHandler(service)

	router.POST("/materials", materialHandler.CreateMaterial)
	router.GET("/materials/:uuid", materialHandler.GetMaterial)
	router.PUT("/materials/:uuid", materialHandler.UpdateMaterial)
	router.GET("/materials", materialHandler.GetAllMaterials)
}

package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/middlewares"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/models"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/service"
)

type Handler struct {
	storeCustomerService service.StoreCustomerService
}

func NewHandler(storeCustomerService *service.StoreCustomerService) *Handler {
	return &Handler{storeCustomerService: *storeCustomerService}
}

func (h *Handler) CreateStoreCustomer(ctx *gin.Context) {
	var storeCustomer models.StoreCustomer

	ctx.ShouldBindJSON(&storeCustomer)

	log.Printf("%+v", storeCustomer)

	createdStoreCustomer, err := h.storeCustomerService.CreateStoreCustomer(ctx, storeCustomer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdStoreCustomer)

}

func (h *Handler) GetStoreCustomerByID(ctx *gin.Context) {
	id := ctx.Param("id")

	storeCustomer, err := h.storeCustomerService.GetStoreCustomerByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, storeCustomer)
}

func (h *Handler) GetStoreCustomersByStoreID(ctx *gin.Context) {
	storeID := ctx.Param("id")

	storeCustomers, err := h.storeCustomerService.GetStoreCustomersByStoreID(ctx, storeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, storeCustomers)
}

func (h *Handler) GetStoreCustomersByUserID(ctx *gin.Context) {
	user := ctx.MustGet("user").(middlewares.User)
	log.Printf("userDataID: %s", user.ID)
	storeCustomers, err := h.storeCustomerService.GetStoreCustomersByUserID(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, storeCustomers)

}

func (h *Handler) UpdateStoreCustomer(ctx *gin.Context) {
	var storeCustomer models.StoreCustomer
	if err := ctx.BindJSON(&storeCustomer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedStoreCustomer, err := h.storeCustomerService.UpdateStoreCustomer(ctx, storeCustomer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedStoreCustomer)
}

func (h *Handler) DeleteStoreCustomer(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.storeCustomerService.DeleteStoreCustomer(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "StoreCustomer deleted successfully"})
}

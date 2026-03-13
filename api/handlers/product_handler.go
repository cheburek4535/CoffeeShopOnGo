package handlers

import (
	"CoffeeShopOnGo/api/models"
	"CoffeeShopOnGo/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.Service
}

func NewProductHandler(svc *service.Service) *ProductHandler {
	return &ProductHandler{service: svc}
}

// @Success      200  {array}   models.ProductResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/products/ [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.service.ShowAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	response := make([]models.ProductResponse, len(products))
	for i, p := range products {
		response[i] = models.ProductResponse{
			ID:       p.ID,
			Name:     p.Name,
			Price:    p.Price,
			VAT:      p.VAT,
			Category: p.Category,
		}
	}
	c.JSON(http.StatusOK, response)
}

// @Param        id   path      int  true  "ID товара"
// @Success      200  {object}  models.ProductResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /api/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid id"})
		return
	}
	product, err := h.service.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Товар не найден"})
		return
	}
	response := models.ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		VAT:      product.VAT,
		Category: product.Category,
		Price:    product.Price,
	}
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	product, err := h.service.AddNewProduct(
		req.Name,
		req.Price,
		req.Category,
		req.VAT,
		req.IsActive,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	response := models.ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
		VAT:      product.VAT,
	}
	c.JSON(http.StatusCreated, response)
}

// @Success      200  {object}  models.ProductResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      402 {object}  models.ErrorResponse
// @Router       /api/products/ [post]
func (h *ProductHandler) UpdatePrice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid id"})
		return
	}

	var req models.UpdatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	product, err := h.service.ChangePrice(id, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	response := models.ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
		VAT:      product.VAT,
	}
	c.JSON(http.StatusOK, response)
}

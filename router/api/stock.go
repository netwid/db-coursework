package api

import (
	"github.com/gin-gonic/gin"
	"github.com/netwid/db-coursework/repository"
	"net/http"
	"strconv"
)

type StockApi struct {
	stockRepo repository.StockRepository
}

func NewStockApi(stockRepo repository.StockRepository) *StockApi {
	return &StockApi{stockRepo: stockRepo}
}

// GetCategories @Summary get categories
// @Tags stock
// @Produce json
// @Success 200 {object} []repository.Category
// @Router /categories [get]
func (s *StockApi) GetCategories(c *gin.Context) {
	categories, err := s.stockRepo.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetStocks @Summary get stocks
// @Tags stock
// @Produce json
// @Success 200 {object} []repository.Stock
// @Router /stocks [get]
func (s *StockApi) GetStocks(c *gin.Context) {
	stocks, err := s.stockRepo.GetStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

// GetStock @Summary get stock
// @Tags stock
// @Param id path int true "Stock ID"
// @Produce json
// @Success 200 {object} repository.FullStock
// @Router /stocks/{id} [get]
func (s *StockApi) GetStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stock, err := s.stockRepo.GetStock(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// GetPrice @Summary get price
// @Tags stock
// @Param id path int true "Stock ID"
// @Produce json
// @Success 200 {object} []repository.Price
// @Router /stocks/{id}/price [get]
func (s *StockApi) GetPrice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	price, err := s.stockRepo.GetPrice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, price)
}

type buy struct {
	StockId int `json:"stock_id"`
	Amount  int `json:"amount"`
}

// Buy @Summary buy
// @Tags stock
// @Produce json
// @Param data body buy true "Data JSON Object"
// @Security ApiKeyAuth
// @Success 200 {string} Ok
// @Router /buy [post]
func (s *StockApi) Buy(c *gin.Context) {
	id, _ := c.Get("id")

	var data buy

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.stockRepo.Buy(id.(int), data.StockId, data.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

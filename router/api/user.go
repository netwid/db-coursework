package api

import (
	"github.com/gin-gonic/gin"
	"github.com/netwid/db-coursework/repository"
	util "github.com/netwid/db-coursework/utils"
	"net/http"
)

type auth struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserApi struct {
	userRepo repository.UserRepository
}

func NewUserApi(userRepo repository.UserRepository) *UserApi {
	return &UserApi{userRepo: userRepo}
}

// Refresh @Summary crutch for avoid redux
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} Test
// @Router /refresh [get]
func Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": c.GetHeader("Authorization")})
}

// Register @Summary register
// @Tags auth
// @Accept json
// @Produce json
// @Param data body auth true "Data JSON Object"
// @Success 200 {string} Test
// @Router /register [post]
func (u *UserApi) Register(c *gin.Context) {
	var data auth

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := &repository.User{Email: data.Email, Password: data.Password}
	err := u.userRepo.Create(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Login @Summary login
// @Tags auth
// @Accept json
// @Produce json
// @Param data body auth true "Data JSON Object"
// @Success 200 {string} Test
// @Router /login [post]
func (u *UserApi) Login(c *gin.Context) {
	var data auth

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := u.userRepo.GetId(data.Email, data.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := util.GenerateToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

type ticket struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreateTicket @Summary ticket
// @Tags ticket
// @Accept json
// @Produce json
// @Param data body ticket true "Data JSON Object"
// @Security ApiKeyAuth
// @Success 200 {string} Ok
// @Router /ticket [post]
func (u *UserApi) CreateTicket(c *gin.Context) {
	id, _ := c.Get("id")

	var data ticket

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.userRepo.CreateTicket(id.(int), data.Title, data.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// GetProfile @Summary Get Profile
// @Tags profile
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} repository.Profile
// @Router /profile [get]
func (u *UserApi) GetProfile(c *gin.Context) {
	id, _ := c.Get("id")

	profile, err := u.userRepo.GetProfile(id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile @Summary Update Profile
// @Tags profile
// @Produce json
// @Security ApiKeyAuth
// @Param data body repository.Profile true "Data JSON Object"
// @Success 200 {string} Ok
// @Router /profile [put]
func (u *UserApi) UpdateProfile(c *gin.Context) {
	id, _ := c.Get("id")

	var data repository.Profile

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.userRepo.UpdateProfile(id.(int), &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// GetPortfolio @Summary Get Portfolio
// @Tags portfolio
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []repository.PortfolioItem
// @Router /portfolio [get]
func (u *UserApi) GetPortfolio(c *gin.Context) {
	id, _ := c.Get("id")

	portfolio, err := u.userRepo.GetPortfolio(id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

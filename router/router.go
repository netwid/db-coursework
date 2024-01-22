package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/netwid/db-coursework/docs"
	"github.com/netwid/db-coursework/middlewares/auth"
	"github.com/netwid/db-coursework/repository"
	"github.com/netwid/db-coursework/router/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(repos repository.Repositories) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userApi := api.NewUserApi(repos.UserRepo)
	stockApi := api.NewStockApi(repos.StockRepo)

	publicApi := r.Group("/api")
	{
		publicApi.POST("/register", userApi.Register)
		publicApi.POST("/login", userApi.Login)

		publicApi.GET("/categories", stockApi.GetCategories)
		publicApi.GET("/stocks", stockApi.GetStocks)
		publicApi.GET("/stocks/:id", stockApi.GetStock)
		publicApi.GET("/stocks/:id/price", stockApi.GetPrice)
	}

	privateApi := r.Group("/api")
	privateApi.Use(auth.JWTAuthMiddleware())
	{
		privateApi.GET("/refresh", api.Refresh)

		privateApi.GET("/profile", userApi.GetProfile)
		privateApi.PUT("/profile", userApi.UpdateProfile)

		privateApi.POST("/ticket", userApi.CreateTicket)

		privateApi.POST("/buy", stockApi.Buy)

		privateApi.GET("/portfolio", userApi.GetPortfolio)
	}

	return r
}

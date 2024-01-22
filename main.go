package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netwid/db-coursework/repository"
	"github.com/netwid/db-coursework/router"
	"log"
)

func init() {

}

// @title           TradePulse API
// @version         1.0
// @description     TradePulse API

// @license.name  Apache 2.0
// @license.url   https://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	connStr := "postgres://dev:dev@localhost:5432/postgres"
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	repos := repository.Repositories{
		UserRepo:  repository.NewUserRepo(dbpool),
		StockRepo: repository.NewStockRepo(dbpool),
	}

	r := router.NewRouter(repos)
	r.Run(":8080")
}

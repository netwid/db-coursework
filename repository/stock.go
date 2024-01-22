package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type StockRepository interface {
	GetCategories() ([]Category, error)
	GetStocks() ([]Stock, error)
	GetStock(id int) (*FullStock, error)
	GetPrice(id int) ([]Price, error)
	Buy(userId, stockId, amount int) error
}

type stockRepo struct {
	db *pgxpool.Pool
}

func NewStockRepo(db *pgxpool.Pool) StockRepository {
	return &stockRepo{
		db: db,
	}
}

const (
	categoryTable = "stock_category"
)

type Category struct {
	Id   int
	Name string
}

func (s *stockRepo) GetCategories() ([]Category, error) {
	var categories []Category

	sql, args, err := mysq.
		Select("id", "name").
		From(categoryTable).
		Where(sq.Eq{"stock_category_parent_id": 1}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var value Category
		if err := rows.Scan(&value.Id, &value.Name); err != nil {
			return nil, err
		}
		categories = append(categories, value)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

type Stock struct {
	Id         int
	Name       string
	CategoryID int `json:"category_id"`
}

func (s *stockRepo) GetStocks() ([]Stock, error) {
	var stocks []Stock

	sql, args, err := mysq.
		Select("id", "name", "category_id").
		From("stock").
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var value Stock
		if err := rows.Scan(&value.Id, &value.Name, &value.CategoryID); err != nil {
			return nil, err
		}
		stocks = append(stocks, value)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}

type FullStock struct {
	Id          int
	Name        string
	Description *string
	Currency    *string
	CategoryID  *int `json:"category_id" db:"category_id"`
}

func (s *stockRepo) GetStock(id int) (*FullStock, error) {
	var stock FullStock

	sql, args, err := mysq.
		Select("stock.id", "stock.name", "description", "currency.name AS currency", "category_id").
		From("stock").
		Join("currency ON stock.currency_id = currency.id").
		Where(sq.Eq{"stock.id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(context.Background(), sql, args...).Scan(&stock.Id, &stock.Name, &stock.Description, &stock.Currency, &stock.CategoryID)
	if err != nil {
		return nil, err
	}

	return &stock, nil
}

type Price struct {
	Date  time.Time
	Price int
}

func (s *stockRepo) GetPrice(id int) ([]Price, error) {
	var prices []Price

	sql, args, err := mysq.
		Select("date", "price").
		From("stock_price").
		Where(sq.Eq{"stock_id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var value Price
		if err := rows.Scan(&value.Date, &value.Price); err != nil {
			return nil, err
		}
		prices = append(prices, value)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

func (s *stockRepo) Buy(userId, stockId, amount int) error {
	sql, args, err := mysq.
		Select("id").
		From("portfolio").
		Where(sq.Eq{"user_id": userId}).
		ToSql()

	if err != nil {
		return err
	}

	var portfolioId int
	err = s.db.QueryRow(context.Background(), sql, args...).Scan(&portfolioId)
	if err != nil {
		return err
	}

	sql, args, err = mysq.Insert("portfolio_item").
		Columns("portfolio_id", "stock_id", "amount").
		Values(portfolioId, stockId, amount).
		Suffix("ON CONFLICT (portfolio_id, stock_id) DO UPDATE SET amount = portfolio_item.amount + ?", amount).
		ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

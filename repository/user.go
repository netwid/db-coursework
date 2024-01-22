package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Email    string
	Password string
}

type UserRepository interface {
	Create(user *User) error
	GetId(email, password string) (int, error)
	CreateTicket(userId int, title, content string) error
	GetProfile(userId int) (*Profile, error)
	UpdateProfile(userId int, profile *Profile) error
	GetPortfolio(userId int) ([]PortfolioItem, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepository {
	return &userRepo{
		db: db,
	}
}

const (
	userTable = "\"user\""
)

func (r *userRepo) Create(user *User) error {
	sql, args, err := mysq.Insert(userTable).
		Columns("email", "password").
		Values(user.Email, user.Password).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return err
	}

	var userId int
	err = r.db.QueryRow(context.Background(), sql, args...).Scan(&userId)
	if err != nil {
		return err
	}

	sql, args, err = mysq.Insert("portfolio").
		Columns("user_id", "currency_id", "name").
		Values(userId, 1, "Default").
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetId(email, password string) (int, error) {
	var id int

	sql, args, err := mysq.
		Select("id").
		From(userTable).
		Where(sq.Eq{"email": email, "password": password}).
		ToSql()

	if err != nil {
		return 0, err
	}

	err = r.db.QueryRow(context.Background(), sql, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userRepo) CreateTicket(userId int, title, content string) error {
	sql, args, err := mysq.Insert("ticket").
		Columns("user_id", "title", "content").
		Values(userId, title, content).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

type Profile struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
	Email   *string `json:"email"`
	Age     *int    `json:"age"`
	Country *string `json:"country"`
	Phone   *string `json:"phone"`
}

func (r *userRepo) GetProfile(userId int) (*Profile, error) {
	sql, args, err := mysq.
		Select("name", "surname", "email", "age", "country", "phone").
		From(userTable).
		Where(sq.Eq{"id": userId}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var profile Profile
	err = r.db.QueryRow(context.Background(), sql, args...).Scan(&profile.Name, &profile.Surname, &profile.Email, &profile.Age, &profile.Country, &profile.Phone)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *userRepo) UpdateProfile(userId int, profile *Profile) error {
	sql, args, err := mysq.
		Update(userTable).
		Set("name", profile.Name).
		Set("surname", profile.Surname).
		Set("email", profile.Email).
		Set("age", profile.Age).
		Set("country", profile.Country).
		Set("phone", profile.Phone).
		Where(sq.Eq{"id": userId}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

type PortfolioItem struct {
	Category   string `json:"category"`
	CategoryId int    `json:"category_id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Amount     int    `json:"amount"`
}

func (r *userRepo) GetPortfolio(userId int) ([]PortfolioItem, error) {
	sql, args, err := mysq.
		Select("stock_category.name", "stock_category.id", "stock.name", "latest_price.price", "portfolio_item.amount").
		From("portfolio_item").
		Join("portfolio ON portfolio.id = portfolio_item.portfolio_id").
		Join("stock ON stock.id = portfolio_item.stock_id").
		Join("stock_category ON stock_category.id = stock.category_id").
		JoinClause(
			// подзапрос для получения последней цены каждого стока
			"JOIN (SELECT stock_id, price FROM stock_price WHERE (stock_id, date) IN (" +
				"SELECT stock_id, MAX(date) FROM stock_price GROUP BY stock_id" +
				")) AS latest_price ON latest_price.stock_id = stock.id",
		).
		Where(sq.Eq{"portfolio.user_id": userId}).
		ToSql()

	var items []PortfolioItem

	rows, err := r.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var value PortfolioItem
		if err := rows.Scan(&value.Category, &value.CategoryId, &value.Name, &value.Price, &value.Amount); err != nil {
			return nil, err
		}
		items = append(items, value)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

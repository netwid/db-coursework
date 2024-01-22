package repository

import sq "github.com/Masterminds/squirrel"

type Repositories struct {
	UserRepo  UserRepository
	StockRepo StockRepository
}

var (
	mysq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

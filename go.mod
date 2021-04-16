module github.com/matthewpi/pgtype

go 1.15

require (
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/jackc/pgio v1.0.0
	github.com/lib/pq v1.3.0
	github.com/matthewpi/pgconn v1.8.2
	github.com/matthewpi/pgx/v4 v4.11.2
	github.com/shopspring/decimal v0.0.0-20200227202807-02e2044944cc
	github.com/stretchr/testify v1.5.1
)

replace github.com/matthewpi/pgx/v4 => ../pgx

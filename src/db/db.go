package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	. "go-api-server/common"
)

const pg_db_url = "postgres://pg:pg@db-server-01:5432/test01"

var Pool *pgxpool.Pool

func Connect() {
	var err error
	connection_url := Default_string(os.Getenv("PG_DB_URL"), pg_db_url)
	Pool, err = pgxpool.New(context.Background(), connection_url)
	Panic(err)
	Log("Initialized database", pg_db_url)
}

func Close() {
	Pool.Close()
}

func Next_Id(ctx context.Context, sequence string) (int64, error) {
	var id int64
	query := fmt.Sprintf(`SELECT nextval('%s');`, sequence)
	row := Pool.QueryRow(ctx, query)
	err := row.Scan(&id)
	return id, err
}

package adapter

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // driver postgres
)

var connectionDB *sqlx.DB

func LoadDatabase() {
	dsn := os.Getenv("DB_URL")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		err = fmt.Errorf("error to connection db %s. %v", dsn, err)
		panic(err)
	}

	// set max idle and connection to database
	maxIdle, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE"))
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTION"))
	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxConn)

	err = db.Ping()

	connectionDB = db
}

func DBConnection() *sqlx.DB {
	return connectionDB
}

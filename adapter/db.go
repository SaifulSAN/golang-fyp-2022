package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	// "github.com/jackc/pgx/v5/stdlib"
	// "github.com/joho/godotenv"
)

func NewDbConnection() *sql.DB {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := url.URL{
		Scheme: "postgres",
		Host:   os.Getenv("dbHost"),
		User:   url.UserPassword(os.Getenv("dbUser"), os.Getenv("dbPassword")),
		Path:   os.Getenv("dbName"),
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		fmt.Println("sql.Open", err)
	}

	// defer func() {
	// 	_ = db.Close()
	// 	fmt.Println("Closed")
	// }()

	return db

	// row := db.QueryRowContext(context.Background(), "SELECT test_some_number FROM test_pgx_conn WHERE test_name = 'hellosaiful'")
	// if err := row.Err(); err != nil {
	// 	fmt.Println("db.QueryRowContext", err)
	// 	return
	// }

	// var someNumber int

	// if err := row.Scan(&someNumber); err != nil {
	// 	fmt.Println("row.Scan", err)
	// 	return
	// }

	// fmt.Println("someNumber", someNumber)
}

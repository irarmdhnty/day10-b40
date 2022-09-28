package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func ConnectDatabase() {

	// postgres://postgres:password@localhost:5432/database_name
	urlDatabase := "postgres://postgres:12345@localhost:5432/my_project"

	var err error
	Conn, err = pgx.Connect(context.Background(), urlDatabase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success connect to database")
}
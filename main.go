package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

type dbInfo struct {
	port     string
	host     string
	name     string
	user     string
	password string
}

func (db *dbInfo) getDbInfo() string {
	info := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.host,
		db.port,
		db.user,
		db.password,
		db.name,
	)

	return info
}

func NewDbInfo() *dbInfo {
	var (
		host     = os.Getenv("DATABASE_HOST")
		port     = os.Getenv("DATABASE_PORT")
		user     = os.Getenv("DATABASE_USER")
		password = os.Getenv("DATABASE_PASSWORD")
		name     = os.Getenv("DATABASE_NAME")
	)

	return &dbInfo{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		name:     name,
	}
}

func initializeDB() (db *sql.DB, err error) {
	dbInformation := NewDbInfo().getDbInfo()
	fmt.Println(dbInformation)
	db, err = sql.Open(os.Getenv("DATABASE_DRIVER"), dbInformation)

	if err != nil {
		log.Printf("Cannot connect to %s database", dbInformation)
		log.Fatal("This is the error: ", err)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("Pinging: ", err)
		return nil, err
	}

	log.Println("Database Connection established")

	return db, nil
}

func main() {
	time.Sleep(7 * time.Second)
	log.Println("Successfully started!")

	database, err := initializeDB()

	if err != nil {
		log.Fatalf("Could not set up database %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS test (id int primary key auto_increment, num int)`,
	)

	defer database.ExecContext(ctx, "DROP TABLE test;")

	if err != nil {
		log.Fatal("Create table err: ", err)
	}

	_, err = database.ExecContext(
		ctx,
		`INSERT INTO test(num) VALUES (?);`,
		800,
	)

	if err != nil {
		log.Fatal("Insert err: ", err)
	}

	data, err := database.QueryContext(
		ctx,
		`SELECT * FROM test;`,
	)

	if err != nil {
		log.Fatal("Select err: ", err)
	}

	var id, num int

	for data.Next() {
		err = data.Scan(&id, &num)
		fmt.Println(err, id, num)
	}

	log.Println(data.Close())

}

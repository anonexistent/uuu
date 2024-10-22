package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	id   int
	name string
}

func (User) TableName() string {
	return "users"
}

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка при чтении файла конфигурации: %v", err)
	}

	var dbPath = viper.Get("TEST_DB_PATH")
	if dbPath == "" {
		panic("TEST_DB_PATH is empty")
	}

	db, err := gorm.Open(postgres.Open(dbPath.(string)), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect database: \n%v\n%v", dbPath, err)
	}

	user, err := GetUserNew(db, 1)
	if err != nil {
		panic(err)
	}
	PrintJson(user)
	////	native
	// conn, err := pgx.Connect(context.Background(), dbPath.(string))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close(context.Background())

	// fmt.Print(RemoveUser(conn, 10))
}

func PrintJson(v any) {
	bytes, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(bytes))
}

func GetUserNew(db *gorm.DB, userId int) (User, error) {
	var user User
	err := db.Take(&user, userId).Error
	return user, err
}

func GetUser(con *pgx.Conn, userId int) (User, error) {
	var user User
	var row = con.QueryRow(
		context.Background(),
		"select id, name from users where id=$1",
		userId,
	)
	err := row.Scan(
		&user.id,
		&user.name,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserList(con *pgx.Conn) ([]User, error) {
	var us = make([]User, 0)
	rows, err := con.Query(
		context.Background(),
		"select * from users",
	)
	if err != nil {
		return us, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.id,
			&user.name,
		)
		if err != nil {
			return []User{}, err
		}
		us = append(us, user)
	}
	return us, nil
}

func CreateUser(con *pgx.Conn, user User) (int, error) {
	var id int
	err := con.QueryRow(
		context.Background(),
		"insert into users values($1,$2) returning id",
		user.id, user.name,
	).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func RemoveUser(con *pgx.Conn, userId int) (int, error) {
	tag, err := con.Exec(
		context.Background(),
		"delete from users where id=$1",
		userId,
	)
	if err != nil {
		return -1, err
	}
	//	count of edited rows
	return int(tag.RowsAffected()), nil
}

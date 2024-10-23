package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MyLogger = logger.New(
	log.New(os.Stderr, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold:             time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	},
)

type User struct {
	Id    uint
	Name  string
	Email string
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

	db, err := gorm.Open(postgres.Open(dbPath.(string)), &gorm.Config{Logger: MyLogger})
	if err != nil {
		log.Panicf("failed to connect database: \n%v\n%v", dbPath, err)
	}

	u, err := GetUserByEmail(db, "test@mail")
	if err != nil {
		os.Exit(-1)
	}
	PrintJson(u)

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

/*
// Get the first record ordered by primary key
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// Get one record, no specified order
db.Take(&user)
// SELECT * FROM users LIMIT 1;

// Get last record, ordered by primary key desc
db.Last(&user)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;
*/
func GetUserByEmail(db *gorm.DB, name string) ([]User, error) {
	us := make([]User, 0)
	err := db.Where("email LIKE ?", "%"+name+"%").Find(&us).Error
	return us, err
}

func GetUserByName(db *gorm.DB, name string) (User, error) {
	var user User
	err := db.Where("name = ?", name).Take(&user).Error
	return user, err
}

func GetUserNew(db *gorm.DB, userId int) (User, error) {
	var user User
	err := db.Take(&user, userId).Error
	return user, err
}

func GetUserListNew(db *gorm.DB) ([]User, error) {
	us := make([]User, 0)
	// err:= db.Preload("Images").Find(&us).Error
	err := db.Find(&us).Error
	return us, err
}

func CreateUserNew(db *gorm.DB, user User) (uint, error) {
	err := db.Create(&user).Error
	return user.Id, err
}

func RemoveUserNew(db *gorm.DB, id uint) (int, error) {
	command := db.Delete(&User{}, id)
	return int(command.RowsAffected), command.Error
}

func GetUser(con *pgx.Conn, userId int) (User, error) {
	var user User
	var row = con.QueryRow(
		context.Background(),
		"select id, name from users where id=$1",
		userId,
	)
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
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
			&user.Id,
			&user.Name,
			&user.Email,
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
		"insert into users values($1,$2,$3) returning id",
		user.Id, user.Name, user.Email,
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

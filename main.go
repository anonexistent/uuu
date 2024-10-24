package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"uuu/config"
	"uuu/lib"
	"uuu/sqlWork"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var dbPath = config.Load().DbPath

	db, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{Logger: lib.MyLogger})
	if err != nil {
		log.Panicf("failed to connect database: \n%v\n%v", dbPath, err)
	}

	u, err := sqlWork.GetUserByEmail(db, "test@mail")
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

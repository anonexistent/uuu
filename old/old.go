package old

import (
	"context"
	"uuu/lib"

	"github.com/jackc/pgx/v5"
)

func GetUser(con *pgx.Conn, userId int) (lib.User, error) {
	var user lib.User
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

func GetUserList(con *pgx.Conn) ([]lib.User, error) {
	var us = make([]lib.User, 0)
	rows, err := con.Query(
		context.Background(),
		"select * from users",
	)
	if err != nil {
		return us, err
	}
	defer rows.Close()

	for rows.Next() {
		var user lib.User
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
		)
		if err != nil {
			return []lib.User{}, err
		}
		us = append(us, user)
	}
	return us, nil
}

func CreateUser(con *pgx.Conn, user lib.User) (int, error) {
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

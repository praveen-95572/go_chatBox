package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Users struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"u_name"`
}

type Chat struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"s_id"`
	RecieverID int    `json:"r_id"`
	Message    string `json:"msg"`
}

// ADD USER
func (m *DBModel) AddUser(name, username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO USERS(name, u_name) VALUES (?, ?)`

	res, err := m.DB.ExecContext(ctx, stmt,
		name,
		username,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

// GET ALL USER
func (m *DBModel) GetAllUser() ([]*Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*Users

	rows, err := m.DB.QueryContext(ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u Users
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Username,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

// GET ALL MSG of USER
func (m *DBModel) GetAllMsg(id, uid int) ([]*Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var chats []*Chat

	rows, err := m.DB.QueryContext(ctx, `SELECT * FROM chat WHERE (s_id = ? AND r_id = ?)   OR (s_id = ? AND r_id = ?)
	order by id;
	`, id, uid, uid, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Chat
		err := rows.Scan(
			&c.ID,
			&c.SenderID,
			&c.RecieverID,
			&c.Message,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, &c)
	}
	return chats, nil
}

func (m *DBModel) PostMsg(id, chat_id int, msg string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO chat (s_id, r_id, msg) VALUES (?, ?, ?)`

	res, err := m.DB.ExecContext(ctx, stmt, id, chat_id, msg)
	if err != nil {
		return 0, err
	}

	idx, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(idx), nil
}

func (m *DBModel) GetMsg(id int) ([]*Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var chats []*Chat

	rows, err := m.DB.QueryContext(ctx, `SELECT * FROM chat WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Chat
		err := rows.Scan(
			&c.ID,
			&c.SenderID,
			&c.RecieverID,
			&c.Message,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, &c)
	}
	return chats, nil
}

func (m *DBModel) AddNewUser(name, uname string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO users(name, u_name) VALUES(?, ?)`

	res, err := m.DB.ExecContext(ctx, stmt, name, uname)
	if err != nil {
		return 0, err
	}

	idx, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(idx), nil
}

func (m *DBModel) GetAllFriend(id int) ([]*Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select DISTINCT users.* from users INNER JOIN (SELECT * from chat where s_id = ? OR r_id = ?) AS tags ON tags.s_id = users.id OR tags.r_id = users.id`

	var users []*Users

	rows, err := m.DB.QueryContext(ctx, stmt, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u Users
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Username,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	fmt.Println("USERS IS ", users)
	return users, nil

}

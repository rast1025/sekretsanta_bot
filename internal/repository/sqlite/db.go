package sqlite

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/rast1025/sekretsanta_bot/internal/models"
)

type DB struct {
	client *sql.DB
}

func NewDB(client *sql.DB) *DB {
	return &DB{client: client}
}

func (d *DB) CreateUser(chatID, username string) error {
	insertUserQuery := `INSERT INTO users (chat_id, username) VALUES (?, ?)`
	_, err := d.client.Exec(insertUserQuery, chatID, username)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) UserExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ?`
	var count int
	err := d.client.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *DB) GetUser(username string) (*models.User, error) {
	query := `SELECT * FROM users WHERE username = ?`

	var user models.User
	err := d.client.QueryRow(query, username).Scan(&user.ID, &user.ChatID, &user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *DB) AddUser(groupID string, userID int) error {
	insertUserQuery := `INSERT INTO groups (user_id, group_id) VALUES (?, ?)`
	_, err := d.client.Exec(insertUserQuery, userID, groupID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && !errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return err
		}
	}

	return nil
}

func (d *DB) GetUsersFromGroup(groupID string) ([]models.User, error) {
	query := `SELECT users.id, users.chat_id, users.username FROM groups JOIN users ON users.id=groups.user_id WHERE groups.group_id = ?`
	rows, err := d.client.Query(query, groupID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.ChatID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil

}

package storage

import (
	"NewsFeed/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (ps *PostgresStorage) AddNewUser(user models.User) error {
	query := `INSERT INTO users (id, number) VALUES ($1, $2)`
	_, err := ps.db.Exec(query, user.ID, user.Number)
	return err
}

func (ps *PostgresStorage) RefreshToken(refreshToken string) (uint32, error) {
	var userID uint32
	log.Println("postgres RefreshToken " + refreshToken)
	err := ps.db.Get(&userID, "SELECT id FROM refresh WHERE refresh_token=$1", refreshToken)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (ps *PostgresStorage) AddRefreshToken(refreshToken string, id uint32) error {
	query := `INSERT INTO refresh (refresh_token, id) VALUES ($1, $2)`
	_, err := ps.db.Exec(query, refreshToken, id)
	return err
}

func (ps *PostgresStorage) UpdateRefreshToken(refreshToken string, id uint32) error {
	query := `UPDATE refresh SET refresh_token=$1 WHERE id=$2`
	_, err := ps.db.Exec(query, refreshToken, id)
	return err
}

func (ps *PostgresStorage) FinishAuth(userID uint32, nickname string) error {
	query := `UPDATE users SET nickname=$1 WHERE id=$2`
	_, err := ps.db.Exec(query, nickname, userID)
	return err
}

func (ps *PostgresStorage) GetUserByID(userID uint32) (models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id=$1`
	err := ps.db.Get(&user, query, userID)
	return user, err
}

// Feed
func (ps *PostgresStorage) GetUserFeedChannels(userID uint32) ([]models.Channel, error) {
	var channelIDs []string
	var channels []models.Channel
	query := `SELECT channel_id FROM subscriptions WHERE user_id=$1`
	err := ps.db.Select(&channelIDs, query, userID)
	if err != nil {
		return nil, err
	}
	query = `SELECT * FROM channels WHERE channel_id= ANY($1)`
	err = ps.db.Select(&channels, query, pq.Array(channelIDs))
	return channels, err
}
func (ps *PostgresStorage) AddFeedChannel(userID uint32, channelTag string, channelName string) error {
	query := `INSERT INTO subscriptions (user_id, channel_id) VALUES ($1, $2)`
	_, err := ps.db.Exec(query, userID, channelTag)
	if err != nil {
		return err
	}
	query = `INSERT INTO channels (channel_id, channel_name) VALUES ($1, $2)`
	_, err = ps.db.Exec(query, channelTag, channelName)
	return err
}
func (ps *PostgresStorage) RemoveFeedChannel(userID uint32, channelTag string) error {
	query := `DELETE FROM subscriptions WHERE user_id=$1 AND channel_id=$2`
	_, err := ps.db.Exec(query, userID, channelTag)
	if err != nil {
		return err
	}
	query = `DELETE FROM channels WHERE channel_id=$1`
	_, err = ps.db.Exec(query, channelTag)
	return err
}

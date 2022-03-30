package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zavyalov-den/url-shortener/internal/config"
)

type DB struct {
	db *pgxpool.Pool
}

func (d *DB) GetURL(short string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var fullURL string
	//language=sql
	query := `
		select full_url from urls where short_url = $1 limit 1;
	`
	err := d.db.QueryRow(ctx, query, short).Scan(&fullURL)
	if err != nil {
		return "", err
	}

	return fullURL, nil
}

func (d *DB) GetUserURLs(id int) []UserURL {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var result []UserURL

	//language=sql
	query := `
		select urls.short_url, urls.full_url from urls 
		join user_urls u on urls.id = u.url_id
		where u.user_id = $1; 
	`
	rows, err := d.db.Query(ctx, query, id)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if rows.Next() {
		fmt.Println(rows)
		values, err := rows.Values()
		if err != nil {
			return nil
		}
		result = append(result, UserURL{
			ShortURL:    values[0].(string),
			OriginalURL: values[1].(string),
		})
	}

	return result
}

func (d *DB) SaveURL(userID int, url UserURL) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var urlID int

	//language=sql
	query := `
		SELECT id from urls where short_url = $1;
	`
	err := d.db.QueryRow(ctx, query, url.ShortURL).Scan(&urlID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// it's okay. happens :)
		} else {
			return fmt.Errorf("select from urls failed: %s", err)
		}
	}
	fmt.Println(urlID)
	if urlID == 0 {
		//language=sql
		query = `
		insert into urls (short_url, full_url) VALUES ($1, $2);
		`

		_, err = d.db.Query(ctx, query, url.ShortURL, url.OriginalURL)
		if err != nil {
			return fmt.Errorf("insert into user_urls err: %s", err)
		}
	}
	//language=sql
	query = `
		insert into user_urls (url_id, user_id) VALUES ($1, $2);
		`

	_, err = d.db.Query(ctx, query, urlID, userID)
	if err != nil {
		return fmt.Errorf("insert into user_urls err: %s", err)
	}

	return nil
}

func (d *DB) SaveBatch(b []BatchRequest) ([]BatchResponse, error) {
	return nil, fmt.Errorf("no batches for in memory storage yet")
}

func (d *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := d.db.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewDB() *DB {
	cfg, err := pgxpool.ParseConfig(config.Config.DatabaseDSN)
	if err != nil {
		panic("failed to init db")
	}

	db, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		panic("db connection failed")
	}

	return &DB{db: db}
}

func (d *DB) InitDB() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// language=sql
	queries := []string{`
		CREATE TABLE if not exists urls (
			id serial primary key,
			short_url text not null unique,
			full_url text not null,
			correlation_id text
		);
		`, `
		CREATE TABLE if not exists user_urls (
		    user_id int, -- references users.id
		    url_id int
		);
-- 
-- 		CREATE TABLE users (
-- 		    id serial primary key,
-- 		    token text
-- 		)
	`}

	for _, query := range queries {
		_, err := d.db.Query(ctx, query)
		if err != nil {
			fmt.Println(err)
		}
	}
}

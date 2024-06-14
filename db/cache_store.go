package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CacheStore interface {
	Truncater
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key string, value string, validityDuration time.Duration) error
	ColelctGarbage(context.Context) error
}

type PostgresCacheStore struct {
	db    *sql.DB
	table string
}

func NewPostgresCacheStore(db *sql.DB, table string) *PostgresCacheStore {
	return &PostgresCacheStore{
		db:    db,
		table: table,
	}
}

func (s *PostgresCacheStore) Truncate(ctx context.Context) error {
	query := fmt.Sprintf("DELETE FROM %s", s.table)
	fmt.Println(query)
	s.db.Exec(query)
	return nil
}

func (s *PostgresCacheStore) Get(ctx context.Context, key string) string {
	cacheValue := ""
	err := s.db.QueryRow("SELECT cache_value FROM cache WHERE cache_key=$1 AND valid_until>$2", key, time.Now()).Scan(&cacheValue)
	if err != nil {
		return ""
	}

	return cacheValue
}

func (s *PostgresCacheStore) Set(ctx context.Context, key string, value string, validityDuration time.Duration) error {
	validUntil := time.Now().Add(validityDuration)
	_, err := s.db.Exec("INSERT INTO cache (created_at, valid_until, cache_key, cache_value) VALUES ($1, $2, $3, $4)", time.Now(), validUntil, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresCacheStore) ColelctGarbage(ctx context.Context) error {
	_, err := s.db.Exec("DELETE FROM cache WHERE valid_until<$1", time.Now())
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"space-bounty-hunters/internal/api"
)

func main() {
	ctx := context.Background()

	pool, err := connect(ctx, databaseURL())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	addr := ":" + getenv("PORT", "8080")
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, api.NewRouter(pool)))
}

func connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 30; i++ {
		if err = pool.Ping(ctx); err == nil {
			return pool, nil
		}
		log.Printf("waiting for database: %v", err)
		time.Sleep(time.Second)
	}
	return nil, err
}

func databaseURL() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(getenv("POSTGRES_USER", "hunter"), os.Getenv("POSTGRES_PASSWORD")),
		Host:     getenv("DB_HOST", "localhost") + ":" + getenv("DB_PORT", "5432"),
		Path:     getenv("POSTGRES_DB", "bounty"),
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

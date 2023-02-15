package handlers

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5"
	"net/http"
	"time"
)

func (s *Service) PingHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", s.cfg.DatabaseDsn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

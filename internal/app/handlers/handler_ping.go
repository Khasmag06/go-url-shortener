package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Service) PingHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", s.cfg.DatabaseDsn)
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

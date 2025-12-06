package apiserver

import (
	sqlstore "RestApi/internal/app/store/sqlStore"
	"database/sql"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	store := sqlstore.NewStore(db)
	srv := NewServer(store)
	return http.ListenAndServe(config.BinAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		return nil, err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

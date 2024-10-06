package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/EM_test/app/internal/models"
	"github.com/Smile8MrBread/EM_test/app/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewConn(url string) (*sql.DB, error) {
	const op = "storage.postgres.New"
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return conn, nil
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Add(ctx context.Context, squad, song, text string) (int64, error) {
	const op = "storage.postgres.Add"

	var id int64
	err := s.db.QueryRowContext(ctx,
		fmt.Sprintf("INSERT INTO Songs(song, squad, text) VALUES ('%s', '%s', '%s') RETURNING id",
			song, squad, text)).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Update(ctx context.Context, id int64, squad, song, text string) error {
	const op = "storage.postgres.Update"

	stmt, err := s.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE songs SET squad = '%s', song = '%s', text = '%s' WHERE id = '%d'",
			squad, song, text, id))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c, err := stmt.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if c == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrSongNotFound)
	}

	return nil
}

func (s *Storage) Text(ctx context.Context, id int64) (string, error) {
	const op = "storage.postgres.Text"

	var text string
	err := s.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT text FROM Songs WHERE id = '%d'", id)).Scan(&text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrSongNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return text, nil
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	const op = "storage.postgres.Delete"

	stmt, err := s.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM songs WHERE id = '%d'", id))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c, err := stmt.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if c == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrSongNotFound)
	}

	return nil
}

func (s *Storage) Library(ctx context.Context, order, field string) ([]models.Song, error) {
	const op = "storage.postgres.Library"

	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM songs ORDER BY %s %s", field, order))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var data []models.Song
	song := models.Song{}
	for rows.Next() {
		if err = rows.Scan(&song.Id, &song.Song, &song.Squad, &song.Text); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		data = append(data, song)
	}

	return data, nil
}

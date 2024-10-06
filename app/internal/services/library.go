package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/EM_test/app/internal/models"
	"github.com/Smile8MrBread/EM_test/app/internal/storage"
	"log/slog"
	"reflect"
	"strings"
)

type Add interface {
	Add(ctx context.Context, squad, song, text string) (int64, error)
	Update(ctx context.Context, id int64, squad, song, text string) error
	Text(ctx context.Context, id int64) (string, error)
	Delete(ctx context.Context, id int64) error
	Library(ctx context.Context, order, field string) ([]models.Song, error)
}

type Lib struct {
	log   *slog.Logger
	adder Add
}

func New(log *slog.Logger, adder Add) *Lib {
	return &Lib{
		log:   log,
		adder: adder,
	}
}

func (l *Lib) Add(ctx context.Context, squad, song, text string) (int64, error) {
	const op = "services.library.Add"
	log := l.log.With(slog.String("op", op))
	log.Info("Adding a new song")

	if err := validateAdd(squad, song, text); err != nil {
		log.Error("Failed to add song", slog.String("error", err.Error()))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := l.adder.Add(ctx, squad, song, text)
	if err != nil {
		log.Error("Failed to add song", slog.String("error", err.Error()))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (l *Lib) Update(ctx context.Context, id int64, squad, song, text string) error {
	const op = "services.library.Update"
	log := l.log.With(slog.String("op", op), slog.Int("id", int(id)))
	log.Info("Updating song")

	if err := validateUpdate(id, squad, song, text); err != nil {
		log.Error("Failed to update the song", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	err := l.adder.Update(ctx, id, squad, song, text)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			log.Error("Song not found", slog.String("error", err.Error()))
			return fmt.Errorf("%s:%w", op, storage.ErrSongNotFound)
		}

		log.Error("Failed to update the song", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (l *Lib) Text(ctx context.Context, id int64) ([]string, error) {
	const op = "services.library.Text"
	log := l.log.With(slog.String("op", op), slog.Int("id", int(id)))
	log.Info("Getting text of song")

	if err := validateText(id); err != nil {
		log.Error("Failed to get text", slog.String("error", err.Error()))
		return nil, ErrInvalidId
	}

	text, err := l.adder.Text(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			log.Error("Song not found", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("Failed to get text")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return strings.Split(text, "\n"), nil
}

func (l *Lib) Delete(ctx context.Context, id int64) error {
	const op = "services.library.Delete"
	log := l.log.With(slog.String("op", op), slog.Int("id", int(id)))
	log.Info("Deleting song")

	if err := validateDelete(id); err != nil {
		log.Error("Failed to delete song", slog.String("error", err.Error()))
		return ErrInvalidId
	}

	err := l.adder.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			log.Error("Song not found", slog.String("error", err.Error()))
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Error("Failed to delete song", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (l *Lib) Library(ctx context.Context, order, field string, chunk int64) ([][]models.Song, error) {
	const op = "services.library.Library"
	log := l.log.With(slog.String("op", op))
	log.Info("Getting all library")

	if err := validateLibrary(order, field, chunk); err != nil {
		log.Error("Failed to get all library", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data, err := l.adder.Library(ctx, order, field)
	if err != nil {
		log.Error("Failed to get all library", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var res [][]models.Song
	var c int64
	for i := chunk; i < int64(len(data)); i += chunk {
		res = append(res, data[i-chunk:i])
		c += chunk
	}
	if c < int64(len(data)) {
		res = append(res, data[c:])
	}

	return res, nil
}

func validateAdd(squad, song, text string) error {
	if reflect.TypeOf(squad).Kind() != reflect.String || len(squad) > 256 || squad == "" {
		return ErrInvalidSquad
	}
	if reflect.TypeOf(song).Kind() != reflect.String || len(song) > 256 || song == "" {
		return ErrInvalidSong
	}
	if reflect.TypeOf(text).Kind() != reflect.String && text != "" {
		return ErrInvalidText
	}

	return nil
}

func validateUpdate(id int64, squad, song, text string) error {
	if reflect.TypeOf(squad).Kind() != reflect.String || len(squad) > 256 || squad == "" {
		return ErrInvalidSquad
	}
	if reflect.TypeOf(song).Kind() != reflect.String || len(song) > 256 || song == "" {
		return ErrInvalidSong
	}
	if reflect.TypeOf(text).Kind() != reflect.String && text != "" {
		return ErrInvalidText
	}
	if reflect.TypeOf(id).Kind() != reflect.Int64 {
		return ErrInvalidId
	}

	return nil
}

func validateText(id int64) error {
	if reflect.TypeOf(id).Kind() != reflect.Int64 {
		return ErrInvalidId
	}

	return nil
}

func validateDelete(id int64) error {
	if reflect.TypeOf(id).Kind() != reflect.Int64 {
		return ErrInvalidId
	}

	return nil
}

func validateLibrary(order, field string, chunk int64) error {
	if reflect.TypeOf(chunk).Kind() != reflect.Int64 {
		return ErrInvalidChunk
	}
	if reflect.TypeOf(order).Kind() != reflect.String || (order != "ASC" && order != "DESC") {
		return ErrInvalidOrder
	}
	if reflect.TypeOf(field).Kind() != reflect.String ||
		(field != "song" && field != "text" && field != "squad" && field != "id") {
		return ErrInvalidField
	}

	return nil
}

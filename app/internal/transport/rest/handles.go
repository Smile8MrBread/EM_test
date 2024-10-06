package rest

import (
	"encoding/json"
	"errors"
	"github.com/Smile8MrBread/EM_test/app/internal/models"
	"github.com/Smile8MrBread/EM_test/app/internal/services"
	"github.com/Smile8MrBread/EM_test/app/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func StartServer(r *chi.Mux, lib *services.Lib) {
	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		s := models.Song{}

		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			w.WriteHeader(500)
			b, _ := json.Marshal(models.ErrorResp{Error: "Internal server error"})
			w.Write(b)
			return
		}

		id, err := lib.Add(r.Context(), s.Squad, s.Song, s.Text)
		if err != nil {
			var b []byte
			w.WriteHeader(400)

			switch {
			case errors.Is(err, services.ErrInvalidId):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid id"})
			case errors.Is(err, services.ErrInvalidSong):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid song"})
			case errors.Is(err, services.ErrInvalidSquad):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid squad"})
			case errors.Is(err, services.ErrInvalidText):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid text"})
			default:
				w.WriteHeader(500)
				b, _ = json.Marshal(models.ErrorResp{Error: "Internal server error"})
			}

			w.Write(b)
			return
		}

		resp, _ := json.Marshal(map[string]int64{"Id": id})
		w.WriteHeader(200)
		w.Write(resp)
	})

	r.Patch("/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		s := models.Song{}
		paramId := chi.URLParam(r, "id")

		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			w.WriteHeader(500)
			b, _ := json.Marshal(models.ErrorResp{Error: "Internal server error"})
			w.Write(b)
			return
		}

		id, err := strconv.Atoi(paramId)
		if err != nil {
			w.WriteHeader(400)
			b, _ := json.Marshal(models.ErrorResp{Error: "Invalid id"})
			w.Write(b)
			return
		}

		err = lib.Update(r.Context(), int64(id), s.Squad, s.Song, s.Text)
		if err != nil {
			var b []byte
			w.WriteHeader(400)

			switch {
			case errors.Is(err, services.ErrInvalidId):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid id"})
			case errors.Is(err, services.ErrInvalidSong):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid song"})
			case errors.Is(err, services.ErrInvalidSquad):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid squad"})
			case errors.Is(err, services.ErrInvalidText):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid text"})
			case errors.Is(err, storage.ErrSongNotFound):
				b, _ = json.Marshal(models.ErrorResp{Error: "Song not found"})
			default:
				w.WriteHeader(500)
				b, _ = json.Marshal(models.ErrorResp{Error: "Internal server error"})
			}

			w.Write(b)
			return
		}

		w.WriteHeader(200)
	})

	r.Get("/text/{id}", func(w http.ResponseWriter, r *http.Request) {
		paramId := chi.URLParam(r, "id")
		id, err := strconv.Atoi(paramId)
		if err != nil {
			w.WriteHeader(400)
			b, _ := json.Marshal(models.ErrorResp{Error: "Invalid id"})
			w.Write(b)
			return
		}

		text, err := lib.Text(r.Context(), int64(id))
		if err != nil {
			var b []byte
			w.WriteHeader(400)

			switch {
			case errors.Is(err, services.ErrInvalidId):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid id"})
			case errors.Is(err, storage.ErrSongNotFound):
				b, _ = json.Marshal(models.ErrorResp{Error: "Song not found"})
			default:
				w.WriteHeader(500)
				b, _ = json.Marshal(models.ErrorResp{Error: "Internal server error"})
			}

			w.Write(b)
			return
		}

		w.WriteHeader(200)
		b, _ := json.Marshal(text)
		w.Write(b)
	})

	r.Delete("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		paramId := chi.URLParam(r, "id")
		id, err := strconv.Atoi(paramId)
		if err != nil {
			w.WriteHeader(400)
			b, _ := json.Marshal(models.ErrorResp{Error: "Invalid id"})
			w.Write(b)
			return
		}

		err = lib.Delete(r.Context(), int64(id))
		if err != nil {
			var b []byte
			w.WriteHeader(400)

			switch {
			case errors.Is(err, services.ErrInvalidId):
				b, _ = json.Marshal(models.ErrorResp{Error: "Invalid id"})
			case errors.Is(err, storage.ErrSongNotFound):
				b, _ = json.Marshal(models.ErrorResp{Error: "Song not found"})
			default:
				w.WriteHeader(500)
				b, _ = json.Marshal(models.ErrorResp{Error: "Internal server error"})
			}

			w.Write(b)
			return
		}

		w.WriteHeader(200)
	})

	r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()

		chunk, err := strconv.Atoi(url.Get("pagination"))
		if err != nil {
			w.WriteHeader(400)
			b, _ := json.Marshal(models.ErrorResp{Error: "invalid n"})
			w.Write(b)
			return
		}

		data, err := lib.Library(r.Context(), url.Get("order"), url.Get("field"), int64(chunk))
		if err != nil {
			w.WriteHeader(500)
			b, _ := json.Marshal(models.ErrorResp{Error: "Internal server error"})
			w.Write(b)
			return
		}

		w.WriteHeader(200)
		b, _ := json.Marshal(data)
		w.Write(b)
	})

	http.ListenAndServe(":8080", r)
}

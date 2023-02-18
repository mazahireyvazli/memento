package memento

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type MementoServerConfig struct {
	*MementoConfig
	Host string
	Port int
}

func RunServer(c *MementoServerConfig) error {
	var cache, err = NewMemento[string](c.MementoConfig)
	defer cache.Close()

	if err != nil {
		return err
	}

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(apiV1 chi.Router) {
		apiV1.Get("/cache/{key}", func(w http.ResponseWriter, r *http.Request) {
			key := chi.URLParam(r, "key")
			status := http.StatusOK

			v, ok := cache.Get(key)
			if !ok {
				status = http.StatusNotFound
			}

			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(status)
			w.Write(v)
		})

		apiV1.Put("/cache/{key}/{value}", func(w http.ResponseWriter, r *http.Request) {
			key := chi.URLParam(r, "key")
			value := chi.URLParam(r, "value")

			go cache.Set(key, []byte(value))

			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNoContent)
			w.Write(nil)
		})

		apiV1.Delete("/cache/:key", func(w http.ResponseWriter, r *http.Request) {
			key := chi.URLParam(r, "key")

			go cache.Delete(key)

			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNoContent)
			w.Write(nil)
		})
	})

	log.Println("Starting server at", fmt.Sprintf("%s:%d", c.Host, c.Port))
	return http.ListenAndServe(fmt.Sprintf("%s:%d", c.Host, c.Port), r)
}

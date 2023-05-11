package main

import (
	"flag"
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"strings"
)

var routes = flag.Bool("routes", false, "Generate router documentation")
var rsc = redisearch.NewClient("localhost:6379", "data")

// /////
func setRoutes(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		replyMessage(w, "OK")
	})

	r.Get("/test/500", func(w http.ResponseWriter, r *http.Request) {
		reply500(w, "This is a test.")
	})

	r.Get("/test/ok", func(w http.ResponseWriter, r *http.Request) {
		replyOK(w)
	})

	r.Get("/search/{query}", func(w http.ResponseWriter, r *http.Request) {
		q := chi.URLParam(r, "query")
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		println("offset = ", offset, " : limit = ", limit)
		result, err := rsSearch(q, limit, offset)
		if err == nil {
			replyJSON(w, r, result)
			return
		}
		println(err.Error())
		reply500(w, "Check your inputs")
	})

	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

}

// /////
func main() {
	flag.Parse()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	setRoutes(r)
	http.ListenAndServe(":3000", r)
}

// /////
func replyMessage(w http.ResponseWriter, message string) {
	j := fmt.Sprintf("{\"message\":\"%s\"}", strings.Replace(message, "\"", "'", 0))
	w.Write([]byte(j))
}

func reply500(w http.ResponseWriter, message string) {
	w.WriteHeader(500)
	replyMessage(w, message)
}

func replyOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	replyMessage(w, "OK")
}

func replyJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.JSON(w, r, v)
}

// /////
func rsSearch(s string, offset int, limit int) (docs []redisearch.Document, err error) {
	if limit > 50 {
		limit = 50
	} else if limit < 10 {
		limit = 10
	}
	docs, _, err = rsc.Search(redisearch.NewQuery(s).Limit(offset, limit))
	if err != nil {
		return nil, err
	}
	return docs, nil
}

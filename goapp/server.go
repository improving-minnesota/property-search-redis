package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
	"regexp"
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

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		filter := r.URL.Query().Get("f")
		println(filter)
		offset, _ := strconv.Atoi(r.URL.Query().Get("o"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("l"))
		result, err := rsInitSearch(query, offset, limit)
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
func rsInitSearch(query string, offset int, limit int) (docs []redisearch.Document, err error) {
	// check offset and limit
	if limit > 100 {
		limit = 100
	}
	if limit < 25 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	// check for usable strings as search "tokens" - must have at least 1
	tokens := rsQuery2Tokens(query)
	if len(tokens) == 0 {
		return docs, errors.New("invalid search - must be one or more words or numbers")
	}
	return rsIterativeSearch(tokens, offset, limit)
}

func rsIterativeSearch(tokens []string, offset int, limit int) (docs []redisearch.Document, err error) {
	phraseExact := strings.Join(tokens, " ")                                                 // multi-word exact phrase
	phraseWild := strings.Join(tokens, "*")                                                  // multi-word wildcard phrase search
	unionExact := strings.Join(tokens, "|")                                                  // union search (OR) without wildcard
	unionWild := "*" + strings.Join(tokens, "*|*") + "*"                                     // union search (OR) with wildcards
	query := "(" + phraseExact + ")|(" + phraseWild + "*)|(" + unionExact + ")|" + unionWild // put it all together
	docs, _, err = rsc.Search(rsNewQuery(query, offset, limit))
	return docs, err
}

func rsNewQuery(query string, offset int, limit int) *redisearch.Query {
	println("query=["+query+"] offset:", offset, " limit:", limit)
	search := redisearch.NewQuery(query)
	if offset > 0 {
		search.Paging.Offset = offset
	}
	if limit > 0 {
		search.Paging.Num = limit
	}
	return search
}

func rsQuery2Tokens(query string) (tokens []string) {
	re := regexp.MustCompile("[^a-zA-Z\\d:]")
	tokens = re.Split(query, -1)
	return
}

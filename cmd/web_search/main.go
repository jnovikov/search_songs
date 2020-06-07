package main

import (
	"context"
	"flag"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	_ "net/http/pprof"
	"search_songs/pkg/searcher"
)

var (
	listen   = flag.String("l", ":4141", "host:port to listen")
	jobCount = flag.Int("j", 4, "job count")
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("Usage: web_search <directory>")
	}

	dir := flag.Arg(0)
	ds, err := searcher.NewDirSearcher(dir, *jobCount)
	if err != nil {
		log.Fatalf("got error: %v", err)
	}

	log.Printf("Going to use %d go-routines\n", ds.JobCount)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		query := c.QueryParam("q")
		if query == "" {
			return c.String(http.StatusUnprocessableEntity, "no query")
		}
		resp := ds.Search(context.TODO(), query)
		return c.JSON(http.StatusOK, resp)
	})
	// Сервер для профилирования!
	go http.ListenAndServe("127.0.0.1:8080", nil)
	e.Logger.Fatal(e.Start(*listen))
}

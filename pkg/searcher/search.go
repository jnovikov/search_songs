package searcher

import (
	"context"
)

type Response struct {
	Line     string
	LineNum  int
	SongName string
}

type Searcher interface {
	Search(ctx context.Context, query string) []Response
}

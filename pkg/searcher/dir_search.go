package searcher

import (
	"bufio"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

func NewDirSearcher(dir string, jc int) (*DirSearcher, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, errors.New("not a directory")
	}

	if jc <= 0 {
		jc = runtime.NumCPU()
	}

	return &DirSearcher{
		Dir:      dir,
		JobCount: jc,
	}, nil
}

type DirSearcher struct {
	Dir      string
	JobCount int
}

func (ds *DirSearcher) scan(r io.Reader, query string) []Response {
	// Используя bufio читайте из r построчно.
	// Для каждой строки проверьте входит ли в нее подстрока query используя strings.Contains
	// Верните массив вхождений. Желательно сделать регистронезависимый поиск(lower)
	scan := bufio.NewScanner(r)
	res := make([]Response, 0)
	ln := 1
	for scan.Scan() {
		line := scan.Text()
		if strings.Contains(line, query) {
			res = append(res, Response{
				Line:     line,
				LineNum:  ln,
				SongName: "",
			})
		}
		ln++
	}
	return res
}

func (ds *DirSearcher) scanRegexp(r io.Reader, queryRe *regexp.Regexp) []Response {
	scan := bufio.NewScanner(r)
	res := make([]Response, 0)
	ln := 1
	for scan.Scan() {
		line := scan.Text()
		if queryRe.Match([]byte(line)) {
			res = append(res, Response{
				Line:     line,
				LineNum:  ln,
				SongName: "",
			})
		}
		ln++
	}
	return res
}

func (ds *DirSearcher) scanFile(filename string, query string) []Response {
	f, err := os.Open(path.Join(ds.Dir, filename))
	if err != nil {
		log.Printf("Open error: %v\n", err)
		return nil
	}
	defer f.Close()
	responses := ds.scan(f, query)
	for i := range responses {
		responses[i].SongName = filename
	}
	return responses
}

func (ds *DirSearcher) Search(ctx context.Context, query string) []Response {
	files := make(chan string)
	infos, err := ioutil.ReadDir(ds.Dir)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil
	}
	go func() {
		for _, f := range infos {
			if f.IsDir() {
				continue
			}
			files <- f.Name()
		}
		close(files)
	}()

	res := make(chan []Response)
	var wg sync.WaitGroup
	wg.Add(ds.JobCount)
	for i := 0; i < ds.JobCount; i++ {
		go func() {
			for f := range files {
				res <- ds.scanFile(f, query)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	responses := make([]Response, 0)
	for r := range res {
		responses = append(responses, r...)
	}

	return responses
}

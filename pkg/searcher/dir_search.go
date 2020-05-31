package searcher

import (
	"context"
    "io"
    "os"
)

type DirSearcher struct {
	Dir string
    JobCount int
}

func (ds *DirSearcher) scan(r io.Reader, query string) []Response {
    // Используя bufio читайте из r построчно.
    // Для каждой строки проверьте входит ли в нее подстрока query используя strings.Contains
    // Верните массив вхождений. Желательно сделать регистронезависимый поиск(lower)
    return nil
}

func (ds *DirSearcher) scanFile(filename string, query string) []Response {
    f, err := os.Open(filename)
    if err != nil {
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
    // TODO: написать полноценный поиск по всем файлам в директории, используя горутины и функцию scanFile.
    // крайне хотелось бы еще и учитывать контекст и иметь возможность завершить досрочно функцию.
    // Количество го-рутин пользователь должен задать с помощью JobCount
	return nil
}

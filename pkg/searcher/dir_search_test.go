package searcher

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

func TestDirsearch_Search(t *testing.T) {
	ds := DirSearcher{Dir: "testdata", JobCount: 1}
	query := "hello"
	responses := ds.Search(context.TODO(), query)
	if len(responses) != 3 {
		t.Errorf("Expected %d responses, got %d", 3, len(responses))
	}
}

func TestDirsearch_scan(t *testing.T) {
	ds := DirSearcher{Dir: "testdata"}
	query := "hello"
	data := `hello world
    world hello
    hello
    world
    `
	responses := ds.scan(strings.NewReader(data), query)
	if len(responses) != 3 {
		// Закончить тест с ошибкой
		t.Fatalf("Expected %d responses, got %d", 3, len(responses))
	}
	if responses[0].LineNum != 1 || responses[0].Line != "hello world" {
		// Error - сообщить об ошибке но продолжить выполнять тест
		t.Errorf("Expected %d line be %s", 0, "hello world")
	}

}

func TestDirsearch_scanFile(t *testing.T) {
	// Your test will be here.
	ds := DirSearcher{Dir: "testdata"}
	res := ds.scanFile("file1.txt", "world")
	if len(res) != 2 {
		t.Errorf("Expected 2 responses got, %d", len(res))
	}
	if res[0].SongName != res[1].SongName || !strings.Contains(res[0].SongName, "file1") {
		t.Errorf("Expected songname be = file1.txt, got %s", res[0].SongName)
	}
}

func BenchmarkDirSearcher_Search(b *testing.B) {
	ds := DirSearcher{Dir: "testdata", JobCount: 20}
	for i := 0; i < b.N; i++ {
		ds.Search(context.TODO(), "а")
	}
}

func BenchmarkDirSearcher_Scan(b *testing.B) {
	ds := DirSearcher{Dir: "testdata", JobCount: 1}
	query := "hello"
	data := `hello world
    world hello
    hello
    world
    `
	for i := 0; i < b.N; i++ {
		responses := ds.scan(strings.NewReader(data), query)
		if len(responses) < 3 {
			b.Fatalf("wrong")
		}
	}
}

func BenchmarkDirSearcher_ScanRegexp(b *testing.B) {
	ds := DirSearcher{Dir: "testdata", JobCount: 1}
	query := "hello"
	data := `hello world
    world hello
    hello
    world
    `
	reg, _ := regexp.Compile(query)
	for i := 0; i < b.N; i++ {
		responses := ds.scanRegexp(strings.NewReader(data), reg)
		if len(responses) < 3 {
			b.Fatalf("wrong")
		}
	}
}

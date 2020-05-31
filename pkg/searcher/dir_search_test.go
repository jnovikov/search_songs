package searcher

import (
    "testing"
    "context"
    "strings"
)

func TestDirsearch_Search(t *testing.T) {
    ds := DirSearcher{Dir: "testdata"}
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
    t.Fatal("Not implemented now")
}

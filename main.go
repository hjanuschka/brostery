package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type chromeRecord struct {
	ID                 int
	URL                string
	Title              string
	VisitCount         int
	TypedCount         int
	LastVisitTime      int64
	LastVisitTimeClean time.Time
	SearchTerm         string
}

func main() {
	p := "/Users/XXX/Library/Application Support/Google/Chrome/Default/History"
	if len(os.Args) > 0 {
		p = os.Args[1]
	}
	pCopy := "working_copy"
	copy(p, pCopy)

	db, err := sql.Open("sqlite3", pCopy)
	if err != nil {
		log.Panic(err)
	}

	// REF: https://gist.github.com/dropmeaword/9372cbeb29e8390521c2#chrome
	sqlStmt := `SELECT 
	id,
	url,
	title,
	visit_count,
	typed_count,
  last_visit_time
FROM urls`
	r, err := db.Query(sqlStmt)
	re := regexp.MustCompile(`q=(.*?)&`)
	for r.Next() {
		row := chromeRecord{}
		err = r.Scan(&row.ID, &row.URL, &row.Title, &row.VisitCount, &row.TypedCount, &row.LastVisitTime)
		row.LastVisitTimeClean = time.Unix(row.LastVisitTime/1000000-11644473600, 0)
		if strings.Contains(row.URL, "google.com/search") && strings.Contains(row.URL, "q=") {
			matches := re.FindStringSubmatch(row.URL)

			if len(matches) > 0 {
				v, e := url.QueryUnescape(matches[1])
				if e == nil {
					row.SearchTerm = v
				}

				fmt.Printf("%d;%s;%s\n", row.ID, row.LastVisitTimeClean, row.SearchTerm)
			}

		}
	}

}

/// REF: https://opensource.com/article/18/6/copying-files-go
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

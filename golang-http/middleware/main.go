package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type contextKey string

const dbContextKey contextKey = "db"

// Logger Middleware
type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handler: handlerToWrap}
}

// ResponseHeader middleware
type ResponseHeader struct {
	handler     http.Handler
	headername  string
	headervalue string
}

func (rh *ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(rh.headername, rh.headervalue)
	rh.handler.ServeHTTP(w, r)
}

func NewResponseHeader(handlerToWrap http.Handler, headerName string, headerValue string) *ResponseHeader {
	return &ResponseHeader{handler: handlerToWrap, headername: headerName, headervalue: headerValue}
}

//Database middleware

type DBMiddleware struct {
	handler http.Handler
	db      *sql.DB
}

func (dm *DBMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), dbContextKey, dm.db)
	dm.handler.ServeHTTP(w, r.WithContext(ctx))
}

func NewDBMiddleware(handlerToWrap http.Handler, db *sql.DB) *DBMiddleware {
	return &DBMiddleware{handler: handlerToWrap, db: db}
}

// get the dbcontext. used in handlers to allow db utilisation
func GetDBContext(r *http.Request) *sql.DB {
	if db, ok := r.Context().Value(dbContextKey).(*sql.DB); ok {
		return db
	}
	return nil
}

// Standard handlers
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You have hit hellohandler"))
}

func CurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	ctime := time.Now().Format(time.Kitchen)
	time.Sleep(500 * time.Millisecond)
	w.Write([]byte(fmt.Sprintf("Current Time : %s", ctime)))
}

func DBHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBContext(r)
	if db != nil {
		var now string
		db.QueryRow("SELECT datetime('now')").Scan(&now)
		w.Write([]byte(fmt.Sprintf("some query exec: %s", now)))
		return
	}
	w.Write([]byte("NEHHHH. No DB"))
}

func main() {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE IF NOT EXISTS demo (id INTEGER)")

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)
	mux.HandleFunc("/time", CurrentTimeHandler)
	mux.HandleFunc("/db", DBHandler)

	newMux := NewLogger(NewResponseHeader(NewDBMiddleware(mux, db), "HOTEL", "Trivago"))

	log.Println("Server is listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", newMux))

}

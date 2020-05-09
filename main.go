package main

import (
	"fmt"
	"log"
	"monkiato/guardian/internal/models"
	"monkiato/guardian/pkg/routes"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var logger *log.Logger
var err error

func main() {
	logger = log.New(os.Stdout, "rest-api ", log.LstdFlags|log.Lshortfile)
	router := mux.NewRouter()

	connectDB()
	defer db.Close()

	logger.Print("creating db...")
	createDB()

	logger.Print("adding routes...")
	auth := router.PathPrefix("/auth").Subrouter()
	routes.NewAuth(db, logger).AddRoutes(auth)

	router.Use(loggingMiddleware)
	router.NotFoundHandler = http.HandlerFunc(notFound)

	handler := cors.Default().Handler(router)

	port := 8080
	logger.Printf("initialization ready. server running at http://localhost:%d", port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func connectDB() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")

	db, err = gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPass))

	if err != nil {
		logger.Panic("failed to connect database: " + err.Error())
	}
	logger.Printf("DB: connected to host: %s, db: %s, user: %s", dbHost, dbName, dbUser)
}

func createDB() {
	// db.Exec("CREATE EXTENSION IF NOT EXISTS hstore;")
	db.AutoMigrate(
		&models.User{})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		logger.Printf("%s  	%d 	%s", r.Method, lrw.statusCode, r.RequestURI)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sorry, this page is not available"))
}

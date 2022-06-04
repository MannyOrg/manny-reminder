package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"manny-reminder/pkg/auth"
	"manny-reminder/pkg/events"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	host     = "manny-master.local"
	port     = 30432
	user     = "postgres"
	password = "JyxXWxxGMz"
	dbname   = "manny-reminder"
)

var bindAddress = ":8080"

func main() {
	err, config := getOAuthConfig()

	l := log.New(os.Stdout, "manny-reminder ", log.LstdFlags)
	db := getDb(err)

	ar := auth.NewRepository(l, db)
	as := auth.NewAuth(l, ar, config)
	ah := auth.NewHandler(as)

	er := events.NewRepository(l, db)
	es := events.NewEvents(er, l, config, as)
	eh := events.NewHandler(es)

	sm := mux.NewRouter()

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/users", ah.GetUsers)
	getR.HandleFunc("/users/add", ah.AddUser)
	getR.HandleFunc("/users/save", ah.SaveUser)
	getR.HandleFunc("/users/events", eh.GetUsersEvents)

	// create a new server
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port $1", bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Println("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	cancelFunc()
	err = s.Shutdown(ctx)
	if err != nil {
		panic(err.Error())
	}
}

func getDb(err error) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func getOAuthConfig() (error, *oauth2.Config) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return err, config
}

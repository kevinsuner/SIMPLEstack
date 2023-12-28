package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

/*** data ***/

var db *sql.DB

/*** operations ***/

/*** views ***/

func admin(w http.ResponseWriter, r *http.Request) {
	// TODO: Check cookie value
	_, err := r.Cookie("simple_stack_token")
	if errors.Is(err, http.ErrNoCookie) {
		t, err := template.New("login").ParseFiles(
			filepath.Join("views", "layouts", "admin_header.tmpl"),
			filepath.Join("views", "layouts", "navbar.tmpl"),
			filepath.Join("views", "layouts", "footer.tmpl"),
			filepath.Join("views", "admin", "login.tmpl"),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
			log.Printf("[ERROR] failed to parse templates: %v\n", err)
			return
		}

		var buf bytes.Buffer
		if err = t.Execute(&buf, nil); err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
			log.Printf("[ERROR] failed to execute template: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to get cookie: %v\n", err)
		return
	}

	t, err := template.New("dashboard").ParseFiles(
		filepath.Join("views", "layouts", "admin_header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "admin", "dashboard.tmpl"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to parse templates: %v\n", err)
		return
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, nil); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to execute template: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").ParseFiles(
		filepath.Join("views", "layouts", "header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "index.tmpl"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to parse templates: %v\n", err)
		return
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, nil); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to execute template: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

/*** init ***/

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[ERROR] failed to load .env file: %v\n", err)
	}

	var connStr string = fmt.Sprintf("postgresql://%s:%s@tcp/%s?sslmode=disable",
		os.Getenv("PQ_USER"), os.Getenv("PQ_PASS"), os.Getenv("PQ_NAME"))

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("[ERROR] failed to initialize db: %v\n", err)
	}
	log.Println("[INFO] successfully connected to db")
}

func main() {
	mux := http.NewServeMux()
	
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", index)
	mux.HandleFunc(fmt.Sprintf("/%s", os.Getenv("ADMIN_URL")), admin)

	log.Printf("[INFO] started http server at port %s\n", os.Getenv("PORT"))
	http.ListenAndServe(os.Getenv("PORT"), mux)
}

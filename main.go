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
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

/*** data ***/

var db *sql.DB

/*** endpoints ***/

func getArticles(w http.ResponseWriter, r *http.Request) {
	// mocked db call
	type Article struct {
		// Think more deeply about the data that an article should hold
		ID string
		Title string
		Author string
		Excerpt string
		CreatedAt string
		UpdatedAt string
		Status string
	}

	articles := []Article{
		{ID: "rand-id-1", Title: "An example title", Author: "Kevin Suñer", Excerpt: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.", CreatedAt: time.Now().Format(time.ANSIC), Status: "unpublished"},
		{ID: "rand-id-2", Title: "Another example title", Author: "Kevin Suñer", Excerpt: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.", CreatedAt: time.Now().Format(time.ANSIC), UpdatedAt: time.Now().Format(time.ANSIC), Status: "published"},
	}

	var isAdmin bool
	if len(r.URL.Query().Get("admin")) > 0 {
		val, err := strconv.ParseBool(r.URL.Query().Get("admin"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse bool: %v", err), http.StatusBadRequest)
			log.Printf("[ERROR] failed to parse bool: %v\n", err)
			return
		}

		if val {
			cookie, err := r.Cookie("simple_stack_token")
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, fmt.Sprintf("failed to authenticate"), http.StatusUnauthorized)
				log.Println("[ERROR] failed to authenticate")
				return
			} else if err != nil {
				http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
				log.Printf("[ERROR] failed to get cookie: %v\n", err)
				return
			}

			if cookie.Value == os.Getenv("SIMPLE_STACK_TOKEN") {
				isAdmin = true
			} else {
				http.Error(w, fmt.Sprintf("failed to authenticate"), http.StatusUnauthorized)
				log.Println("[ERROR] failed to authenticate")
				return
			}
		}
	}

	t, err := template.New("articles").ParseFiles(filepath.Join("views", "articles", "articles.tmpl"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse template: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to parse template: %v\n", err)
		return
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, struct{Articles []Article; IsAdmin bool}{Articles: articles, IsAdmin: isAdmin}); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to execute template: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to parse form: %v\n", err)
		return
	}

	if r.Form.Get("username") == os.Getenv("USER_NAME") &&
		r.Form.Get("password") == os.Getenv("USER_PASS") {
		cookie := http.Cookie{}
		cookie.Name = "simple_stack_token"
		cookie.Value = os.Getenv("SIMPLE_STACK_TOKEN")
		cookie.Expires = time.Now().Add(time.Hour * 1)
		cookie.Secure = true
		cookie.HttpOnly = true
		cookie.Path = "/"
		http.SetCookie(w, &cookie)
		w.Header().Add("HX-Redirect", fmt.Sprintf("/%s", os.Getenv("ADMIN_URL")))
		return
	} else {
		http.Error(w, fmt.Sprintf("failed to authenticate"), http.StatusUnauthorized)
		log.Println("[ERROR] failed to authenticate")
		return
	}
}

/*** views ***/

func admin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("simple_stack_token")
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

	if cookie.Value == os.Getenv("SIMPLE_STACK_TOKEN") {
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
		return
	} else {
		http.Error(w, fmt.Sprintf("failed to authenticate"), http.StatusUnauthorized)
		log.Println("[ERROR] failed to authenticate")
		return
	}
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
	
	/*** endpoints ***/
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/get/articles", getArticles)
	
	/*** views ***/
	mux.HandleFunc("/", index)
	mux.HandleFunc(fmt.Sprintf("/%s", os.Getenv("ADMIN_URL")), admin)

	log.Printf("[INFO] started http server at port %s\n", os.Getenv("PORT"))
	http.ListenAndServe(os.Getenv("PORT"), mux)
}

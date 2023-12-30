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

/*** models ***/

type Article struct {
	ID uint
	CreatedAt sql.NullString
	UpdatedAt sql.NullString
	Title string
	Excerpt string
	Author string
	Status string
	Content string
}

/*** endpoints ***/

func putArticle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("simple_stack_token")
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, fmt.Sprintf("failed to authenticate: %v", err), http.StatusUnauthorized)
		log.Printf("[ERROR] failed to authenticate: %v\n", err)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to get cookie: %v\n", err)
		return
	}

	if cookie.Value == os.Getenv("SIMPLE_STACK_TOKEN") {
		if len(r.URL.Query().Get("id")) > 0 {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))		
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to parse id: %v", err), http.StatusBadRequest)
				log.Printf("[ERROR] failed to parse id: %v\n", err)
				return
			}

			if err = r.ParseForm(); err != nil {
				http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusInternalServerError)
				log.Printf("[ERROR] failed to parse form: %v\n", err)
				return
			}

			checkEmpty := func(str ...string) error {
				for _, s := range str {
					if s == "" {
						return errors.New("empty string")
					}
				}
				return nil
			}

			title := r.Form.Get("title")
			excerpt := r.Form.Get("excerpt")
			author := r.Form.Get("author")
			status := r.Form.Get("status")
			content := r.Form.Get("content")
			if err = checkEmpty(title, excerpt, author, status, content); err != nil {
				http.Error(w, fmt.Sprintf("failed to validate form values: %v", err), http.StatusBadRequest)
				log.Printf("[ERROR] failed to validate form values: %v\n", err)
				return
			}

			_, err = db.Exec(
				`UPDATE articles SET updated_at=$1, title=$2, excerpt=$3, author=$4, status=$5, content=$6 WHERE id = $7`,
				time.Now().Format(time.ANSIC), title, excerpt, author, status, content, id)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to insert article: %v", err), http.StatusInternalServerError)
				log.Printf("[ERROR] failed to insert article: %v\n", err)
				return
			}

			w.Header().Add("HX-Redirect", fmt.Sprintf("/%s", os.Getenv("ADMIN_URL")))
			return
		} else {
			http.Error(w, "failed due to an id of zero length", http.StatusBadRequest)
			log.Println("[ERROR] failed due to an id of zero length")
			return
		}
	} else {
		http.Error(w, "failed to authenticate", http.StatusUnauthorized)
		log.Println("[ERROR] failed to authenticate", err)
		return
	}
}

func postArticle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("simple_stack_token")
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, fmt.Sprintf("failed to authenticate: %v", err), http.StatusUnauthorized)
		log.Printf("[ERROR] failed to authenticate: %v\n", err)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to get cookie: %v\n", err)
		return
	}

	if cookie.Value == os.Getenv("SIMPLE_STACK_TOKEN") {
		if err = r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusInternalServerError)
			log.Printf("[ERROR] failed to parse form: %v\n", err)
			return
		}

		checkEmpty := func(str ...string) error {
			for _, s := range str {
				if s == "" {
					return errors.New("empty string")
				}
			}
			return nil
		}

		title := r.Form.Get("title")
		excerpt := r.Form.Get("excerpt")
		author := r.Form.Get("author")
		status := r.Form.Get("status")
		content := r.Form.Get("content")
		if err = checkEmpty(title, excerpt, author, status, content); err != nil {
			http.Error(w, fmt.Sprintf("failed to validate form values: %v", err), http.StatusBadRequest)
			log.Printf("[ERROR] failed to validate form values: %v\n", err)
			return
		}

		_, err = db.Exec(
			`INSERT INTO articles (created_at, title, excerpt, author, status, content) VALUES ($1, $2, $3, $4, $5, $6)`,
			time.Now().Format(time.ANSIC), title, excerpt, author, status, content)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to insert article: %v", err), http.StatusInternalServerError)
			log.Printf("[ERROR] failed to insert article: %v\n", err)
			return
		}

		w.Header().Add("HX-Redirect", fmt.Sprintf("/%s", os.Getenv("ADMIN_URL")))
		return
	} else {
		http.Error(w, "failed to authenticate", http.StatusUnauthorized)
		log.Println("[ERROR] failed to authenticate", err)
		return
	}
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		`SELECT id, created_at, updated_at, title, excerpt, author, status FROM articles ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to select articles: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to select articles: %v\n", err)
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		if err = rows.Scan(
			&article.ID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Title,
			&article.Excerpt,
			&article.Author,
			&article.Status); err != nil {
			http.Error(w, fmt.Sprintf("failed to scan value: %v", err), http.StatusInternalServerError)
			log.Printf("[ERROR] failed to scan rows: %v\n", err)
			return
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("failed while iterating: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed while iterating: %v\n", err)
		return
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
				http.Error(w, fmt.Sprintf("failed to authenticate: %v", err), http.StatusUnauthorized)
				log.Printf("[ERROR] failed to authenticate: %v\n", err)
				return
			} else if err != nil {
				http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
				log.Printf("[ERROR] failed to get cookie: %v\n", err)
				return
			}

			if cookie.Value == os.Getenv("SIMPLE_STACK_TOKEN") {
				isAdmin = true
			} else {
				http.Error(w, "failed to authenticate", http.StatusUnauthorized)
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
	data := struct{Articles []Article; IsAdmin bool}{Articles: articles, IsAdmin: isAdmin}
	if err = t.Execute(&buf, data); err != nil {
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

func updateArticle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("simple_stack_token")
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, fmt.Sprintf("failed to authenticate: %v", err), http.StatusUnauthorized)
		log.Printf("[ERROR] failed to authenticate: %v\n", err)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to get cookie: %v\n", err)
		return
	}

	if cookie.Value == os.Getenv("SIMPLE_STACK_TOKEN") {
		if len(r.URL.Query().Get("id")) > 0 {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))		
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to parse id: %v", err), http.StatusBadRequest)
				log.Printf("[ERROR] failed to parse id: %v\n", err)
				return
			}

			var article Article
			err = db.QueryRow(
				`SELECT id, title, excerpt, author, status, content FROM articles WHERE id = $1`, id).Scan(
				&article.ID,
				&article.Title,
				&article.Excerpt,
				&article.Author,
				&article.Status,
				&article.Content,
			)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to query article: %v", err), http.StatusInternalServerError)
				log.Printf("[ERROR] failed to query article: %v\n", err)
				return
			}

			t, err := template.New("update_article").ParseFiles(
				filepath.Join("views", "layouts", "admin_header.tmpl"),
				filepath.Join("views", "layouts", "navbar.tmpl"),
				filepath.Join("views", "layouts", "footer.tmpl"),
				filepath.Join("views", "admin", "update_article.tmpl"),
			)

			if err != nil {
				http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
				log.Printf("[ERROR] failed to parse templates: %v\n", err)
				return
			}

			var buf bytes.Buffer
				if err = t.Execute(&buf, article); err != nil {
				http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
				log.Printf("[ERROR] failed to execute template: %v\n", err)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(buf.Bytes())
			return
		} else {
			http.Error(w, "failed due to an id of zero length", http.StatusBadRequest)
			log.Println("[ERROR] failed due to an id of zero length")
			return
		}
	} else {
		http.Error(w, "failed to authenticate", http.StatusUnauthorized)
		log.Println("[ERROR] failed to authenticate")
		return
	}
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("simple_stack_token")
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, fmt.Sprintf("failed to authenticate: %v", err), http.StatusUnauthorized)
		log.Printf("[ERROR] failed to authenticate: %v\n", err)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to get cookie: %v\n", err)
		return
	}

	if cookie.Value == os.Getenv("SIMPLE_STACK_TOKEN") {
		t, err := template.New("create_article").ParseFiles(
			filepath.Join("views", "layouts", "admin_header.tmpl"),
			filepath.Join("views", "layouts", "navbar.tmpl"),
			filepath.Join("views", "layouts", "footer.tmpl"),
			filepath.Join("views", "admin", "create_article.tmpl"),
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
		http.Error(w, "failed to authenticate", http.StatusUnauthorized)
		log.Println("[ERROR] failed to authenticate")
		return
	}
}

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
		http.Error(w, "failed to authenticate", http.StatusUnauthorized)
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

	var connStr string = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("PQ_USER"), os.Getenv("PQ_PASS"), os.Getenv("PQ_IP"), os.Getenv("PQ_NAME"))

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
	mux.HandleFunc("/post/article", postArticle)
	mux.HandleFunc("/put/article", putArticle)
	
	/*** views ***/
	mux.HandleFunc("/", index)
	mux.HandleFunc(fmt.Sprintf("/%s", os.Getenv("ADMIN_URL")), admin)
	mux.HandleFunc("/create/article", createArticle)
	mux.HandleFunc("/update/article", updateArticle)

	log.Printf("[INFO] started http server at port %s\n", os.Getenv("PORT"))
	http.ListenAndServe(os.Getenv("PORT"), mux)
}

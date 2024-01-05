package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	emptyString			error = errors.New("empty string")
	invalidCredentials	error = errors.New("invalid username or password")
)

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse id: %v", err), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM articles WHERE id = $1", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete article: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", "/dashboard")
}

func PutArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse id: %v", err), http.StatusBadRequest)
		return
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	isEmpty := func(str ...string) error {
		for _, s := range str {
			if len(s) == 0 {
				return emptyString
			}
		}
		return nil
	}

	var (
		title	string = r.Form.Get("title")
		slug	string = r.Form.Get("slug")
		excerpt	string = r.Form.Get("excerpt")
		author	string = r.Form.Get("author")
		status	string = r.Form.Get("status")
		content	string = r.Form.Get("content")
	)

	if err = isEmpty(title, slug, excerpt, author, status, content); err != nil {
		http.Error(w, fmt.Sprintf("failed to validate form values: %v", err), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		`UPDATE articles SET updated_at=$1, title=$2, slug=$3, excerpt=$4, author=$5, status=$6, content=$7 WHERE id = $8`,
		time.Now().Format(time.ANSIC), title, slug, excerpt, author, status, content, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update article: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
	<div class="alert alert-success" role="alert">
		<p>¡Hey! The article has been successfully edited</p>
		<hr>
		<a href="/dashboard" class="link-success mb-0">Back to Dashboard &#x2192;</a>
	</div>`))
}

func PostArticle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	isEmpty := func(str ...string) error {
		for _, s := range str {
			if len(s) == 0 {
				return emptyString
			}
		}
		return nil
	}

	var (
		title	string = r.Form.Get("title")
		slug	string = r.Form.Get("slug")
		excerpt	string = r.Form.Get("excerpt")
		author	string = r.Form.Get("author")
		status	string = r.Form.Get("status")
		content	string = r.Form.Get("content")
	)

	if err := isEmpty(title, slug, excerpt, author, status, content); err != nil {
		http.Error(w, fmt.Sprintf("failed to validate form values: %v", err), http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		`INSERT INTO articles (created_at, title, slug, excerpt, author, status, content) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		time.Now().Format(time.ANSIC), title, slug, excerpt, author, status, content)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to post article: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
	<div class="alert alert-success" role="alert">
		<p>¡Hooray! A new article has been created</p>
		<hr>
		<a href="/dashboard" class="link-success mb-0">Back to Dashboard &#x2192;</a>
	</div>`))
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		`SELECT id, created_at, updated_at, title, slug, excerpt, author, status FROM articles ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get articles: %v", err), http.StatusInternalServerError)
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
			&article.Slug,
			&article.Excerpt,
			&article.Author,
			&article.Status); err != nil {
			http.Error(w, fmt.Sprintf("failed to scan value: %v", err), http.StatusInternalServerError)
			return
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("failed while iterating: %v", err), http.StatusInternalServerError)
		return
	}

	var isAdmin bool
	if len(r.URL.Query().Get("admin")) > 0 {
		val, err := strconv.ParseBool(r.URL.Query().Get("admin"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse bool: %v", err), http.StatusBadRequest)
			return
		}

		if val {
			cookie, err := r.Cookie("simple_stack_token")
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, fmt.Sprintf("failed to authenticate: %v", err), http.StatusUnauthorized)
				return
			} else if err != nil {
				http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
				return
			}

			if cookie.Value != os.Getenv("SIMPLE_STACK_TOKEN") {
				http.Error(w, fmt.Sprintf("failed to authenticate: %v", invalidToken), http.StatusUnauthorized)
				return
			}

			isAdmin = true
		}
	}

	t, err := template.New("articles").ParseFiles(filepath.Join("views", "articles", "articles.tmpl"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse template: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	data := struct{Articles []Article; IsAdmin bool}{Articles: articles, IsAdmin: isAdmin}
	if err = t.Execute(&buf, data); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	if r.Form.Get("username") != os.Getenv("USER_NAME") ||
		r.Form.Get("password") != os.Getenv("USER_PASS") {
		http.Error(w, fmt.Sprintf("failed to authenticate: %v", invalidCredentials), http.StatusUnauthorized)
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "simple_stack_token"
	cookie.Value = os.Getenv("SIMPLE_STACK_TOKEN")
	cookie.Expires = time.Now().Add(time.Hour * 1)
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.Path = "/"
	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/dashboard")
}

func InitEndpoints(mux *http.ServeMux) {
	/*** Private **/
	mux.Handle("/post/article", CheckCookie(http.HandlerFunc(PostArticle)))
	mux.Handle("/put/article", CheckCookie(http.HandlerFunc(PutArticle)))
	mux.Handle("/delete/article", CheckCookie(http.HandlerFunc(DeleteArticle)))

	/*** Public ***/
	mux.HandleFunc("/authenticate", Authenticate)
	mux.HandleFunc("/get/articles", GetArticles)
}

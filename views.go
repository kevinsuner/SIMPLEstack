package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
)

func ViewArticle(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/article/")
	if len(slug) < 0 {
		http.Error(w, fmt.Sprintf("failed to trim slug: %v", errEmptyString), http.StatusBadRequest)
		return
	}

	var article Article
	err := db.QueryRow(
		`select created_at, updated_at, title, description, author, status, content from "articles" where slug = $1`, slug).Scan(
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.Title,
		&article.Description,
		&article.Author,
		&article.Status,
		&article.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get article: %v", err), http.StatusInternalServerError)
		return
	}

	t, err := template.New("article").ParseFiles(
		filepath.Join("views", "layouts", "header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "articles", "article.tmpl"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = goldmark.Convert([]byte(article.Content), &buf); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse markdown: %v", err), http.StatusInternalServerError)
		return
	}
	html := buf.String() 

	buf = bytes.Buffer{}
	if err = t.Execute(&buf, map[string]interface{}{
		"meta": Meta{
			Description: article.Description,
			Author: article.Author,
			Type: "article",
			URL: "https://"+r.Host,
			Title: fmt.Sprintf("%s | SIMPLEstack", article.Title),
			CreatedAt: article.CreatedAt.String,
			UpdatedAt: article.UpdatedAt.String,
		},
		"article": article,
		"html": template.HTML(html),
	}); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func EditArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse id: %v", err), http.StatusBadRequest)
		return
	}

	var article Article
	err = db.QueryRow(
		`select id, title, slug, description, author, status, content from "articles" where id = $1`, id).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Description,
		&article.Author,
		&article.Status,
		&article.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get article: %v", err), http.StatusInternalServerError)
		return
	}

	t, err := template.New("edit").ParseFiles(
		filepath.Join("views", "layouts", "admin_header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "articles", "edit.tmpl"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, map[string]interface{}{
		"article": article,
	}); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("create").ParseFiles(
		filepath.Join("views", "layouts", "admin_header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "articles", "create.tmpl"),
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, nil); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin_token")
	if errors.Is(err, http.ErrNoCookie) {
		t, err := template.New("login").ParseFiles(
			filepath.Join("views", "layouts", "admin_header.tmpl"),
			filepath.Join("views", "layouts", "navbar.tmpl"),
			filepath.Join("views", "layouts", "footer.tmpl"),
			filepath.Join("views", "admin", "login.tmpl"),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		if err = t.Execute(&buf, nil); err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cookie: %v", err), http.StatusInternalServerError)
		return
	}

	if cookie.Value != os.Getenv("ADMIN_TOKEN") {
		t, err := template.New("login").ParseFiles(
			filepath.Join("views", "layouts", "admin_header.tmpl"),
			filepath.Join("views", "layouts", "navbar.tmpl"),
			filepath.Join("views", "layouts", "footer.tmpl"),
			filepath.Join("views", "admin", "login.tmpl"),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		if err = t.Execute(&buf, nil); err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
		return
	}

	var pages float64
	err = db.QueryRow(fmt.Sprintf(`
		select count(id)::float / %d as pages from "articles"`, ARTICLES_LIMIT)).Scan(&pages)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to scan value: %v", err), http.StatusInternalServerError)
		return
	}

	t, err := template.New("dashboard").Funcs(templateFuncs).ParseFiles(
		filepath.Join("views", "layouts", "admin_header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "admin", "dashboard.tmpl"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, map[string]interface{}{
		"pages": int(math.Ceil(pages)),
	}); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	var pages float64
	err := db.QueryRow(fmt.Sprintf(`
		select count(id)::float / %d as pages from "articles" where status = 'published'`, ARTICLES_LIMIT)).Scan(&pages)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to scan value: %v", err), http.StatusInternalServerError)
		return
	}

	t, err := template.New("homepage").Funcs(templateFuncs).ParseFiles(
		filepath.Join("views", "layouts", "header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "homepage.tmpl"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, map[string]interface{}{
		"meta": Meta{
			Description: "unimplemented!",
			Author: "Kevin Suñer",
			Type: "website",
			URL: "https://%s"+r.Host,
			Title: "Home | SIMPLEstack",
		},
		"pages": int(math.Ceil(pages)),
	}); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func InitViews(mux *http.ServeMux) {
	/*** Private **/
	mux.HandleFunc(fmt.Sprintf("/%s", os.Getenv("ADMIN_URL")), AdminDashboard)
	mux.Handle("/create/article", CheckCookie(http.HandlerFunc(CreateArticle)))
	mux.Handle("/edit/article", CheckCookie(http.HandlerFunc(EditArticle)))

	/*** Public ***/
	mux.HandleFunc("/", Homepage)
	mux.HandleFunc("/article/", ViewArticle)
}

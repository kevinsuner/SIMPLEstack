package main

import (
	"bytes"
	"fmt"
	"html/template"
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
		http.Error(w, fmt.Sprintf("failed to trim slug: %v", emptyString), http.StatusBadRequest)
		return
	}

	var article Article
	err := db.QueryRow(
		`SELECT created_at, updated_at, title, excerpt, author, status, content FROM articles WHERE slug = $1`, slug).Scan(
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.Title,
		&article.Excerpt,
		&article.Author,
		&article.Status,
		&article.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get article: %v", err), http.StatusInternalServerError)
		return
	}

	t, err := template.New("view_article").ParseFiles(
		filepath.Join("views", "layouts", "header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "articles", "view_article.tmpl"),
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
	data := struct{Article Article; HTML template.HTML}{Article: article, HTML: template.HTML(html)}
	if err = t.Execute(&buf, data); err != nil {
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
		`SELECT id, title, slug, excerpt, author, status, content FROM articles WHERE id = $1`, id).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Excerpt,
		&article.Author,
		&article.Status,
		&article.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get article: %v", err), http.StatusInternalServerError)
		return
	}

	t, err := template.New("edit_article").ParseFiles(
		filepath.Join("views", "layouts", "admin_header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "admin", "edit_article.tmpl"),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse templates: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, article); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("create_article").ParseFiles(
		filepath.Join("views", "layouts", "admin_header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "admin", "create_article.tmpl"),
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
	t, err := template.New("admin_dashboard").ParseFiles(
		filepath.Join("views", "layouts", "admin_header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "admin", "admin_dashboard.tmpl"),
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

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("admin_login").ParseFiles(
		filepath.Join("views", "layouts", "admin_header.tmpl"),
		filepath.Join("views", "layouts", "navbar.tmpl"),
		filepath.Join("views", "layouts", "footer.tmpl"),
		filepath.Join("views", "admin", "admin_login.tmpl"),
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

func Homepage(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("homepage").ParseFiles(
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
	if err = t.Execute(&buf, nil); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func InitViews(mux *http.ServeMux) {
	/*** Private **/
	mux.HandleFunc(fmt.Sprintf("/%s/login", os.Getenv("ADMIN_URL")), AdminLogin)
	mux.Handle(fmt.Sprintf("/%s/dashboard", os.Getenv("ADMIN_URL")), CheckCookie(http.HandlerFunc(AdminDashboard)))
	mux.Handle(fmt.Sprintf("/%s/create/article", os.Getenv("ADMIN_URL")), CheckCookie(http.HandlerFunc(CreateArticle)))
	mux.Handle(fmt.Sprintf("/%s/edit/article", os.Getenv("ADMIN_URL")), CheckCookie(http.HandlerFunc(EditArticle)))

	/*** Public ***/
	mux.HandleFunc("/", Homepage)
	mux.HandleFunc("/article/", ViewArticle)
}

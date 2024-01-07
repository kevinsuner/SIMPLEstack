package main

import (
	"database/sql"
	"html/template"
)

type Article struct {
	ID uint
	CreatedAt sql.NullString
	UpdatedAt sql.NullString
	Title string
	Slug string
	Description string
	Author string
	Status string
	Content string
}

type Meta struct {
	Description string
	Author string
	Type string
	URL string
	Title string
	CreatedAt string
	UpdatedAt string
}

type Pagination struct {
	Offset int
}

type TemplateData struct {
	Meta Meta
	Article Article
	Articles []Article
	Pagination []Pagination
	HTML template.HTML
	IsAdmin bool
}

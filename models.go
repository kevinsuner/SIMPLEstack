package main

import "database/sql"

type Article struct {
	ID uint
	CreatedAt sql.NullString
	UpdatedAt sql.NullString
	Title string
	Slug string
	Excerpt string
	Author string
	Status string
	Content string
}

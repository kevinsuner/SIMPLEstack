package main

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

/*** data ***/

type Task struct {
	ID uint
	Title string
}

var id uint = 1
var tasks = []Task{
	{ID: id, Title: "An example task"},
}

/*** operations ***/

func addTask(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("failed to parse form", err.Error())
		return
	}

	id++ // Yes I'm lazy, this is not a proper ID :)
	tasks = append(tasks, Task{ID: id, Title: r.Form.Get("title")})

	t, err := template.New("tasks").Parse(`
		{{range .Tasks}}
			<input name="title" type="text" value="{{.Title}}" hx-put="/update?id={{.ID}}" hx-trigger="keyup changed delay:2000ms" hx-target="#tasks" hx-swap="innerHTML">
			<button hx-delete="/delete?id={{.ID}}" hx-trigger="click" hx-target="#tasks" hx-swap="innerHTML">
				Remove task!
			</button>
			<br />
		{{end}}
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to parse tpl", err.Error())
		return
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, struct{Tasks []Task}{Tasks: tasks})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute tpl", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("failed to parse form", err.Error())
		return
	}

	for idx, v := range tasks {
		if strconv.Itoa(int(v.ID)) == r.URL.Query().Get("id") {
			tasks[idx].Title = r.Form.Get("title")
			break;
		}
	}

	t, err := template.New("tasks").Parse(`
		{{range .Tasks}}
			<input name="title" type="text" value="{{.Title}}" hx-put="/update?id={{.ID}}" hx-trigger="keyup changed delay:2000ms" hx-target="#tasks" hx-swap="innerHTML">
			<button hx-delete="/delete?id={{.ID}}" hx-trigger="click" hx-target="#tasks" hx-swap="innerHTML">
				Remove task!
			</button>
			<br />
		{{end}}
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to parse tpl", err.Error())
		return
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, struct{Tasks []Task}{Tasks: tasks})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute tpl", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	for idx, v := range tasks {
		if strconv.Itoa(int(v.ID)) == r.URL.Query().Get("id") {
			// Yeet it out of existence
			tasks = append(tasks[:idx], tasks[idx+1:]...)
			break;
		}
	}

	t, err := template.New("tasks").Parse(`
		{{range .Tasks}}
			<input name="title" type="text" value="{{.Title}}" hx-put="/update?id={{.ID}}" hx-trigger="keyup changed delay:2000ms" hx-target="#tasks" hx-swap="innerHTML">
			<button hx-delete="/delete?id={{.ID}}" hx-trigger="click" hx-target="#tasks" hx-swap="innerHTML">
				Remove task!
			</button>
			<br />
		{{end}}
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to parse tpl", err.Error())
		return
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, struct{Tasks []Task}{Tasks: tasks})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute tpl", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("tasks").Parse(`
		<!DOCTYPE html>
		<html>
			<head>
				<script src="https://unpkg.com/htmx.org@1.9.10" integrity="sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC" crossorigin="anonymous"></script>
			</head>
			<body>
				<div id="tasks">
					{{range .Tasks}}
						<input name="title" type="text" value="{{.Title}}" hx-put="/update?id={{.ID}}" hx-trigger="keyup changed delay:2000ms" hx-target="#tasks" hx-swap="innerHTML">
						<button hx-delete="/delete?id={{.ID}}" hx-trigger="click" hx-target="#tasks" hx-swap="innerHTML">
							Remove task!
						</button>
						<br />
						{{else}}
						<h4>No tasks found</h4>
					{{end}}
				</div>
				<form hx-post="/add" hx-target="#tasks" hx-swap="innerHTML">
					<input id="title" name="title" type="text">
					<input type="submit" value="Add task!">
				</form>
			<body>
		</html>
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to parse tpl", err.Error())
		return
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, struct{Tasks []Task}{Tasks: tasks})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute tpl", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

/*** init ***/

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/add", addTask)
	mux.HandleFunc("/update", updateTask)
	mux.HandleFunc("/delete", deleteTask)
	http.ListenAndServe(":8080", mux)
}

package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = "5432"
	user   = "vivek"
	dbname = "vivek"
)

var tpl = template.Must(template.ParseGlob("templates/*"))

type Todo struct {
	Id          int
	Title       string
	Description string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/delete/", deleteTodo)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		if r.PostForm.Get("todo_title") == "" {
			fmt.Fprintf(w, "please provide the title")
		} else {
			postgresqlExecute(fmt.Sprintf("insert into todos(title, description) values('%v', '%v')", r.PostForm.Get("todo_title"), r.PostForm.Get("todo_desc")))
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		rows := postgresqlQuery("select id, title, description from todos")
		defer rows.Close()
		var (
			id          int
			title       string
			description string
			todos       []Todo
		)
		for rows.Next() {
			err := rows.Scan(&id, &title, &description)
			checkError(err)
			todos = append(todos, Todo{id, title, description})
		}
		tpl.ExecuteTemplate(w, "index.gohtml", todos)
	}
}
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	postgresqlExecute(fmt.Sprintf("delete from todos where id = %v", id))
	http.Redirect(w, r, "/", http.StatusFound)
}

func postgresqlExecute(command string) {
	psqlConnection := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, os.Getenv("PG_PASS"), dbname)
	db, err := sql.Open("postgres", psqlConnection)
	checkError(err)
	defer db.Close()
	db.Exec(command)
}

func postgresqlQuery(query string) *sql.Rows {
	psqlConnection := fmt.Sprintf("host=localhost port=5432 user=vivek password=%v dbname=vivek sslmode=disable", os.Getenv("PG_PASS"))
	db, err := sql.Open("postgres", psqlConnection)
	checkError(err)
	rows, err := db.Query(query)
	checkError(err)
	return rows
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

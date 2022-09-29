package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"my-project/connection"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	connection.ConnectDatabase()

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/project", project).Methods("GET")
	r.HandleFunc("/add-project", addProject).Methods("POST")
	r.HandleFunc("/detail/{id}", detail).Methods("GET")
	r.HandleFunc("/delete/{id}", delete).Methods("GET")
	r.HandleFunc("/edit/{id}", update).Methods("GET")
	r.HandleFunc("/edit-project/{id}", editProject).Methods("POST")

	fmt.Println("server on in port 5000")
	http.ListenAndServe("localhost:5000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, description, technologies, duration FROM blog ORDER BY id DESC")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.Name, &each.Desc, &each.Technologies, &each.Duration)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	card := map[string]interface{}{
		"Add": result,
	}

	tmpl.Execute(w, card)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	tmpl.Execute(w, "")
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/addProject.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	tmpl.Execute(w, "")
}

type Project struct {
	ID                int
	Name              string
	Start_date        time.Time
	End_date          time.Time
	StartDate         string
	EndDate           string
	Duration          string
	Desc              string
	Technologies      []string
	Format_Start_date string
	Format_End_date   string
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var start_date = r.PostForm.Get("startDate")
	var end_date = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("desc")
	var technologies []string
	technologies = r.Form["technologies"]

	layout := "2006-01-02"
	dateStart, _ := time.Parse(layout, start_date)
	dateEnd, _ := time.Parse(layout, end_date)

	hours := dateEnd.Sub(dateStart).Hours()
	daysInHours := hours / 24
	monthInDay := math.Round(daysInHours / 30)
	yearInMonth := math.Round(monthInDay / 12)

	var duration string

	if yearInMonth > 0 {
		duration = strconv.FormatFloat(yearInMonth, 'f', 0, 64) + " Years"
	} else if monthInDay > 0 {
		duration = strconv.FormatFloat(monthInDay, 'f', 0, 64) + " Months"
	} else if daysInHours > 0 {
		duration = strconv.FormatFloat(daysInHours, 'f', 0, 64) + " Days"
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
	} else {
		duration = "0 Days"
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO blog(name, start_date, end_date, description, technologies, duration) VALUES ($1,$2,$3,$4,$5,$6)", name, start_date, end_date, desc, technologies, duration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/detail.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Detail = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, duration FROM blog WHERE id=$1", id).Scan(
		&Detail.ID, &Detail.Name, &Detail.Start_date, &Detail.End_date, &Detail.Desc, &Detail.Technologies, &Detail.Duration)

	// fmt.Printf("%T %v", id, id)
	// fmt.Println(err)
	fmt.Println(Detail)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	Detail.Format_Start_date  = Detail.Start_date.Format("2 January 2006")
	Detail.Format_End_date  = Detail.End_date.Format("2 January 2006")

	data := map[string]interface{}{
		"Details": Detail,
	}
	// fmt.Println(data)
	tmpl.Execute(w, data)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM blog WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var name = r.PostForm.Get("inputName")
	var start_date = r.PostForm.Get("startDate")
	var end_date = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("desc")
	var technologies []string
	technologies = r.Form["technologies"]

	layout := "2006-01-02"
	dateStart, _ := time.Parse(layout, start_date)
	dateEnd, _ := time.Parse(layout, end_date)

	hours := dateEnd.Sub(dateStart).Hours()
	daysInHours := hours / 24
	monthInDay := math.Round(daysInHours / 30)
	yearInMonth := math.Round(monthInDay / 12)

	var duration string

	if yearInMonth > 0 {
		duration = strconv.FormatFloat(yearInMonth, 'f', 0, 64) + " Years"
	} else if monthInDay > 0 {
		duration = strconv.FormatFloat(monthInDay, 'f', 0, 64) + " Months"
	} else if daysInHours > 0 {
		duration = strconv.FormatFloat(daysInHours, 'f', 0, 64) + " Days"
	} else {
		duration = "0 Days"
	}

	_, err = connection.Conn.Exec(context.Background(), "UPDATE blog SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, duration=$6 WHERE id=$7", name, dateStart, dateEnd, desc, technologies, duration, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/editProject.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Edit = Project{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, duration FROM blog WHERE id=$1", id).Scan(
		&Edit.ID, &Edit.Name, &Edit.Start_date, &Edit.End_date, &Edit.Desc, &Edit.Technologies, &Edit.Duration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	Edit.Format_Start_date  = Edit.Start_date.Format("2006-01-02")
	Edit.Format_End_date  = Edit.End_date.Format("2006-01-02")

	data := map[string]interface{}{
		"Id":   id,
		"Edit": Edit,
	}
	// fmt.Println(data)
	tmpl.Execute(w, data)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

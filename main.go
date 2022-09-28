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

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, duration FROM blog ORDER BY id")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.Name, &each.Start_date, &each.End_date, &each.Desc, &each.Technologies, &each.Duration)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.StartDate = each.Start_date.Format("2 January 2006")
		each.EndDate = each.End_date.Format("2 January 2006")

		result = append(result, each)
	}

	fmt.Println(result)

	card := map[string]interface{}{
		"Add": result,
	}
	w.WriteHeader(http.StatusOK)
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
	ID           int
	Name         string
	Start_date   time.Time
	End_date     time.Time
	Duration     string
	Desc         string
	StartDate    string
	EndDate      string
	Technologies []string
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
	technologies := r.Form["technologies"]

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
		&Detail.ID, &Detail.Name, &Detail.StartDate, &Detail.EndDate, &Detail.Desc, &Detail.Technologies, &Detail.Duration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	Detail.StartDate = Detail.Start_date.Format("2 January 2006")
	Detail.EndDate = Detail.End_date.Format("2 January 2006")

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
	technologies := r.Form["technologies"]

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

	_, err = connection.Conn.Exec(context.Background(), "UPDATE blog SET id=$1, name=$2, start_date=$3, end_date=$4, description=$5, technologies=$6, duration=$7 WHERE id=$1", id, name, start_date, end_date, desc, technologies, duration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
		&Edit.ID, &Edit.Name, &Edit.StartDate, &Edit.EndDate, &Edit.Desc, &Edit.Technologies, &Edit.Duration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	Edit.StartDate = Edit.Start_date.Format("2 January 2006")
	Edit.EndDate = Edit.End_date.Format("2 January 2006")

	data := map[string]interface{}{
		"Id":   id,
		"Edit": Edit,
	}
	// fmt.Println(data)
	tmpl.Execute(w, data)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

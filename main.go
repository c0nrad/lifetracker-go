package main

import (
	"fmt"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type Accomplishment struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Body string
	Name string
}

var session *mgo.Session
var templates = template.Must(template.ParseFiles("./templates/base.html", "./templates/index.html", "./templates/accomplishment.html", "./templates/addAccomplishment.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving / %v\n", r.Method)
	recent := findRecent()

	err := templates.ExecuteTemplate(w, "base.html", recent)
	if err != nil {
		panic(err)
	}
}

func accomplishmentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving", r.URL, r.Method, "\n")

	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println(r.Form, r.Form["accomplishment"])
		c := session.DB("test").C("accomplishments")
		err := c.Insert(&Accomplishment{"", r.Form["accomplishment"][0], r.Form["name"][0]})
		if err != nil {
			panic(err)
		}
	} else {

	}
	fmt.Fprintf(w, "Good Work!")
}

func initMGO() *mgo.Session {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	return session
}

func findRecent() []Accomplishment {
	var results []Accomplishment
	err := session.DB("test").C("accomplishments").Find(bson.M{}).All(&results)

	if err != nil {
		panic(err)
	}

	return results
}

func main() {
	fmt.Println("LOL NO WAI")
	session = initMGO()

	http.Handle("/static/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/a/", accomplishmentHandler)

	fmt.Println("Starting server on port :7774")
	http.ListenAndServe(":7774", nil)
}

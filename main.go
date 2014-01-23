package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"labix.org/v2/mgo"
	"net/http"
	"os"
)

var mongoSession *mgo.Session
var currentUser User
var templates = template.Must(template.ParseFiles("./templates/calendar.html", "./templates/howItWorks.html", "./templates/signup.html", "./templates/login.html", "./templates/base.html", "./templates/index.html", "./templates/accomplishmentTemplate.html", "./templates/addAccomplishment.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	recent := findRecent()
	recentJSON, err := json.Marshal(recent)

	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"Recent":     recent,
		"RecentJSON": string(recentJSON),
		"User":       currentUser,
	}

	if err := templates.ExecuteTemplate(w, "base.html", params); err != nil {
		return err
	}
	return nil
}

func initMGO() *mgo.Session {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
	}

	err = session.DB("test").C("users").EnsureIndex(index)
	return session
}

func makeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		currentUser = user
		fmt.Println("Serving:", r.Method, r.URL)
		sessionCookie, err := r.Cookie("ltsession")
		if err != nil {
		} else {
			currentUser = sessionTokenToUser(sessionCookie.Value)
		}

		if err = fn(w, r); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}

func main() {
	mongoSession = initMGO()

	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", makeHandler(indexHandler))
	http.HandleFunc("/a", makeHandler(accomplishmentHandler))
	http.HandleFunc("/login", makeHandler(loginHandler))
	http.HandleFunc("/signup", makeHandler(signupHandler))
	http.HandleFunc("/logout", makeHandler(logoutHandler))

	if len(os.Args) >= 2 {
		fmt.Println("Starting server on port", os.Args[1])
		http.ListenAndServe(":"+os.Args[1], nil)
	} else {
		fmt.Println("Starting server on port :7776")
		http.ListenAndServe(":7776", nil)
	}
}

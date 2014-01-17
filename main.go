package main

import (
	"fmt"
	"html/template"
	"labix.org/v2/mgo"
	"net/http"
	"os"
)

var mongoSession *mgo.Session
var currentUser User
var templates = template.Must(template.ParseFiles("./templates/howItWorks.html", "./templates/signup.html", "./templates/login.html", "./templates/base.html", "./templates/index.html", "./templates/accomplishment.html", "./templates/addAccomplishment.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	recent := findRecent()

	params := map[string]interface{}{
		"Recent": recent,
		"User":   currentUser,
	}

	if err := templates.ExecuteTemplate(w, "base.html", params); err != nil {
		return err
	}
	return nil
}

func emailHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	f, err := os.OpenFile("emails.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(r.Form["email"][0] + "\n"); err != nil {
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

	return session
}

func loadUser(r *http.Request) {
	sessionCookie, err := r.Cookie("ltsession")
	if err != nil {
		fmt.Println("loadUser: ", err)
	} else {
		currentUser = sessionTokenToUser(sessionCookie.Value)
		if err != nil {
			fmt.Println("loadUser:", err)
			currentUser = User{}
		}
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/email", makeHandler(emailHandler))
	http.HandleFunc("/", makeHandler(indexHandler))
	http.HandleFunc("/a/", makeHandler(accomplishmentHandler))
	http.HandleFunc("/login", makeHandler(loginHandler))
	http.HandleFunc("/signup", makeHandler(signupHandler))
	http.HandleFunc("/logout", makeHandler(logoutHandler))

	fmt.Println("Starting server on port :7776")
	http.ListenAndServe(":7776", nil)
}

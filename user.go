package main

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
	"net/http"
	"regexp"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string
	Email    string
	Password string
}

const RE_BASIC_EMAIL = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`

func authenticate(email, password string) (User, error) {
	var user User
	err := mongoSession.DB("test").C("users").Find(bson.M{"email": email, "password": password}).One(&user)
	return user, err
}

func validateEmail(email string) bool {
	exp, _ := regexp.Compile(RE_BASIC_EMAIL)
	return exp.MatchString(email)
}

func addUser(name, email, password string) (User, error) {
	var newUser User
	if name == "" {
		return newUser, errors.New("Name too small")
	}
	if len(password) < 5 {
		return newUser, errors.New("Password too short")
	}
	if !validateEmail(email) {
		return newUser, errors.New("Invalid email address")
	}

	newUser = User{"", name, email, password}
	err := mongoSession.DB("test").C("users").Insert(newUser)
	return newUser, err
}

func isEmptyUser(a User) bool {
	var b User
	return a == b
}

func login(email, password string, w http.ResponseWriter) error {
	currentUser, err := authenticate(email, password)
	if err != nil {
		return err
	}
	s, err := generateSession(currentUser)
	if err != nil {
		fmt.Println("loginHandler: error generating session")
		return err
	}
	addSessionCookie(s, &w)
	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		r.ParseForm()
		_ = login(r.FormValue("email"), r.FormValue("password"), w)
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func logoutHandler(w http.ResponseWriter, r *http.Request) error {
	var user User

	if currentUser == user {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	if err := destroySession(currentUser); err != nil {
		return err
	}

	currentUser = user
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func signupHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		if err := templates.ExecuteTemplate(w, "signup.html", nil); err != nil {
			return err
		}
	} else {
		r.ParseForm()
		err := errors.New("")
		currentUser, err = addUser(r.FormValue("name"), r.FormValue("email"), r.FormValue("password"))
		if err != nil {
			return err
		}

		err = login(r.FormValue("email"), r.FormValue("password"), w)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
	return nil
}

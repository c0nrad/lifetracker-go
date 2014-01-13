package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"time"
)

type Accomplishment struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Body      string
	Name      string
	Date      time.Time
	ImagePath string
}

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Email    string
	Password string
}

var session *mgo.Session
var templates = template.Must(template.ParseFiles("./templates/base.html", "./templates/index.html", "./templates/accomplishment.html", "./templates/addAccomplishment.html"))

func generateImageHash(filename string) string {
	h := md5.New()
	io.WriteString(h, filename)
	io.WriteString(h, time.Now().String())
	return hex.EncodeToString(h.Sum(nil))
}

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
		filename := ""
		r.ParseForm()
		fmt.Println(r.Form, r.Form["accomplishment"])

		file, handler, err := r.FormFile("file")
		if err != nil {
			// No file upload
		} else {
			data, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err, "2")
			}

			filename = "./static/uploads/" + generateImageHash(handler.Filename)
			err = ioutil.WriteFile(filename, data, 0777)
			if err != nil {
				fmt.Println(err)
			}
		}

		c := session.DB("test").C("accomplishments")
		newAccomplishment := &Accomplishment{"", r.Form["accomplishment"][0], r.Form["name"][0], time.Now(), filename}
		fmt.Println("Inserting new accomplishment:", newAccomplishment)
		err = c.Insert(newAccomplishment)
		if err != nil {
			panic(err)
		}
	} else {

	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func emailHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	f, err := os.OpenFile("emails.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(r.Form["email"][0] + "\n"); err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Email saved!")
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
	err := session.DB("test").C("accomplishments").Find(bson.M{}).Sort("-date").All(&results)

	if err != nil {
		panic(err)
	}

	return results
}

func main() {
	fmt.Println("LOL NO WAI")
	session = initMGO()

	http.Handle("/static/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/email", emailHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/a/", accomplishmentHandler)

	fmt.Println("Starting server on port :7776")
	http.ListenAndServe(":7776", nil)
}

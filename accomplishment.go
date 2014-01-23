package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type Accomplishment struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Body      string
	Name      string
	Date      time.Time
	ImagePath string
	UserID    bson.ObjectId ",omitempty"
}

func generateImageHash(filename string) string {
	h := md5.New()
	io.WriteString(h, filename)
	io.WriteString(h, time.Now().String())
	return hex.EncodeToString(h.Sum(nil))
}

func saveImage(r *http.Request) string {
	filename := ""
	file, handler, err := r.FormFile("file")
	if err != nil {
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("saveImage:", err)
			return ""
		}

		filename = "./static/uploads/" + generateImageHash(handler.Filename)
		if err = ioutil.WriteFile(filename, data, 0777); err != nil {
			fmt.Println("saveImage:", err)
			return ""
		}
	}
	return filename
}

func validateAccomplishment(a *Accomplishment) error {
	if a.Body == "" {
		return errors.New("There must be an accomplishment")
	}

	if len(a.Body) > 200 {
		a.Body = a.Body[0:200]
	}

	if a.Name == "" {
		a.Name = "Anonymous"
	}

	return nil
}

func buildAccomplishment(accomplishment, name, filename string) (*Accomplishment, error) {
	newAccomplishment := &Accomplishment{"", accomplishment, name, time.Now(), filename, ""}

	if isEmptyUser(currentUser) {
		newAccomplishment.UserID = ""
	} else {
		newAccomplishment.UserID = currentUser.ID
		newAccomplishment.Name = currentUser.Name
	}

	if err := validateAccomplishment(newAccomplishment); err != nil {
		return nil, err
	}

	return newAccomplishment, nil
}

func accomplishmentHandler(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "POST" {
		r.ParseForm()

		filename := saveImage(r)
		accomplishment := r.FormValue("accomplishment")
		name := r.FormValue("name")

		newAccomplishment, err := buildAccomplishment(accomplishment, name, filename)
		if err != nil {
			return err
		}

		fmt.Println("Inserting new accomplishment:", newAccomplishment)
		if err := mongoSession.DB("test").C("accomplishments").Insert(newAccomplishment); err != nil {
			return err
		}
	} else if r.Method == "GET" {

	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func findRecent() []Accomplishment {
	var results []Accomplishment

	err := errors.New("")
	if isEmptyUser(currentUser) {
		err = mongoSession.DB("test").C("accomplishments").Find(bson.M{"userid": bson.M{"$exists": false}}).Sort("-date").All(&results)
	} else {
		err = mongoSession.DB("test").C("accomplishments").Find(bson.M{"userid": currentUser.ID}).Sort("-date").All(&results)
	}

	if err != nil {
		panic(err)
	}

	return results
}

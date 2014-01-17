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

func accomplishmentHandler(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "POST" {
		filename := ""
		r.ParseForm()
		fmt.Println(r.Form, r.Form["accomplishment"])

		file, handler, err := r.FormFile("file")
		if err != nil {
		} else {
			data, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			filename = "./static/uploads/" + generateImageHash(handler.Filename)
			if err = ioutil.WriteFile(filename, data, 0777); err != nil {
				return err
			}
		}

		fmt.Println(r.Form)

		accomplishment := r.FormValue("accomplishment")
		name := r.FormValue("name")

		if accomplishment == "" {
			return errors.New("There must be an Accomplishment")
		}

		if name == "" {
			name = "Anonymous"
		}

		var newAccomplishment *Accomplishment
		if isEmptyUser(currentUser) {
			newAccomplishment = &Accomplishment{"", accomplishment, name, time.Now(), filename, ""}
		} else {
			newAccomplishment = &Accomplishment{"", accomplishment, name, time.Now(), filename, currentUser.ID}
		}
		fmt.Println("Inserting new accomplishment:", newAccomplishment)
		if err = mongoSession.DB("test").C("accomplishments").Insert(newAccomplishment); err != nil {
			fmt.Println("AQUI AQUI")
			return err
		}
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

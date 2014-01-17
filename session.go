package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type Session struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Token  string
	UserID bson.ObjectId //Index me
}

func generateSession(user User) (Session, error) {
	h := md5.New()
	io.WriteString(h, time.Now().String())
	io.WriteString(h, user.Email+user.Password+string(user.ID))
	token := hex.EncodeToString(h.Sum(nil))

	newSession := Session{"", token, user.ID}
	err := mongoSession.DB("test").C("sessions").Insert(newSession)
	return newSession, err
}

func sessionTokenToSession(sessionToken string) (Session, error) {
	var s Session
	err := mongoSession.DB("test").C("sessions").Find(bson.M{"token": sessionToken}).One(&s)
	return s, err
}

func destroySession(user User) error {
	fmt.Println(user)
	return mongoSession.DB("test").C("sessions").Remove(bson.M{"userid": user.ID})
}

func sessionTokenToUser(sessionToken string) User {
	var user User
	s, err := sessionTokenToSession(sessionToken)
	if err != nil {
		return user
	}
	user, err = sessionToUser(s)
	return user
}

func sessionToUser(s Session) (User, error) {
	var user User
	err := mongoSession.DB("test").C("users").Find(bson.M{"_id": s.UserID}).One(&user)
	return user, err
}
func userToSession(user User) (Session, error) {
	var s Session
	err := mongoSession.DB("test").C("sessions").Find(bson.M{"userid": user.ID}).One(&s)
	return s, err
}

func addSessionCookie(s Session, w *http.ResponseWriter) {
	cookie := &http.Cookie{Name: "ltsession", Value: s.Token, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(*w, cookie)
}

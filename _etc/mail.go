package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

var email = os.Getenv("AOP_EMAIL_ADDRESS")
var login = os.Getenv("AOP_EMAIL_LOGIN")
var password = os.Getenv("AOP_EMAIL_PASSWORD")

var expectedScore = "7";

func sendMail(visitorEmail string, body string) error {
	auth := smtp.PlainAuth(
		"",
		login,
		password,
		"smtp.gmail.com",
	)

	msg := strings.Join([]string{
		fmt.Sprintf("To: %s", email),
		fmt.Sprintf("Reply-To: %s", visitorEmail),
		fmt.Sprintf("Subject: Message de %s sur ahouhpuc.fr", visitorEmail),
		"",
		body,
	}, "\r\n")

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		email,
		[]string{email},
		[]byte(msg),
	)
}

func mailHandler(w http.ResponseWriter, r *http.Request) {
	visitorEmail := r.PostFormValue("email")
	body := r.PostFormValue("body")
	redirect := r.PostFormValue("redirect")
	score := r.PostFormValue("score") // Our score input is filtering spambots (a bit)

	if len(visitorEmail) > 0 && len(body) > 0 && len(redirect) > 0 {
		var err = error(nil)
		if (score == expectedScore) {
			err = sendMail(visitorEmail, body)
		}
		if err == nil {
			http.Redirect(w, r, redirect, http.StatusFound)
		} else {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}

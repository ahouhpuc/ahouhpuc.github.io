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

func sendMail(replyTo string, body string) error {
	auth := smtp.PlainAuth(
		"",
		login,
		password,
		"smtp.gmail.com",
	)

	msg := strings.Join([]string{
		fmt.Sprintf("To: %s", email),
		fmt.Sprintf("Reply-To: %s", replyTo),
		"Subject: Message de ahouhpuc.fr",
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
	replyTo := r.PostFormValue("email")
	body := r.PostFormValue("body")
	redirect := r.PostFormValue("redirect")

	if len(replyTo) > 0 && len(body) > 0 && len(redirect) > 0 {
		err := sendMail(replyTo, body)
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

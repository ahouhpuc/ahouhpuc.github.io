// gofmt -tabs=false -tabwidth=2 -w=true server.go && go run server.go

package main

import (
  "fmt"
  "log"
  "net/http"
  "net/smtp"
  "os"
  "strings"
)

func sendMail(replyTo string, body string) error {
  email := os.Getenv("AOP_EMAIL_ADDRESS")
  login := os.Getenv("AOP_EMAIL_LOGIN")
  password := os.Getenv("AOP_EMAIL_PASSWORD")

  if len(email) == 0 || len(login) == 0 || len(password) == 0 {
    panic("AOP_EMAIL_ADDRESS, AOP_EMAIL_LOGIN and AOP_EMAIL_PASSWORD environment variable are mandatory.")
  }

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

func handler(w http.ResponseWriter, r *http.Request) {
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
    http.Error(w, "Not Found", http.StatusNotFound)
  }
}

func main() {
  http.HandleFunc("/", handler)
  log.Fatal(http.ListenAndServe(":5000", nil))
}

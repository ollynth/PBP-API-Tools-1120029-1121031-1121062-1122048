package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

type NormalResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Users struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func sendMessage(w http.ResponseWriter, code int, message string) {
	var response NormalResponse
	response.Status = code
	response.Message = message
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SetRedis(rdb *redis.Client, key string, value string, expiration int) {
	err := rdb.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		fmt.Println("Error: ", err.Error)
		return
	}
	log.Println("Set Redis Success")
}

func GetRedis(key string) string {
	dataGet, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println("Error: ", err.Error)
		return ""
	}
	log.Println("Get Redis Success")
	return dataGet
}

func SendMail(to string, subject string, content string) {
	// set manual
	email := "if-21062@students.ithb.ac.id"
	password := "zktpwgsklzlqcgoq"

	message := gomail.NewMessage()
	message.SetHeader("From", email)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", content)

	sender := gomail.NewDialer("smtp.gmail.com", 587, email, password)

	if err := sender.DialAndSend(message); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Email sent successfully to", to)
}

func GetAllUserData() {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT email FROM users")
	if err != nil {
		log.Println(err)
		return
	}

	var user Users
	var userEmails []string
	for rows.Next() {
		if err := rows.Scan(&user.Email); err != nil {
			log.Println(err)
			return
		} else {
			fmt.Println(user.Email)
			userEmails = append(userEmails, user.Email)
		}
	}
	for _, email := range userEmails {
		SetRedis(rdb, "userEmail", email, 0)
	}
}

func getAllUserEmailsFromRedis() []string {
	GetAllUserData()
	userEmailsMap := make(map[string]bool)
	var userEmails []string
	keys := rdb.Keys(context.Background(), "userEmail:*").Val()

	for _, key := range keys {
		email := GetRedis(key)
		if _, exists := userEmailsMap[email]; !exists {
			userEmailsMap[email] = true
		}
	}

	for email := range userEmailsMap {
		userEmails = append(userEmails, email)
	}

	for _, mail := range userEmails {
		fmt.Println("Hasil Get : ", mail)
	}
	return userEmails
}

func Announce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		sendMessage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	subject := r.FormValue("subject")
	text := r.FormValue("text")

	if subject == "" {
		sendMessage(w, http.StatusBadRequest, "Subject cannot be empty")
		return
	}

	if text == "" {
		sendMessage(w, http.StatusBadRequest, "Message cannot be empty")
		return
	}

	fmt.Println("Subject: ", subject)
	fmt.Println("Text: ", text)

	userEmails := getAllUserEmailsFromRedis()

	for _, email := range userEmails {
		fmt.Println(email)
		SendMail(email, subject, text)
	}

	sendMessage(w, http.StatusOK, "Announcement sent successfully to all users")
}

func sendEmailToAllEmployees(message string) {
	employeeEmails := getAllUserEmailsFromRedis()
	if employeeEmails == nil {
		log.Println("Error getting employee emails")
		return
	}

	for _, email := range employeeEmails {
		go SendMail(email, "Reminder", message)
		log.Println("Task Send email successfully to ", email)
	}
}

func Task() {
	c := cron.New()

	c.AddFunc("@every 10m", func() {
		sendEmailToAllEmployees("Waktunya untuk istirahat sejenak!")
	})
	c.Start()
}

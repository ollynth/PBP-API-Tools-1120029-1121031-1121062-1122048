package main

import (
	"ToolsAPI/controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/announce", controllers.Announce).Methods("POST")

	controllers.Task()

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))

	// SMTP server details
	// smtpServer := "smtp.gmail.com"
	// smtpPort := 587

	// // Sender details
	// from := "if-21062@students.ithb.ac.id"
	// password := "zktpwgsklzlqcgoq" // Use an app password if 2-factor authentication is enabled

	// // Recipient email address
	// to := "xiolala5@gmail.com"

	// // Email content
	// subject := "Test Email"
	// body := "This is a test email sent from Golang."

	// // Authenticate with the SMTP server
	// auth := smtp.PlainAuth("", from, password, smtpServer)

	// // Compose the email message
	// message := []byte("To: " + to + "\r\n" +
	// 	"Subject: " + subject + "\r\n" +
	// 	"\r\n" +
	// 	body + "\r\n")

	// // Connect to the SMTP server and send the email
	// err := smtp.SendMail(smtpServer+":"+fmt.Sprint(smtpPort), auth, from, []string{to}, message)
	// if err != nil {
	// 	fmt.Println("Error sending email:", err)
	// 	return
	// }

	// fmt.Println("Email sent successfully")
}

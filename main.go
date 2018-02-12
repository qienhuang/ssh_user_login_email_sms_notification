package main

/*
  go run main.go
  # Before you build it, don't forget to set disableLogDisplay = true
  go build main.go
  cp main /bin/sshlog
  nano /etc/profile.d/log.sh
  # Add the following line in log.sh
  #!/bin/bash
  /bin/sshlog

  # Kevin Huang at qien.huang.ny@gmail.com, 2/1/2018
*/

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os/exec"
	"strings"
)

var (

	disableLogDisplay = false	// Set to "true" before you build it
)

type Email struct {
	SmtpServerAddr string
	SmtpServerPort string
	SenderID       string
	SenderPassword string
	Receivers      []string
	Subject        string
	Body           string
}

func (s *Email) GetHost() string {
	return s.SmtpServerAddr + ":" + s.SmtpServerPort
}

func (mail *Email) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.SenderID)
	if len(mail.Receivers) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.Receivers, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	message += "\r\n" + mail.Body

	return message
}

func main() {
	if disableLogDisplay {
		log.SetOutput(ioutil.Discard)
	}

	mail := Email{}
	mail.SmtpServerAddr = "smtp.gmail.com"
	mail.SmtpServerPort = "465"
	mail.SenderID = "YourGmailAccount"		// Replace this string with your Gmail account
	mail.SenderPassword = "YourPassword"    // Replace this string with your Gmail password
	mail.Receivers = []string{ "TheReceiverEmail", "SecondReceiver"}  // Replace this string to receiver. For T-Mobile SMS, it will be xxxxxxxxxx@tmomail.net
	mail.Subject = "User just login to my Raspberry Pi"
	mail.Body = "The current logged-in users:\r\n"

	// Get current logged-in users
	out, err := exec.Command("w").Output()
	if err == nil {
		log.Println("output: " + string(out))
		mail.Body += string(out) + "********************\r\n"
	}

	// Get system host name
	out, err = exec.Command("hostname").Output()
	if err == nil {
		log.Println("Host: " + string(out))
		mail.Body += "Hostname:" + string(out)
	}

	// Get Current Wan IP
	out, err = exec.Command("curl", "ident.me").Output()

	if err == nil {
		log.Println("Current WAN IP:" + string(out))
		mail.Body += "Current WAN IP:" + string(out)
	} else {
		log.Println(err)
	}

	messageBody := mail.BuildMessage()

	//log.Println(mail.GetHost())

	//build an auth
	auth := smtp.PlainAuth("", mail.SenderID, mail.SenderPassword, mail.GetHost())

	// Gmail will reject connection if you didn't "Allow less secure apps."  https://myaccount.google.com/lesssecureapps
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         mail.GetHost(),
	}

	conn, err := tls.Dial("tcp", mail.GetHost(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, mail.GetHost())
	if err != nil {
		log.Panic(err)
	}

	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	if err = client.Mail(mail.SenderID); err != nil {
		log.Panic(err)
	}

	for _, r := range mail.Receivers {
		if err = client.Rcpt(r); err != nil {
			log.Panic(err)
		}
	}

	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}

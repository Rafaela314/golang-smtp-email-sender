package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"text/template"
)

func SendSimpleMail() error {

	// retrieve credentials from env
	from := os.Getenv("MAIL")       //{i.e: sender_example_mail_address@gmail.com}
	password := os.Getenv("PASSWD") //e-mail password

	// to is a list of receiver's addresses
	to := []string{"example_receiver_email_address@gmail.com"}

	//host is the sender mail server's address
	host := "smtp.gmail.com"

	// port is the default port of smtp server
	port := "587"

	//message is the message to be send in the mail. It must be converted into slice bytes.
	message := []byte("Simple email example message")

	//PlainAuth autenticate to host
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, to, message)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully sent email!")
	return nil
}

func SendMailWithHtmlTemplate() error {

	// retrieve credentials from env
	from := os.Getenv("MAIL")       //{i.e: sender_example_mail_address@gmail.com}
	password := os.Getenv("PASSWD") //e-mail password

	// receiver email address list
	to := []string{"example_receiver_email_address@gmail.com"}

	// smtp server configuration
	host := "smtp.gmail.com"
	port := "587"

	//PlainAuth autenticate to host
	auth := smtp.PlainAuth("", from, password, host)

	t, _ := template.ParseFiles("mailTemplate.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Sender test name",
		Message: "Test message in a HTML template",
	})

	err := smtp.SendMail(host+":"+port, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Successfully sent email!")
	return nil
}

func SendMailWithCSVAttachment() error {

	// Connect to the remote SMTP server.
	c, err := smtp.Dial("mail.example.com:25")
	if err != nil {
		return err
	}

	//close the connection after execution
	defer c.Quit()

	var body bytes.Buffer

	// Set the sender, recipients and subject first
	body.WriteString(fmt.Sprintf("From: %s\r\n", "sender_example_mail@gmail.com"))
	body.WriteString(fmt.Sprintf("To: %s\r\n", "receiver_example_mail@gmail.com"))
	body.WriteString(fmt.Sprintf("Subject: %s\r\n", "email_subject_example"))

	boundary := "my-boundary-779"
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))

	body.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	body.WriteString(fmt.Sprintf("Content-Type: text/plain; charset=\"utf-8\"\r\n"))
	content := "message that will be displayed on e-mail body"
	body.WriteString(fmt.Sprintf("\r\n%s", content))

	//attach csv file
	fileContent, err := ioutil.ReadFile("examplefile.csv")
	if err != nil {
		return err
	}

	body.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	body.WriteString(fmt.Sprintf("Content-Type: text/csv\r\n"))

	body.WriteString("Content-Transfer-Encoding: base64\r\n")
	body.WriteString("Content-Disposition: attachment; filename=CranePoRecoReport.csv\r\n")
	body.WriteString("Content-ID: <CranePoRecoReport.csv>\r\n\r\n")

	encodedBytes := make([]byte, base64.StdEncoding.EncodedLen(len(fileContent)))
	base64.StdEncoding.Encode(encodedBytes, fileContent)
	body.Write(encodedBytes)
	body.WriteString(fmt.Sprintf("\r\n--%s", boundary))

	body.WriteString("--")

	err = smtp.SendMail("example_mail@gmail.com:smtp", nil, "sender_example_mail@gmail.com", []string{"receiver_example_mail@gmail.com"}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Successfully sent email with csv attachment!")
	return nil
}

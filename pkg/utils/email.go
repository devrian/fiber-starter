package utils

import (
	"bytes"
	"fiber-starter/config"
	"fiber-starter/pkg/common"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
	"sync"
)

type MailRequest struct {
	to      []string
	subject string
	body    string
}

type MailData struct {
	Receivers []string
	Usage     string
	Data      interface{}
}

type DataEmailToken struct {
	TokenURL     string
	ExpiredTime  int
	Description  string
	IsChannelApp bool
	Title        string
}

const (
	MAIL_TEMPLATE_PATH      = "public/templates/"
	MAIL_FOR_USERACTIVATION = "user-activation"
	MAIL_FOR_RESETPASS      = "reset-password"
	MAIL_MIME               = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

var mailSubject = map[string]string{
	MAIL_FOR_USERACTIVATION: "[sample] User Activation",
	MAIL_FOR_RESETPASS:      "[sample] Reset Password",
}

var ListUsedFor = []string{MAIL_FOR_USERACTIVATION, MAIL_FOR_RESETPASS}

func (r *MailRequest) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *MailRequest) sendMailAsync(ch chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	sent := true
	mailHost := os.Getenv("MAIL_HOST")
	mailFromAddress := os.Getenv("MAIL_FROM_ADDRESS")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	for _, to := range r.to {
		email := []string{to}
		SMTP := fmt.Sprintf("%s:%d", mailHost, 587)
		body := "To: " + to + "\r\nFrom: " + mailFromAddress + "\r\nSubject: " + r.subject + "\r\n" + MAIL_MIME + "\r\n" + r.body

		if err := smtp.SendMail(SMTP, smtp.PlainAuth("", mailUsername, mailPassword, mailHost), mailFromAddress, email, []byte(body)); err != nil {
			fmt.Println(err)
			sent = false
		}
	}

	ch <- bool(sent)
}

func (r *MailRequest) sendMail() bool {
	// make a channel
	ch := make(chan bool)
	var wg sync.WaitGroup

	// do async request with channel
	wg.Add(1)
	go r.sendMailAsync(ch, &wg)

	// close the channel in the background
	go func() {
		wg.Wait()
		close(ch)
	}()

	// read from channel as they come in until its closed
	var sent bool
	for res := range ch {
		sent = res
	}

	return sent
}

func NewRequestSendMail(to []string, usedFor string) (*MailRequest, error) {
	mailReq := &MailRequest{}
	for _, mailTo := range to {
		if !common.IsEmail(mailTo) {
			return mailReq, fmt.Errorf("email %s is invalid", mailTo)
		}
	}

	mailReq.to = to
	mailReq.subject = mailSubject[usedFor]

	return mailReq, nil
}

func (r *MailRequest) Send(usedFor string, items interface{}) {
	templateName := fmt.Sprintf("%s/%s.html", MAIL_TEMPLATE_PATH, usedFor)

	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}

	if ok := r.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", r.to)
	} else {
		log.Printf("Failed to send the email to %s\n", r.to)
	}
}

func SetMailData(receivers []string, usage string, data interface{}) MailData {
	return MailData{
		Receivers: receivers,
		Usage:     usage,
		Data:      data,
	}
}

func PushQueueMail(config *config.RabbitMQ, qName string, data MailData) error {
	queue := QueueRabbitMQ{RabbitMQ: config, Data: data}
	err := PushQueueToRabbitMQ(queue, qName)

	return err
}

func CheckValidUsedFor(used string) error {
	var err error
	if !common.CheckStringContains(used, ListUsedFor) {
		usedStr := strings.Join(ListUsedFor, ", ")
		errMsg := fmt.Sprintf("invalid param type for arguments in action not includes in %s.", usedStr)
		err = fmt.Errorf("%s", errMsg)
	}

	return err
}

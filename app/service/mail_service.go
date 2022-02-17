package service

import (
	"fiber-starter/pkg/utils"
)

type MailService interface {
	MailSender(receivers []string, usage string, data interface{}) error
}

type mailService struct {
}

func NewMailService() *mailService {
	return &mailService{}
}

func (s *mailService) MailSender(receivers []string, usage string, data interface{}) error {
	emailReq, err := utils.NewRequestSendMail(receivers, usage)
	if err != nil {
		return err
	}

	emailReq.Send(usage, data)

	return err
}

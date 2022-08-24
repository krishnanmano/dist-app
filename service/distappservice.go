package service

import "dist-app/model"

type IDistAppService interface {
	SaveMessage(msg model.Message)
	GetMessages() []model.Message
}

type distAppService struct{}

func NewDistAppService() *distAppService {
	return &distAppService{}
}

func (d distAppService) SaveMessage(msg model.Message) {
	msg.SaveMessage()
}

func (d distAppService) GetMessages() []model.Message {
	return model.Message{}.GetMessages()
}

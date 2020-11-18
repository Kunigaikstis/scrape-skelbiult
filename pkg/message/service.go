package message

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Service struct {
	repo   Repository
	sender *tgbotapi.BotAPI
}

func NewService(repository Repository, botAPI *tgbotapi.BotAPI) *Service {
	return &Service{
		repo:   repository,
		sender: botAPI,
	}
}

func (s *Service) SendMessageToAllChats(message string) error {
	chats, err := s.repo.GetAll()

	if err != nil {
		return err
	}

	for _, chat := range chats {
		_, err = s.sender.Send(tgbotapi.NewMessage(chat.Id, message))

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) AddChat(chat Chat) error {
	err := s.repo.Add(chat)

	return err
}

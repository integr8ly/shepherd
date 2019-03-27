package usecases

import (
	"reflect"

	"github.com/integr8ly/shepherd/pkg/chat"
	"github.com/integr8ly/shepherd/pkg/gchat"
	"github.com/sirupsen/logrus"
)

type HelpUseCase struct {
}

const (
	helpResp = `available commands are:
@shepherd help

`
)

func (hu *HelpUseCase) GeneralHelp(cmd chat.Command, event interface{}) (string, error) {
	_, ok := event.(*gchat.Event)
	if !ok {
		logrus.Error("failed to handle event. Expected gchat.Event but got ", reflect.TypeOf(event))
		return "sorry could not complete the request.", nil
	}
	return helpResp, nil
}

//func(chat.Command, interface{}) (string, error)

func (hu *HelpUseCase) Register() map[string]chat.ActionHandlerFunc {
	handlers := map[string]chat.ActionHandlerFunc{}
	handlers[chat.CommandGeneralHelp] = hu.GeneralHelp
	return handlers
}

func NewHelpUseCase() *HelpUseCase {
	return &HelpUseCase{}
}

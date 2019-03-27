package gchat

import (
	"github.com/integr8ly/shepherd/pkg/chat"
	"github.com/sirupsen/logrus"
	chat2 "google.golang.org/api/chat/v1"
)

type ActionHandler struct {
	handlers map[string]chat.ActionHandlerFunc
}

func NewActionHandler(chatService *chat2.SpacesService, actionHandlers map[string]chat.ActionHandlerFunc) *ActionHandler {
	return &ActionHandler{
		handlers: actionHandlers,
	}
}

func (ah *ActionHandler) Handle(m chat.Message) string {
	return "hello from shepherd"
}

func (ah *ActionHandler) Platform() string {
	return "hangout"
}

func (ah *ActionHandler) parseCommand(event *Event, argumentText string) (chat.Command, error) {
	cmd := chat.Command{
		Requester:   event.User.DisplayName,
		RequesterID: event.User.Name,
	}

	logrus.Debug("parse command text ", argumentText)
	for k, v := range chat.Commands {
		if k.MatchString(argumentText) {
			logrus.Debug("found command ", k.String())
			v.Requester = event.User.DisplayName
			v.RequesterID = event.User.Name
			sm := k.FindStringSubmatch(argumentText)
			logrus.Debug("sub matches ", sm, len(sm))
			v.Args = sm[1:]
			namedArgs := map[string]string{}
			for i, name := range k.SubexpNames() {
				if name != "" {
					namedArgs[name] = sm[i]
				}
			}
			v.MappedArgs = namedArgs
			return v, nil
		}
	}
	return cmd, chat.NewUknownCommand(argumentText)
}

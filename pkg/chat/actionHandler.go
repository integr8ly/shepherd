package chat

import (
	"regexp"
)

type ActionHandler struct {
	handlers map[string]Handler
}

type ActionHandlerFunc func(Command, interface{}) (string, error)

func (ah *ActionHandler) RegisterHandler(h Handler) {
	ah.handlers[h.Platform()] = h
}

func NewActionHandler() *ActionHandler {
	return &ActionHandler{handlers: map[string]Handler{}}
}

func (ah *ActionHandler) Handle(m Message) string {
	switch m.Platform() {
	case "hangout":
		return ah.handlers["hangout"].Handle(m)

	}
	return ""
}

type Handler interface {
	Handle(m Message) string
	Platform() string
}

type Message interface {
	Platform() string
}

type Command struct {
	ActionType  string
	TeamID      string
	Requester   string
	RequesterID string
	Name        string
	Args        []string
	MappedArgs  map[string]string
	Room        string
}

func (c Command) NoEmptyArgs() bool {
	for _, a := range c.Args {
		if a == "" {
			return false
		}
	}
	return true
}

const (
	CommandAdminHelp   = "admin help"
	CommandGeneralHelp = "general help"
)

var (
	AdminHelpRegexp   = regexp.MustCompile(`^admin help$`)
	GeneralHelpRegexp = regexp.MustCompile(`^help$`)
	Commands          = map[*regexp.Regexp]Command{
		GeneralHelpRegexp: {ActionType: "general", Name: CommandGeneralHelp},
		AdminHelpRegexp:   {ActionType: "admin", Name: CommandAdminHelp},
	}
)

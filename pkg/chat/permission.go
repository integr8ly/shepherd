package chat

import (
	"github.com/integr8ly/shepherd/pkg/domain"
	"github.com/sirupsen/logrus"
)

type Permissions struct {
}

func NewPermissions() *Permissions {
	return &Permissions{}
}

func (p *Permissions) CanUserDoCmd(u *domain.User, cmd Command) (bool, error) {
	logrus.Infof("can user do cmd %s %s ", cmd.ActionType, u.Role)
	if cmd.ActionType == "general" {
		return true, nil
	}

	if u.Role == "admin" {
		return true, nil
	}

	return u.Role == cmd.ActionType, nil
}

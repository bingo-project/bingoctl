package bootstrap

import (
	"{[.RootPackage]}/internal/pkg/facade"
	"{[.RootPackage]}/pkg/mail"
)

func InitMail() {
	facade.Mail = mail.NewMailer(facade.Config.Mail)
}

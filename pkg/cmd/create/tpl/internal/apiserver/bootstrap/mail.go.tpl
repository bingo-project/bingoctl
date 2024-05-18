package bootstrap

import (
	"{[.RootPackage]}/internal/apiserver/facade"
	"{[.RootPackage]}/pkg/mail"
)

func InitMail() {
	facade.Mail = mail.NewMailer(facade.Config.Mail)
}

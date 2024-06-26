package common

import (
	"context"
	"fmt"
	"time"

	"github.com/bingo-project/component-base/log"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/hibiken/asynq"

	"{[.RootPackage]}/internal/apiserver/facade"
	"{[.RootPackage]}/internal/apiserver/global"
	v1 "{[.RootPackage]}/internal/apiserver/http/request/v1"
	"{[.RootPackage]}/internal/apiserver/job"
	"{[.RootPackage]}/internal/apiserver/store"
	"{[.RootPackage]}/internal/pkg/errno"
)

type EmailBiz interface {
	SendEmailVerifyCode(ctx context.Context, req *v1.SendEmailRequest) error
}

type emailBiz struct {
	ds store.IStore
}

var _ EmailBiz = (*emailBiz)(nil)

func NewEmail(ds store.IStore) *emailBiz {
	return &emailBiz{ds: ds}
}

func (b *emailBiz) SendEmailVerifyCode(ctx context.Context, req *v1.SendEmailRequest) error {
	// Check waiting time
	keyWaiting := fmt.Sprintf("%s:%s", global.CacheKeyVerifyCodeWaiting, req.Email)
	exist := facade.Cache.Get(keyWaiting)
	if exist != nil {
		return errno.ErrTooManyRequests
	}

	// Generate code
	code := random.RandNumeral(6)
	subject := "Email Verification code " + code
	msg := fmt.Sprintf("Your verification code is: %s, please note that it will expire in 5 minutes.", code)

	// Email task payload
	payload := &job.EmailVerificationCodePayload{
		To:      req.Email,
		Subject: subject,
		Content: msg,
	}

	// Enqueue email task
	t := asynq.NewTask(job.EmailVerificationCode, []byte(convertor.ToString(payload)))
	_, err := facade.Queue.Enqueue(t)
	if err != nil {
		log.C(ctx).Errorw("enqueue failed", "err", err)

		return err
	}

	// Cache code
	keyTtl := fmt.Sprintf("%s:%s", global.CacheKeyVerifyCodeTtl, req.Email)
	facade.Cache.Set(keyTtl, code, time.Minute*time.Duration(facade.Config.Code.TTL))
	facade.Cache.Set(keyWaiting, code, time.Minute*time.Duration(facade.Config.Code.Waiting))

	log.C(ctx).Infow("SendEmailVerifyCode succeed", "email", req.Email, "msg", msg)

	return nil
}

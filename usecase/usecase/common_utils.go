package usecase

import (
	"fmt"
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type CommonUtils struct {
	p   ipresenter.IPresenter
	lr  igateway.ILineRepository
}

// Newコンストラクタ
func NewCommonUtils(
	p ipresenter.IPresenter,
	lr igateway.ILineRepository) *CommonUtils {
	return &CommonUtils{p: p, lr: lr}
}

func (cu *CommonUtils) SendReplyMessage(
	token string,
	m []map[string]interface{},
	botType dto.BotType) error {
	rm := dto.NewReplyMessage(token, m)
	if err := cu.p.ReplyMessage(rm, botType); err != nil {
		return err
	}
	return nil
}

func (cu *CommonUtils) MulticastMessage(
	users []entity.User,
	m []map[string]interface{},
	botType dto.BotType) error {
	userCnt := len(users)
	quotaComsuption, err := cu.lr.GetQuotaComsuption(botType)
	if err != nil {
		return fmt.Errorf("failed to get the quota comsuption; err = %w", err)
	}
	// To avoid sending more than 1000 notifications 
	// https://www.linebiz.com/jp/service/line-official-account/plan/
	conf.Log.Info("publish message counts", zap.Int("user count", userCnt), zap.Int("message cnt", len(m)), zap.Int64("Quota Comsuption", quotaComsuption))
	if (userCnt * (len(m) / 3 + 1)  + int(quotaComsuption) <= 1000) {
		userIds := make([]string, userCnt)
		for i, user := range users {
			userIds[i] = user.Id
		}
		pm := dto.NewMulticastMessage(userIds, m)
		if err = cu.p.MulticastMessage(pm, botType); err != nil {
			return fmt.Errorf("failed to send the multicast message; err = %w", err)
		}
	} else {
		return fmt.Errorf("there is no remaining quota; userCnt = %d, quotaComsuption = %d", userCnt, quotaComsuption)
	}
	return nil
}
package wedding

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"go.uber.org/zap"
)

type WeddingPushUsecase struct {
	p ipresenter.IPresenter
	mr igateway.IMessageRepository
	lr igateway.ILineRepository
	ur igateway.IUserRepository
}

// Newコンストラクタ
func NewWeddingPushUsecase(
	p ipresenter.IPresenter,
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	ur igateway.IUserRepository) *WeddingPushUsecase {
	return &WeddingPushUsecase{p:p, mr:mr, lr:lr, ur:ur}
}

func (wpu *WeddingPushUsecase) PublishReminder() error {
	messages := wpu.mr.FindMessageByKey(dto.WeddingBotType, "reminder")
	users, err := wpu.ur.FindByWillJoin(true)
	if err != nil {
		return fmt.Errorf("failed to get user by WillJoin; err = %w", err)
	}
	userCnt := len(*users)
	quotaComsuption, err := wpu.lr.GetQuotaComsuption(dto.WeddingBotType)
	if err != nil {
		return fmt.Errorf("failed to get the quota comsuption; err = %w", err)
	}
	// To avoid sending more than 1000 notifications 
	// https://www.linebiz.com/jp/service/line-official-account/plan/
	conf.Log.Info("publish message counts", zap.Int("user count", userCnt), zap.Int("message cnt", len(messages)), zap.Int64("Quota Comsuption", quotaComsuption))
	if (userCnt * (len(messages) / 3 + 1)  + int(quotaComsuption) <= 1000) {
		userIds := make([]string, userCnt)
		for i, user := range *users {
			userIds[i] = user.Id
		}
		pm := dto.NewMulticastMessage(userIds, messages)
		if err = wpu.p.MulticastMessage(pm, dto.WeddingBotType); err != nil {
			return fmt.Errorf("failed to send the multicast message; err = %w", err)
		}
	} else {
		return fmt.Errorf("there is no remaining quota; userCnt = %d, quotaComsuption = %d", userCnt, quotaComsuption)
	}
	return nil
}
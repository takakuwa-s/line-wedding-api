package wedding

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
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

func (wpu *WeddingPushUsecase) SendReminder() error {
	message := wpu.mr.FindReminderMessage()
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
	if (userCnt * (len(message) / 3 + 1)  + quotaComsuption <= 1000) {
		userIds := make([]string, userCnt)
		for i, user := range *users {
			userIds[i] = user.Id
		}
		pm := dto.NewMulticastMessage(userIds, message)
		if err = wpu.p.MulticastMessage(pm, dto.WeddingBotType); err != nil {
			return fmt.Errorf("failed to send the multicast message; err = %w", err)
		}
	} else {
		return fmt.Errorf("there is no remaining quota; userCnt = %d, quotaComsuption = %d", userCnt, quotaComsuption)
	}
	return nil
}
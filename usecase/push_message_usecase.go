package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
)


type PushMessageUsecase struct {
	p ipresenter.IPresenter
	mr igateway.IMessageRepository
	lr igateway.ILineRepository
	ur igateway.IUserRepository
}

// Newコンストラクタ
func NewPushMessageUsecase(
	p ipresenter.IPresenter,
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	ur igateway.IUserRepository) *PushMessageUsecase {
	return &PushMessageUsecase{p:p, mr:mr, lr:lr, ur:ur}
}

func (pmu *PushMessageUsecase) SendReminder() error {
	message := pmu.mr.FindReminderMessage()
	users, err := pmu.ur.FindByWillJoin(true)
	if err != nil {
		return fmt.Errorf("failed to get user by WillJoin; err = %w", err)
	}
	userCnt := len(*users)
	quotaComsuption, err := pmu.lr.GetQuotaComsuption()
	if err != nil {
		return fmt.Errorf("failed to get the quota comsuption; err = %w", err)
	}
	if (userCnt * (len(message) / 3 + 1)  + quotaComsuption <= 1000) {
		userIds := make([]string, userCnt)
		for i, user := range *users {
			userIds[i] = user.Id
		}
		pm := dto.NewMulticastMessage(userIds, message)
		if err = pmu.p.MulticastMessage(pm); err != nil {
			return fmt.Errorf("failed to send the multicast message; err = %w", err)
		}
	} else {
		return fmt.Errorf("there is no remaining quota; userCnt = %d, quotaComsuption = %d", userCnt, quotaComsuption)
	}
	return nil
}
package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
)

type WeddingPushUsecase struct {
	mr igateway.IMessageRepository
	ur igateway.IUserRepository
	cu *CommonUtils
}

// Newコンストラクタ
func NewWeddingPushUsecase(
	mr igateway.IMessageRepository,
	ur igateway.IUserRepository,
	cu *CommonUtils) *WeddingPushUsecase {
	return &WeddingPushUsecase{mr:mr, ur:ur, cu:cu}
}

func (wpu *WeddingPushUsecase) PublishInvitation() error {
	messages := wpu.mr.FindMessageByKey(dto.WeddingBotType, "invitation")
	users, err := wpu.ur.FindByWillJoinAndFollowStatus(true, true)
	if err != nil {
		return fmt.Errorf("failed to get user by WillJoin and FollowStatus; err = %w", err)
	}
	return wpu.cu.MulticastMessage(users, messages, dto.WeddingBotType)
}

func (wpu *WeddingPushUsecase) PublishReminder() error {
	messages := wpu.mr.FindMessageByKey(dto.WeddingBotType, "reminder")
	users, err := wpu.ur.FindByWillJoinAndFollowStatus(true, true)
	if err != nil {
		return fmt.Errorf("failed to get user by WillJoin and FollowStatus; err = %w", err)
	}
	return wpu.cu.MulticastMessage(users, messages, dto.WeddingBotType)
}
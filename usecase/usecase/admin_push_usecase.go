package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
)

type AdminPushUsecase struct {
	mr igateway.IMessageRepository
	ur igateway.IUserRepository
	cu *CommonUsecase
}

// Newコンストラクタ
func NewAdminPushUsecase(
	mr igateway.IMessageRepository,
	ur igateway.IUserRepository,
	cu *CommonUsecase) *AdminPushUsecase {
	return &AdminPushUsecase{mr:mr, ur:ur, cu:cu}
}

func (apu *AdminPushUsecase) SendFollowNotification(follower *entity.User) error {
	messages := apu.mr.FindMessageByKey(dto.AdminBotType, "wedding_follow")
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), follower.Name)
	users, err := apu.ur.FindByIsAdmin(true)
	if err != nil {
		return fmt.Errorf("failed to get user by IsAdmin; err = %w", err)
	}
	return apu.cu.MulticastMessage(users, messages, dto.AdminBotType)
}

func (apu *AdminPushUsecase) SendUnFollowNotification(unFollower *entity.User) error {
	messages := apu.mr.FindMessageByKey(dto.AdminBotType, "wedding_unfollow")
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), unFollower.Name)
	users, err := apu.ur.FindByIsAdmin(true)
	if err != nil {
		return fmt.Errorf("failed to get user by IsAdmin; err = %w", err)
	}
	return apu.cu.MulticastMessage(users, messages, dto.AdminBotType)
}
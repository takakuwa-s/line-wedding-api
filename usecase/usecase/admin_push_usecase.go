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
	cu *CommonUtils
}

// Newコンストラクタ
func NewAdminPushUsecase(
	mr igateway.IMessageRepository,
	ur igateway.IUserRepository,
	cu *CommonUtils) *AdminPushUsecase {
	return &AdminPushUsecase{mr:mr, ur:ur, cu:cu}
}

func (apu *AdminPushUsecase) SendFollowNotification(follower *entity.User, isFirst bool) error {
	var messages []map[string]interface{}
	if isFirst {
		messages = apu.mr.FindMessageByKey(dto.AdminBotType, "wedding_follow")
	} else {
		messages = apu.mr.FindMessageByKey(dto.AdminBotType, "wedding_refollow")
	}
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), follower.Name)
	messages[1]["originalContentUrl"] = fmt.Sprintf(messages[1]["originalContentUrl"].(string), follower.IconUrl)
	messages[1]["previewImageUrl"] = fmt.Sprintf(messages[1]["previewImageUrl"].(string), follower.IconUrl)
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
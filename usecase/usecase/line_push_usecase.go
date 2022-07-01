package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"go.uber.org/zap"
)

type LinePushUsecase struct {
	mr igateway.IMessageRepository
	ur igateway.IUserRepository
	p  ipresenter.IPresenter
	lg igateway.ILineGateway
}

// Newコンストラクタ
func NewLinePushUsecase(
	mr igateway.IMessageRepository,
	ur igateway.IUserRepository,
	p ipresenter.IPresenter,
	lg igateway.ILineGateway) *LinePushUsecase {
	return &LinePushUsecase{mr: mr, ur: ur, p: p, lg: lg}
}

func (lpu *LinePushUsecase) SendSlideshowSuccessNotification(contentUrl, thumbUrl string) error {
	messages := lpu.mr.FindMessageByKey("slideshow_create_success")
	messages[1]["originalContentUrl"] = fmt.Sprintf(messages[1]["originalContentUrl"].(string), contentUrl)
	messages[1]["previewImageUrl"] = fmt.Sprintf(messages[1]["previewImageUrl"].(string), thumbUrl)
	return lpu.multicastMessageToAdmin(messages)
}

func (lpu *LinePushUsecase) SendSlideshowErrorNotification(note string) error {
	messages := lpu.mr.FindMessageByKey("slideshow_create_error")
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), note)
	return lpu.multicastMessageToAdmin(messages)
}

func (lpu *LinePushUsecase) PublishMessageToAttendee(messages []map[string]interface{}) error {
	users, err := lpu.ur.FindByAttendanceAndFollow(true, true)
	if err != nil {
		return err
	}
	return lpu.multicastMessage(users, messages)
}

func (lpu *LinePushUsecase) SendFollowNotification(follower *entity.User, isFirst bool) error {
	var messages []map[string]interface{}
	if isFirst {
		messages = lpu.mr.FindMessageByKey("wedding_follow")
	} else {
		messages = lpu.mr.FindMessageByKey("wedding_refollow")
	}
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), follower.LineName)
	messages[1]["originalContentUrl"] = fmt.Sprintf(messages[1]["originalContentUrl"].(string), follower.IconUrl)
	messages[1]["previewImageUrl"] = fmt.Sprintf(messages[1]["previewImageUrl"].(string), follower.IconUrl)
	return lpu.multicastMessageToAdmin(messages)
}

func (lpu *LinePushUsecase) SendUnFollowNotification(unFollower *entity.User) error {
	messages := lpu.mr.FindMessageByKey("wedding_unfollow")
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), unFollower.LineName)
	return lpu.multicastMessageToAdmin(messages)
}

func (lpu *LinePushUsecase) multicastMessageToAdmin(m []map[string]interface{}) error {
	users, err := lpu.ur.FindByIsAdmin(true)
	if err != nil {
		return err
	}
	return lpu.multicastMessage(users, m)
}

func (lpu *LinePushUsecase) multicastMessage(
	users []entity.User,
	m []map[string]interface{}) error {
	userCnt := len(users)
	if userCnt > 500 {
		// https://developers.line.biz/ja/reference/messaging-api/#send-multicast-request-body
		return fmt.Errorf("userCnt is more than 500 limitation; userCnt = %d", userCnt)
	}

	quotaComsuption, err := lpu.lg.GetQuotaComsuption()
	if err != nil {
		return err
	}

	// To avoid sending more than 1000 notifications
	// https://www.linebiz.com/jp/service/line-official-account/plan/
	conf.Log.Info("publish message counts", zap.Int("user count", userCnt), zap.Int("message cnt", len(m)), zap.Int64("Quota Comsuption", quotaComsuption))
	if userCnt*(len(m)/3+1)+int(quotaComsuption) <= 1000 {
		userIds := make([]string, userCnt)
		for i, user := range users {
			userIds[i] = user.Id
		}
		if err = lpu.p.MulticastMessage(userIds, m); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("there is no remaining quota; userCnt = %d, quotaComsuption = %d", userCnt, quotaComsuption)
	}
	return nil
}

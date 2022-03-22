package gateway

import (
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
)

type UserRepository struct {
}

// Newコンストラクタ
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) SaveUser(user *dto.User) {
	conf.Log.Info("Successfully save the user", zap.Any("user", user))
}

func (ur *UserRepository) UpdateFollowStatusById(id string, status bool) {
	conf.Log.Info("Successfully update the follow status", zap.String("id", id), zap.Bool("status", status))
}
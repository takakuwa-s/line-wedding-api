package gateway

import (
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type UserRepository struct {
}

// Newコンストラクタ
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) SaveUser(user *entity.User) error {
	conf.Log.Info("Successfully save the user", zap.Any("user", user))
	return nil
}

func (ur *UserRepository) UpdateFollowStatusById(id string, status bool) error {
	conf.Log.Info("Successfully update the follow status", zap.String("id", id), zap.Bool("status", status))
	return nil
}

func (ur *UserRepository) FindByWillJoin(willJoin bool) (*[]entity.User, error) {
	conf.Log.Info("Successfully find the users with WillJoin flag", zap.Bool("WillJoin", willJoin))
	user := &[]entity.User{
		{Id: "U544c7c84c496d89b3f56b034b75f8dae"},
	}
	return user, nil
}
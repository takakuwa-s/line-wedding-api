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

func (ur *UserRepository) FindByWillJoinAndFollowStatus(willJoin, followStatus bool) (*[]entity.User, error) {
	conf.Log.Info("Successfully find the users with WillJoin and FollowStatus flag", zap.Bool("WillJoin", willJoin), zap.Bool("FollowStatus", followStatus))
	user := &[]entity.User{
		{Id: "U544c7c84c496d89b3f56b034b75f8dae"},
	}
	return user, nil
}

func (ur *UserRepository) FindByIsAdmin(isAdmin bool) (*[]entity.User, error) {
	conf.Log.Info("Successfully find the users with IsAdmin flag", zap.Bool("IsAdmin", isAdmin))
	user := &[]entity.User{
		{Id: "U544c7c84c496d89b3f56b034b75f8dae"},
	}
	return user, nil
}

func (ur *UserRepository) FindById(id string) (*entity.User, error) {
	conf.Log.Info("Successfully find the users by Id", zap.String("id", id))
	user := &entity.User{
		Id: "U544c7c84c496d89b3f56b034b75f8dae",
		Name: "たかくわ",
	}
	return user, nil
}
package gateway

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type UserRepository struct {
	f *dto.Firestore
}

// Newコンストラクタ
func NewUserRepository(f *dto.Firestore) *UserRepository {
	return &UserRepository{f: f}
}

func (ur *UserRepository) SaveUser(user *entity.User) error {
	if _, err := ur.f.Firestore.Collection("users").Doc(user.Id).Set(conf.Ctx, user); err != nil {
		return fmt.Errorf("failed adding a new user; user =  %v, err = %w", user, err)
	}
	conf.Log.Info("Successfully save the user", zap.Any("user", user))
	return nil
}

func (ur *UserRepository) UpdateFollowStatusById(id string, status bool) error {
	if _, err := ur.f.Firestore.Collection("users").Doc(id).Update(conf.Ctx, []firestore.Update{
		{
			Path:  "FollowStatus",
			Value: status,
		},
		{
			Path:  "UpdatedAt",
			Value: time.Now(),
		},
		}); err != nil {
		return fmt.Errorf("failed update the user; id =  %s, status =  %t, err = %w", id, status, err)
	}
	conf.Log.Info("Successfully update the follow status", zap.String("id", id), zap.Bool("status", status))
	return nil
}

func (ur *UserRepository) FindById(id string) (*entity.User, error) {
	dsnap, err := ur.f.Firestore.Collection("users").Doc(id).Get(conf.Ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed get a user by id; id = %s, err = %w", id, err)
		}
	}
	var user entity.User
	dsnap.DataTo(&user)
	conf.Log.Info("Successfully find the users by Id", zap.String("id", id), zap.Any("user", user))
	return &user, nil
}


func (ur *UserRepository) executeQuery(query *firestore.Query) ([]entity.User, error) {
	var users []entity.User
	iter := query.Documents(conf.Ctx)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get a user ; err = %w", err)
		}
		var u entity.User
		dsnap.DataTo(&u)
		users = append(users, u)
	}
	return users, nil
}

func (ur *UserRepository) FindByIsAdmin(isAdmin bool) ([]entity.User, error) {
	query := ur.f.Firestore.Collection("users").Where("IsAdmin", "==", isAdmin)
	users, err := ur.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the users with IsAdmin flag", zap.Bool("IsAdmin", isAdmin), zap.Any("users", users))
	return users, nil
}

func (ur *UserRepository) FindByAttendanceAndFollowStatus(attendance, followStatus bool) ([]entity.User, error) {
	query := ur.f.Firestore.Collection("users").Where("Attendance", "==", attendance).Where("FollowStatus", "==", followStatus)
	users, err := ur.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the users with Attendance and FollowStatus flag", zap.Bool("Attendance", attendance), zap.Bool("FollowStatus", followStatus), zap.Any("user", users))
	return users, nil
}
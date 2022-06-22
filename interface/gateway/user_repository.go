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
	cr *CommonRepository
	f  *dto.Firestore
}

// Newコンストラクタ
func NewUserRepository(cr *CommonRepository, f *dto.Firestore) *UserRepository {
	return &UserRepository{cr: cr, f: f}
}

func (ur *UserRepository) SaveUser(user *entity.User) error {
	user.UpdatedAt = time.Now()
	return ur.cr.Save("users", user.Id, user)
}

func (ur *UserRepository) UpdateBoolFieldById(id, field string, val bool) error {
	if _, err := ur.f.Firestore.Collection("users").Doc(id).Update(conf.Ctx, []firestore.Update{
		{
			Path:  field,
			Value: val,
		},
		{
			Path:  "UpdatedAt",
			Value: time.Now(),
		},
	}); err != nil {
		return fmt.Errorf("failed update the user; id =  %s, field = %s, val = %t, err = %w", id, field, val, err)
	}
	conf.Log.Info("Successfully update the user", zap.String("id", id), zap.Bool(field, val))
	return nil
}

func (ur *UserRepository) FindById(id string) (*entity.User, error) {
	dsnap, err := ur.cr.FindById("users", id)
	if err != nil {
		return nil, err
	}
	if dsnap == nil {
		return nil, nil
	}
	var user entity.User
	dsnap.DataTo(&user)
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
	if users == nil {
		return []entity.User{}, nil
	}
	return users, nil
}

func (ur *UserRepository) FindByIds(ids []string) ([]entity.User, error) {
	var users []entity.User
	for _, list := range ur.cr.SplitSlice(ids) {
		query := ur.f.Firestore.Collection("users").Where("Id", "in", list)
		u, err := ur.executeQuery(&query)
		if err != nil {
			return nil, err
		}
		users = append(users, u...)
	}
	conf.Log.Info("Successfully find users with", zap.Int("user count", len(users)), zap.Strings("ids", ids))
	return users, nil
}

func (ur *UserRepository) FindByIsAdmin(isAdmin bool) ([]entity.User, error) {
	query := ur.f.Firestore.Collection("users").Where("IsAdmin", "==", isAdmin)
	users, err := ur.executeQuery(&query)
	if err != nil {
		return nil, fmt.Errorf("failed to find users ; isAdmin = %t, err = %w", isAdmin, err)
	}
	conf.Log.Info("Successfully find the users with IsAdmin flag", zap.Bool("IsAdmin", isAdmin), zap.Any("users", users))
	return users, nil
}

func (ur *UserRepository) FindByAttendanceAndFollow(attendance, follow bool) ([]entity.User, error) {
	query := ur.f.Firestore.Collection("users").Where("Attendance", "==", attendance).Where("Follow", "==", follow)
	users, err := ur.executeQuery(&query)
	if err != nil {
		return nil, fmt.Errorf("failed to find users ; attendance = %t, follow = %t, err = %w", attendance, follow, err)
	}
	conf.Log.Info("Successfully find the users with Attendance and follow flag", zap.Bool("Attendance", attendance), zap.Bool("follow", follow), zap.Any("user", users))
	return users, nil
}

func (ur *UserRepository) FindByFlagOrderByName(limit int, startId, flag string, val bool) ([]entity.User, error) {
	query := ur.f.Firestore.Collection("users").OrderBy("FamilyNameKana", firestore.Asc).OrderBy("FirstNameKana", firestore.Asc)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if startId != "" {
		dsnap, err := ur.f.Firestore.Collection("users").Doc(startId).Get(conf.Ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return []entity.User{}, nil
			} else {
				return nil, fmt.Errorf("failed to get the users by startId; id =  %s err = %w", startId, err)
			}
		}
		query = query.StartAfter(dsnap)
	}
	if flag != "" {
		query = query.Where(flag, "==", val)
	}
	users, err := ur.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the users",
		zap.Any("user", users),
		zap.Int("limit", limit),
		zap.String("startId", startId),
		zap.String("flag", flag),
		zap.Bool("val", val))
	return users, nil
}

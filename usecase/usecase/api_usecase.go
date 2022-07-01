package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"go.uber.org/zap"
)

type ApiUsecase struct {
	mr  igateway.IMessageRepository
	ur  igateway.IUserRepository
	lg  igateway.ILineGateway
	fr  igateway.IFileRepository
	br  igateway.IBinaryRepository
	lpu *LinePushUsecase
	su  *SlideShowUsecase
}

// Newコンストラクタ
func NewApiUsecase(
	mr igateway.IMessageRepository,
	ur igateway.IUserRepository,
	lg igateway.ILineGateway,
	fr igateway.IFileRepository,
	br igateway.IBinaryRepository,
	lpu *LinePushUsecase,
	su *SlideShowUsecase) *ApiUsecase {
	return &ApiUsecase{mr: mr, ur: ur, lg: lg, fr: fr, br: br, lpu: lpu, su: su}
}

func (au *ApiUsecase) GetInitialData(id string) (*dto.InitApiResponse, error) {
	// Get user
	user, err := au.ur.FindById(id)
	if err != nil {
		return nil, err
	}
	// Get file list
	uploaded := true
	files, err := au.GetFileList(12, "", "", "", "", &uploaded, user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return dto.NewInitApiResponse(user, files), nil
}

func (au *ApiUsecase) UpdateUser(r *dto.UpdateUserRequest) (*entity.User, error) {
	// Get user
	user, err := au.ur.FindById(r.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = au.lg.GetUserProfileById(r.Id)
		if err != nil {
			return nil, err
		}
	}
	user = r.ToUser(user)
	user.Registered = true
	if err = au.ur.SaveUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (au *ApiUsecase) PatchUser(userId, field string, val bool) error {
	// Check if user exists
	if _, err := au.GetUser(userId); err != nil {
		return err
	}

	return au.ur.UpdateBoolFieldById(userId, field, val)
}

func (au *ApiUsecase) GetUser(id string) (*entity.User, error) {
	// Get user
	user, err := au.ur.FindById(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("not found the user with id = %s", id)
	}
	return user, nil
}

func (au *ApiUsecase) GetUsers(limit int, startId, flag string, val bool) ([]entity.User, error) {
	// Get users
	users, err := au.ur.FindByFlagOrderByName(limit, startId, flag, val)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (au *ApiUsecase) GetFileList(limit int, startId, userId, orderBy, fileType string, uploaded *bool, needCreaterName bool) ([]dto.FileResponce, error) {
	if orderBy == "" {
		orderBy = "UpdatedAt"
	}
	files, err := au.fr.FindByLimitAndStartIdAndUserIdAndFileTypeAndUploaded(limit, startId, userId, orderBy, fileType, uploaded)
	if err != nil {
		return nil, err
	}
	uMap := make(map[string]string)
	if needCreaterName {
		set := make(map[string]struct{})
		for _, f := range files {
			set[f.Creater] = struct{}{}
		}
		var ids []string
		for id, _ := range set {
			ids = append(ids, id)
		}
		users, err := au.ur.FindByIds(ids)
		conf.Log.Info("a", zap.Any("set", set), zap.Any("ids", ids), zap.Any("users", users))
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			uMap[u.Id] = u.FamilyName + u.FirstName
		}
	}
	return dto.NewFileResponceList(files, uMap), nil
}

func (au *ApiUsecase) DeleteFile(id string) error {
	file, err := au.fr.FindById(id)
	if err != nil {
		return err
	}
	if file == nil {
		return fmt.Errorf("not found the file with id = %s", id)
	}
	if err := au.fr.DeleteFileById(id); err != nil {
		return err
	}
	if file.Uploaded {
		if err := au.br.DeleteBinary(file.Name, string(file.FileType)); err != nil {
			return err
		}
	}
	return nil
}

func (au *ApiUsecase) DeleteFileList(ids []string) error {
	files, err := au.fr.FindByIds(ids)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("not found the files with ids = %s", ids)
	}
	if len(files) != len(ids) {
		return fmt.Errorf("some ids is invalid and not found; ids = %s", ids)
	}
	for _, f := range files {
		if f.Uploaded {
			if err := au.br.DeleteBinary(f.Name, string(f.FileType)); err != nil {
				conf.Log.Error("failed to delete the file binary", zap.Any("file", f))
			}
		}
		if err := au.fr.DeleteFileById(f.Id); err != nil {
			conf.Log.Error("failed to delete the file metadata", zap.Any("file", f))
		}
	}
	return nil
}

func (au *ApiUsecase) PublishMessageToAttendee(messageKey string) error {
	var messages []map[string]interface{}
	if messageKey == "slideshow" {
		var err error
		messages, err = au.su.CreateSlideshowMessage()
		if err != nil {
			return err
		}
	} else {
		messages = au.mr.FindMessageByKey(messageKey)
	}
	if len(messages) == 0 {
		return fmt.Errorf("not found the message; %v", messageKey)
	}
	return au.lpu.PublishMessageToAttendee(messages)
}

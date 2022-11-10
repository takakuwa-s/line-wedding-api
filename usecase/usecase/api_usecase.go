package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
)

type ApiUsecase struct {
	mr  igateway.IMessageRepository
	ur  igateway.IUserRepository
	lg  igateway.ILineGateway
	fr  igateway.IFileRepository
	br  igateway.IBinaryRepository
	bpg igateway.IBackgroundProcessGateway
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
	bpg igateway.IBackgroundProcessGateway,
	lpu *LinePushUsecase,
	su *SlideShowUsecase) *ApiUsecase {
	return &ApiUsecase{mr: mr, ur: ur, lg: lg, fr: fr, br: br, bpg: bpg, lpu: lpu, su: su}
}

func (au *ApiUsecase) GetInitialData(id string) (*dto.InitApiResponse, error) {
	// Get user
	user, err := au.ur.FindById(id)
	if err != nil {
		return nil, err
	}
	// Get file list
	fileStatuses := []string{string(entity.Open), string(entity.Uploaded)}
	forBrideAndGroom := false
	files, err := au.GetFileList(12, "", "", "", "", fileStatuses, &forBrideAndGroom, user.IsAdmin)
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

func (au *ApiUsecase) GetFileList(limit int, startId, userId, orderBy, fileType string, fileStatuses []string, forBrideAndGroom *bool, needCreaterName bool) ([]dto.FileResponce, error) {
	if orderBy == "" {
		orderBy = "UpdatedAt"
	}
	files, err := au.fr.FindByLimitAndStartIdAndUserIdAndFileTypeAndForBrideAndGroomAndFileStatusIn(
		limit, startId, userId, orderBy, fileType, forBrideAndGroom, fileStatuses)
	if err != nil {
		return nil, err
	}
	uMap := make(map[string]string)
	if needCreaterName && len(files) > 0 {
		set := make(map[string]struct{})
		for _, f := range files {
			set[f.Creater] = struct{}{}
		}
		var ids []string
		for id, _ := range set {
			ids = append(ids, id)
		}
		users, err := au.ur.FindByIds(ids)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			uMap[u.Id] = u.FamilyName + u.FirstName
		}
	}
	return dto.NewFileResponceList(files, uMap), nil
}

func (au *ApiUsecase) DeleteFileList(ids []string) error {
	err := au.fr.UpdateFileStatusByIdIn(entity.Deleted, ids)
	if err != nil {
		return err
	}
	if err := au.bpg.StartDeletingFiles(ids); err != nil {
		return err
	}
	return nil
}

func (au *ApiUsecase) PatchFile(id string, forBrideAndGroom bool) error {
	file, err := au.fr.FindById(id)
	if err != nil {
		return err
	}
	if file == nil {
		return fmt.Errorf("not found the file with id = %s", id)
	}
	if err := au.fr.UpdateForBrideAndGroomById(forBrideAndGroom, id); err != nil {
		return err
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
